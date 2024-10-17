package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"

	"solana-program-scanner/block_height_manager"
)

type BlockGetter struct {
	workerNumber int
	cli          rpc.RpcClient
	bhm          block_height_manager.BlockHeightManager
	fc           FlowController
}

const (
	Delay = time.Second * 10
)

var (
	transactionVersion = uint8(0)
	rewards            = false
	getBlockConfig     = rpc.GetBlockConfig{
		Encoding:                       rpc.GetBlockConfigEncodingJsonParsed,
		TransactionDetails:             rpc.GetBlockConfigTransactionDetailsFull,
		Rewards:                        &rewards,
		Commitment:                     rpc.CommitmentFinalized,
		MaxSupportedTransactionVersion: &transactionVersion,
	}
)

func NewBlockGetter(workerNumber int, bhm block_height_manager.BlockHeightManager, fc FlowController) *BlockGetter {
	cli := rpc.NewRpcClient(conf.Solana.RpcEndpoint)
	return &BlockGetter{
		workerNumber: workerNumber,
		cli:          cli,
		bhm:          bhm,
		fc:           fc,
	}
}

func (bg *BlockGetter) getBlockHeightBySlot(slot uint64) int64 {
	ctx := context.Background()
	resp, err := bg.cli.GetBlockWithConfig(ctx, slot, getBlockConfig)
	if err != nil {
		Logger.Fatal(err.Error())
	}

	if resp.Error != nil {
		Logger.Fatal(resp.Error.Error())
	}

	return *resp.Result.BlockHeight
}

func (bg *BlockGetter) run(ctx context.Context, wg *sync.WaitGroup, taskCh chan uint64, blockCh chan *rpc.GetBlock) {
	defer wg.Done()
	defer close(blockCh)

	var wgWorker sync.WaitGroup
	wgWorker.Add(bg.workerNumber)

	worker := func(id int) {
		defer wgWorker.Done()

		for {
			select {
			case <-ctx.Done():
				return

			case slot := <-taskCh:
				if slot == 0 {
					Logger.Info(fmt.Sprintf("worker:%d all task finish", id))
					return
				}

				Logger.Info(fmt.Sprintf("GetBlock:%d start", slot))
				failCnt := 0
				for {
					resp, err := bg.cli.GetBlockWithConfig(ctx, slot, getBlockConfig)
					if err != nil {
						bg.fc.onErr()
						failCnt += 1
						Logger.Info(fmt.Sprintf("get slot:%d failed:%d with err:%v", slot, failCnt, err))
						continue
					}

					if resp.Error != nil {
						Logger.Error(fmt.Sprintf("TODO check! get slot:%d JsonRpc err:%s", slot, resp.Error.Error()))
						// ignore this slot: TODO check
						bg.fc.onDone(time.Now())
						break
					}

					bg.fc.onDone(time.Now())
					Logger.Info(fmt.Sprintf("get slot:%d succeed", slot))

					for {
						cnt := 0
						if cnt%100 == 0 {
							Logger.Debug(fmt.Sprintf("current height:%d, ParentSlot:%d, height:%d", bg.bhm.Get(), resp.Result.ParentSlot, *resp.Result.BlockHeight))
						}
						cnt += 1

						if bg.bhm.CanCommit(*resp.Result.BlockHeight) {
							blockCh <- resp.Result
							bg.bhm.Commit(*resp.Result.BlockHeight)
							break
						}
						time.Sleep(time.Millisecond * 100)
					}

					break
				}
			}
		}
	}

	for i := 0; i < bg.workerNumber; i++ {
		go worker(i)
	}

	Logger.Info("wait all workers done")
	wgWorker.Wait()
	Logger.Info("all workers done")
}
