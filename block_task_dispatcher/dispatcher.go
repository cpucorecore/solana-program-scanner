package block_task_dispatcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"
	"go.uber.org/zap"

	"solana-program-scanner/config"
	"solana-program-scanner/global_stop"
	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/msg_broker"
	"solana-program-scanner/solana"
)

var getSlotConfig = rpc.GetSlotConfig{Commitment: rpc.Commitment(config.G.Solana.Commitment)}

type TaskDispatcher struct {
	conf         *config.BlockTaskDispatcherConf
	slotProducer msg_broker.MsgProducer[uint64]

	mu        sync.RWMutex
	stopped   bool
	ctx       context.Context
	rpcClient *rpc.RpcClient
}

func New(
	conf *config.BlockTaskDispatcherConf,
	solanaRpcUrl string,
	slotProducer msg_broker.MsgProducer[uint64],
) *TaskDispatcher {
	d := &TaskDispatcher{
		conf:         conf,
		slotProducer: slotProducer,
		ctx:          context.Background(),
		rpcClient:    solana.NewCommonClient(solanaRpcUrl),
	}
	global_stop.Subscribe(d)
	return d
}

func (d *TaskDispatcher) dispatchTasks(startSlot uint64, endSlot uint64, closeAfterDispatch bool) bool {
	monitor.NewestSlot.Set(float64(endSlot))

	slotProducerCh := d.slotProducer.ProducerChan()
	for slot := startSlot; slot <= endSlot; slot++ {
		if d.Stopped() {
			log.Logger.Info("TaskDispatcher stopped", zap.Uint64("slot", slot))
			d.slotProducer.Close()
			return true
		}

		slotProducerCh <- slot
	}

	if closeAfterDispatch {
		d.slotProducer.Close()
	}

	return false
}

func (d *TaskDispatcher) run(startSlot uint64) {
	ticker := time.NewTicker(time.Millisecond * time.Duration(d.conf.GetSlotIntervalByMs))
	defer ticker.Stop()

	var stopped bool
	var endSlot uint64
	for {
		select {
		case <-ticker.C:
		}

		rpcResult, err := d.rpcClient.GetSlotWithConfig(d.ctx, getSlotConfig)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("solanaClient GetSlot err: %s", err.Error()))
			continue
		}
		endSlot = rpcResult.Result

		if endSlot < startSlot {
			continue
		}
		log.Logger.Info(fmt.Sprintf("rpc getSlot: %d", endSlot))

		if stopped = d.dispatchTasks(startSlot, endSlot, false); stopped {
			return
		}

		startSlot = endSlot + 1
	}
}

func (d *TaskDispatcher) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	if !d.conf.Enable {
		log.Logger.Warn("TaskDispatcher is disabled")
		return
	}

	log.Logger.Info(fmt.Sprintf("TaskDispatcher Run with [startSlot/endSlot]=[%d/%d]", d.conf.StartSlot, d.conf.EndSlot))

	startSlot := d.conf.StartSlot
	endSlot := d.conf.EndSlot

	if endSlot > 0 {
		d.dispatchTasks(startSlot, endSlot, true)
		return
	}

	d.run(startSlot)
}

func (d *TaskDispatcher) Stopped() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.stopped
}

func (d *TaskDispatcher) NotifyStopped() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.stopped = true
}
