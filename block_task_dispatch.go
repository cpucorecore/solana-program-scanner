package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"
)

const (
	Commitment = rpc.CommitmentFinalized
)

type BlockTaskDispatch struct {
	cli rpc.RpcClient
	fc  FlowController
}

func NewBlockTaskDispatch(fc FlowController) *BlockTaskDispatch {
	cli := rpc.NewRpcClient(conf.Solana.RpcEndpoint)
	return &BlockTaskDispatch{
		cli: cli,
		fc:  fc,
	}
}

func (btd *BlockTaskDispatch) keepDispatchTaskMock(wg *sync.WaitGroup, startSlot uint64, count uint64, taskCh chan uint64) {
	defer wg.Done()

	start := startSlot
	end := startSlot + count
	for start < end {
		taskCh <- start
		start += 1
	}

	close(taskCh)
}

func (btd *BlockTaskDispatch) keepDispatchingTask(ctx context.Context, wg *sync.WaitGroup, startSlot uint64, count uint64, taskCh chan uint64) {
	defer wg.Done()

	const QueryInterval = time.Second * 10
	endCursor := uint64(0)
	if count > 0 {
		endCursor = startSlot + count
	}

	cursor := startSlot
	taskCh <- cursor
	cursor++

	for {
		resp, err := btd.cli.GetSlotWithConfig(ctx, rpc.GetSlotConfig{Commitment: Commitment})
		if err != nil {
			Logger.Error(fmt.Sprintf("GetSlot err: %s", err.Error()))
			btd.fc.onErr()
		}

		if resp.Error != nil {
			Logger.Error(fmt.Sprintf("GetSlot JsonRpc err:%s", resp.Error.Error()))
			btd.fc.onErr()
		}

		btd.fc.onDone(time.Now())

		if resp.Result >= cursor {
			end := resp.Result
			for ; cursor <= end; cursor++ {
				if cursor == endCursor {
					close(taskCh)
					return
				}
				taskCh <- cursor
			}
		}

		time.Sleep(QueryInterval)
	}
}
