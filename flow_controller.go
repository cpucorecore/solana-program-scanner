package main

import (
	"encoding/json"
	"fmt"
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
	TpsMax      uint
	TpsWaitUnit time.Duration
	ErrWaitUnit time.Duration

	Stop bool

	OkCounter     mutexCounter
	ErrCounter    mutexCounter
	ErrCounterTmp mutexCounter

	TpsCounter TpsCounter
}

type mutexCounter struct {
	Cnt int
	mu  sync.Mutex
}

func (mc *mutexCounter) get() (cnt int) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	cnt = mc.Cnt
	return
}

func (mc *mutexCounter) up() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.Cnt++
}

func (mc *mutexCounter) down() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.Cnt--
}

func (mc *mutexCounter) reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.Cnt = 0
}

func NewFlowController(tpsMax uint, tpsCountWindow int, tpsWaitUnit time.Duration, errWaitUnit time.Duration) FlowController {
	return &flowController{
		TpsMax:      tpsMax,
		TpsWaitUnit: tpsWaitUnit,
		ErrWaitUnit: errWaitUnit,
		TpsCounter:  NewTpsCounter(tpsCountWindow),
	}
}

func (fc *flowController) onDone(timestamp time.Time) {
	fc.ErrCounterTmp.reset()
	fc.OkCounter.up()
	tps := fc.TpsCounter.onDone(timestamp)
	if fc.TpsMax > 0 && tps > float32(fc.TpsMax) {
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
		Logger.Info(fmt.Sprintf("flow controller dump: %s", string(fcJson))) // TODO merge
		Logger.Info(fc.TpsCounter.String())                                  // TODO merge
	}
}

func (fc *flowController) stopLog() {
	fc.Stop = true
}

var _ FlowController = &flowController{}
