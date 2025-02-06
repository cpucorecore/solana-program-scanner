package dispatcher

import (
	"fmt"
	"sync"
	"time"

	"solana-program-scanner/log"
	"solana-program-scanner/msg_broker"
	"solana-program-scanner/parser"
)

type BlockDispatcher struct {
	blockParsedConsumer msg_broker.MsgBroker[*parser.Block]
	dispatchers         []Dispatchable[*parser.Block]
}

func NewBlockDispatcher(blockParsedConsumer msg_broker.MsgBroker[*parser.Block]) *BlockDispatcher {
	return &BlockDispatcher{
		blockParsedConsumer: blockParsedConsumer,
	}
}

func (d *BlockDispatcher) RegisterDispatcher(dispatchers ...Dispatchable[*parser.Block]) {
	d.dispatchers = append(d.dispatchers, dispatchers...)
}

func (d *BlockDispatcher) Stop() {
	for _, dispatcher := range d.dispatchers {
		dispatcher.Stop()
	}
}

func (d *BlockDispatcher) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	dispatchersLen := len(d.dispatchers)
	blockCh := d.blockParsedConsumer.ConsumerChan()
	for {
		select {
		case block, ok := <-blockCh:
			if !ok {
				log.Logger.Info("blockParsedConsumer @ done")
				d.Stop()
				return
			}

			log.Logger.Info(fmt.Sprintf("BlockDispatcher dispatch slot:%d", block.Block.Slot))
			var errCnt int
			for {
				errCnt = 0
				for _, dispatcher := range d.dispatchers {
					err := dispatcher.Dispatch(block)
					if err != nil {
						log.Logger.Error(fmt.Sprintf("BlockDispatcher-%d Dispatch slot:%d err %v", dispatcher.Id(), block.Block.Slot, err))
						errCnt++
					}
				}

				if errCnt == dispatchersLen {
					log.Logger.Warn("all dispatcher failed, will retry")
					time.Sleep(100 * time.Millisecond)
					continue
				}

				break
			}
		}
	}
}
