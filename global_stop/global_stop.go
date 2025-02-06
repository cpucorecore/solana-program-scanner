package global_stop

import (
	"sync"
)

type SubscriptionStopped interface {
	NotifyStopped()
}

type globalStop struct {
	mutex    sync.Mutex
	stopped  bool
	stopSubs []SubscriptionStopped
}

func (gs *globalStop) Register(subscriber SubscriptionStopped) {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.stopSubs = append(gs.stopSubs, subscriber)
}

func (gs *globalStop) Stop() {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	gs.stopped = true
	for _, stopSub := range gs.stopSubs {
		stopSub.NotifyStopped()
	}
}

func (gs *globalStop) Stopped() bool {
	gs.mutex.Lock()
	defer gs.mutex.Unlock()
	return gs.stopped
}

var gStop globalStop

func Stop() {
	gStop.Stop()
}

func Stopped() bool {
	return gStop.Stopped()
}

func Subscribe(subscriber SubscriptionStopped) {
	gStop.Register(subscriber)
}
