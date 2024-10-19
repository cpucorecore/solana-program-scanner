package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FlowController interface {
	onDone(time.Time)
	onErr()
	startLog(time.Duration)
	stopLog()
}

type flowController struct {
	TpsMax         uint
	TpsCountWindow int
	TpsWaitUnit    time.Duration
	ErrWaitUnit    time.Duration

	Stop bool

	OkCounter     mutexCounter
	ErrCounter    mutexCounter
	ErrCounterTmp mutexCounter
	TpsCounter    TpsCounter
}

type tpsCounter struct {
	countWindow int
	timestamps  []time.Time
	mu          sync.Mutex
}

type TpsCounter interface {
	onDone(time.Time) float32
	String() string
}

var _ TpsCounter = &tpsCounter{}

func NewTpsCounter(countWindow int) TpsCounter {
	if countWindow < 10 {
		countWindow = 10
	}

	return &tpsCounter{
		countWindow: countWindow,
	}
}

func (tc *tpsCounter) onDone(timestamp time.Time) (tps float32) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.timestamps = append(tc.timestamps, timestamp)
	reqCnt := len(tc.timestamps)
	if reqCnt == tc.countWindow {
		tps = tc.curTps()
		tc.timestamps = tc.timestamps[1:]
	}

	return
}

func (tc *tpsCounter) curTps() (tps float32) {
	reqCnt := len(tc.timestamps)
	if reqCnt <= 1 {
		return
	}

	timeSecond := float32(1.0*tc.timestamps[reqCnt-1].Sub(tc.timestamps[0])) / float32(time.Second)
	tps = float32(reqCnt) / timeSecond
	return
}

func (tc *tpsCounter) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for _, timestamp := range tc.timestamps {
		sb.WriteString(strconv.FormatInt(timestamp.UnixMilli(), 10))
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")

	reqCnt := len(tc.timestamps)
	if reqCnt >= 2 {
		sb.WriteString(fmt.Sprintf("tps:%v", tc.curTps()))
	}

	return sb.String()
}

type mutexCounter struct {
	cnt int
	mu  sync.Mutex
}

func (mc *mutexCounter) get() (cnt int) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	cnt = mc.cnt
	return
}

func (mc *mutexCounter) up() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.cnt++
}

func (mc *mutexCounter) down() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.cnt--
}

func (mc *mutexCounter) reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.cnt = 0
}

func NewFlowController(tpsMax uint, tpsCountWindow int, tpsWaitUnit time.Duration, errWaitUnit time.Duration) FlowController {
	return &flowController{
		TpsMax:         tpsMax,
		TpsCountWindow: tpsCountWindow,
		TpsWaitUnit:    tpsWaitUnit,
		ErrWaitUnit:    errWaitUnit,
		TpsCounter:     NewTpsCounter(tpsCountWindow),
	}
}

func (fc *flowController) onDone(timestamp time.Time) {
	fc.ErrCounterTmp.reset()
	fc.OkCounter.up()
	tps := fc.TpsCounter.onDone(timestamp)
	if tps > float32(fc.TpsMax) {
		time.Sleep(time.Millisecond * 100) // TODO calc the sleep time
	}
}

func (fc *flowController) onErr() {
	fc.ErrCounterTmp.up()
	fc.ErrCounter.up()
	time.Sleep(fc.ErrWaitUnit * time.Duration(fc.ErrCounterTmp.get()))
}

func (fc *flowController) startLog(logInterval time.Duration) {
	fc.Stop = false
	for !fc.Stop {
		time.Sleep(logInterval)
		fcJson, _ := json.Marshal(fc)
		Logger.Info(fmt.Sprintf("flow controller dump: %s", string(fcJson)))
	}
}

func (fc *flowController) stopLog() {
	fc.Stop = true
}

var _ FlowController = &flowController{}
