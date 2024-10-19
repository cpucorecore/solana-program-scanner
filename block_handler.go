package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/blocto/solana-go-sdk/rpc"
)

type BlockHandler struct {
	blockCh    chan *rpc.GetBlock
	processors []ObjProcessor[*rpc.GetBlock]
}

func (bh *BlockHandler) registerProcessor(processor ObjProcessor[*rpc.GetBlock]) {
	bh.processors = append(bh.processors, processor)
}

func (bh *BlockHandler) keepHandling(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		for _, processor := range bh.processors {
			processor.done()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case block := <-bh.blockCh:
			if block == nil {
				Logger.Info("blockCh @ done")
				return
			}

			for _, processor := range bh.processors {
				err := processor.process(block)
				if err != nil {
					Logger.Error(fmt.Sprintf("%s process block:%d failed", processor.id(), block.BlockHeight))
					return
				}
			}
		}
	}
}

func NewBlockHandler(blockCh chan *rpc.GetBlock) *BlockHandler {
	return &BlockHandler{
		blockCh: blockCh,
	}
}
