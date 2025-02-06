package rate_limiter

import (
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"solana-program-scanner/config"
	"solana-program-scanner/getter/pubsub"
	"solana-program-scanner/monitor"
	"solana-program-scanner/types"
	"solana-program-scanner/utils"
)

type RpcRateLimiter interface {
	SlotBegin(slot uint64)
	SlotDone(slot uint64)
	SlotErr(slot uint64, errCode int)
	AccountErr(account string, errCode int)

	pubsub.SlotBeginSub
	pubsub.SlotDoneSub
	pubsub.SlotErrSub
	pubsub.GetMarketErrSub
}

type rpcRateLimiter struct {
	errWaitUnit      time.Duration
	errWaitUnitQuick time.Duration

	mu            sync.Mutex
	slotStartTime map[uint64]time.Time

	errCounter utils.MutexCounter
}

func New(conf *config.RpcRateLimiterConf) RpcRateLimiter {
	return &rpcRateLimiter{
		errWaitUnit:      time.Millisecond * time.Duration(conf.ErrWaitUnitByMs),
		errWaitUnitQuick: time.Millisecond * time.Duration(conf.ErrWaitUnitByMsQuick),
		slotStartTime:    make(map[uint64]time.Time),
	}
}

func (l *rpcRateLimiter) SlotBegin(slot uint64) {
	monitor.CurrentSlot.Set(float64(slot))

	l.mu.Lock()
	defer l.mu.Unlock()
	l.slotStartTime[slot] = time.Now()
}

func (l *rpcRateLimiter) SlotDone(slot uint64) {
	l.errCounter.Reset()

	l.mu.Lock()
	defer l.mu.Unlock()
	monitor.GetBlockDuration.Observe(float64(time.Since(l.slotStartTime[slot]).Milliseconds()))
	delete(l.slotStartTime, slot)
}

func (l *rpcRateLimiter) SlotErr(slot uint64, errCode int) {
	l.mu.Lock()
	monitor.GetBlockFailedDuration.Observe(float64(time.Since(l.slotStartTime[slot]).Milliseconds()))
	delete(l.slotStartTime, slot)
	l.mu.Unlock()
	monitor.GetBlockErrCounter.With(prometheus.Labels{"err_code": strconv.Itoa(errCode)}).Inc()

	l.waitByErr(errCode)
}

func (l *rpcRateLimiter) AccountErr(account string, errCode int) {
	l.waitByErr(errCode)
}

func (l *rpcRateLimiter) PubSlotBegin(slot uint64) {
	l.SlotBegin(slot)
}

func (l *rpcRateLimiter) PubSlotDone(slot uint64) {
	l.SlotDone(slot)
}

func (l *rpcRateLimiter) PubSlotErr(slot uint64, errCode int) {
	l.SlotErr(slot, errCode)
}

func (l *rpcRateLimiter) PubGetMarketErr(marketAddr string, errCode int) {
	l.AccountErr(marketAddr, errCode)
}

func (l *rpcRateLimiter) waitByErr(errType int) {
	errCnt := l.errCounter.Up()
	switch errType {
	case types.RpcErr, types.RpcErrRpsLimit:
		time.Sleep(l.errWaitUnit * time.Duration(1<<errCnt))

	case types.RpcErrSlotCleanup, types.RpcErrSlotNotAvailable:
		time.Sleep(l.errWaitUnitQuick * time.Duration(errCnt))
	}
}

var _ RpcRateLimiter = &rpcRateLimiter{}
