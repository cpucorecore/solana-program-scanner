package main

import (
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
	targetTps      int
	tpsCountWindow int
	tpsWaitUnit    time.Duration
	errWaitUnit    time.Duration

	stop bool

	okCounter     mutexCounter
	errCounter    mutexCounter
	errCounterTmp mutexCounter
	tpsCounter    TpsCounter
}

type tpsCounter struct {
	countWindow int
	timestamps  []time.Time
	mu          sync.Mutex
}

type TpsCounter interface {
	onDone(time.Time) int
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

func (mtt *tpsCounter) onDone(timestamp time.Time) (tps int) {
	mtt.mu.Lock()
	defer mtt.mu.Unlock()

	mtt.timestamps = append(mtt.timestamps, timestamp)
	reqCnt := len(mtt.timestamps)
	if reqCnt == mtt.countWindow {
		tps = reqCnt / int(mtt.timestamps[reqCnt-1].Sub(mtt.timestamps[0])/time.Second)
		mtt.timestamps = mtt.timestamps[1:]
	}

	return
}

func (mtt *tpsCounter) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for _, timestamp := range mtt.timestamps {
		sb.WriteString(strconv.FormatInt(timestamp.UnixMilli(), 10))
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")

	reqCnt := len(mtt.timestamps)
	if reqCnt >= 2 {
		sb.WriteString(fmt.Sprintf("tps:%v", reqCnt*1.0/int(mtt.timestamps[reqCnt-1].Sub(mtt.timestamps[0])/time.Second)))
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

func NewFlowController(targetTps int, tpsCountWindow int, tpsWaitUnit time.Duration, errWaitUnit time.Duration) FlowController {
	return &flowController{
		targetTps:      targetTps,
		tpsCountWindow: tpsCountWindow,
		tpsWaitUnit:    tpsWaitUnit,
		errWaitUnit:    errWaitUnit,
		tpsCounter:     NewTpsCounter(tpsCountWindow),
	}
}

func (fc *flowController) onDone(timestamp time.Time) {
	fc.errCounterTmp.reset()
	fc.okCounter.up()
	tps := fc.tpsCounter.onDone(timestamp)
	if tps > fc.targetTps {
		time.Sleep(time.Millisecond * 100) // TODO calc the sleep time
	}
}

func (fc *flowController) onErr() {
	fc.errCounterTmp.up()
	fc.errCounter.up()
	time.Sleep(fc.errWaitUnit * time.Duration(fc.errCounterTmp.get()))
}

func (fc *flowController) startLog(logInterval time.Duration) {
	fc.stop = false
	for !fc.stop {
		time.Sleep(logInterval)
		Logger.Info(fmt.Sprintf(`flow controller summary:
okCounter:%d
errCounter:%d
errCounterTmp:%d
tpsCounter:
%s
`, fc.okCounter.get(), fc.errCounter.get(), fc.errCounterTmp.get(), fc.tpsCounter))
	}
}

func (fc *flowController) stopLog() {
	fc.stop = true
}

var _ FlowController = &flowController{}
