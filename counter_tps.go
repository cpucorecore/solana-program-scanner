package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type tpsCounter struct {
	CountWindow int
	Timestamps  []time.Time
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
		CountWindow: countWindow,
	}
}

func (tc *tpsCounter) onDone(timestamp time.Time) (tps float32) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	tc.Timestamps = append(tc.Timestamps, timestamp)
	reqCnt := len(tc.Timestamps)
	if reqCnt == tc.CountWindow {
		tps = tc.curTps()
		tc.Timestamps = tc.Timestamps[1:]
	}

	return
}

func (tc *tpsCounter) curTps() (tps float32) {
	reqCnt := len(tc.Timestamps)
	if reqCnt <= 1 {
		return
	}

	timeSecond := float32(1.0*tc.Timestamps[reqCnt-1].Sub(tc.Timestamps[0])) / float32(time.Second)
	tps = float32(reqCnt) / timeSecond
	return
}

func (tc *tpsCounter) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for _, timestamp := range tc.Timestamps {
		sb.WriteString(strconv.FormatInt(timestamp.UnixMilli(), 10))
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")

	reqCnt := len(tc.Timestamps)
	if reqCnt >= 2 {
		sb.WriteString(fmt.Sprintf("tps:%v", tc.curTps()))
	}

	return sb.String()
}
