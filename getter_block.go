package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"

	"solana-program-scanner/block_height_manager"
)

type GetterBlock struct {
	workerNumber int
	cli          rpc.RpcClient
	bhm          block_height_manager.BlockHeightManager
	fc           FlowController
}

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

func NewGetterBlock(workerNumber int, bhm block_height_manager.BlockHeightManager, fc FlowController) *GetterBlock {
	cli := rpc.NewRpcClient(gc.Solana.RpcEndpoint)
	return &GetterBlock{
		workerNumber: workerNumber,
		cli:          cli,
		bhm:          bhm,
		fc:           fc,
	}
}

func (bg *GetterBlock) getBlockHeight(slot uint64) int64 {
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

func (bg *GetterBlock) keepBlockGetting(ctx context.Context, wg *sync.WaitGroup, taskCh chan uint64, blockCh chan *rpc.GetBlock) {
	defer wg.Done()
	defer func() {
		close(blockCh)
	}()

	const queryDuration = time.Millisecond * 10
	const logDuration = time.Second * 10
	const logPoint = int(logDuration / queryDuration)

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

				Logger.Info(fmt.Sprintf("id:%d GetBlock:%d start", id, slot))
				failCnt := 0
				for {
					resp, err := bg.cli.GetBlockWithConfig(ctx, slot, getBlockConfig)
					if err != nil {
						bg.fc.onErr()
						failCnt += 1
						Logger.Info(fmt.Sprintf("GetBlock:%d failed:%d with err:%v", slot, failCnt, err))
						continue
					}

					if resp.Error != nil {
						Logger.Error(fmt.Sprintf("id:%d TODO check! get slot:%d JsonRpc err:%s", id, slot, resp.Error.Error()))
						break
					}

					bg.fc.onDone(time.Now())
					Logger.Info(fmt.Sprintf("id:%d GetBlock:%d:%d succeed", id, slot, *resp.Result.BlockHeight))

					cnt := 0
					for {
						if cnt%logPoint == 0 {
							Logger.Info(fmt.Sprintf("id:%d current height:%d, ParentSlot:%d, height:%d", id, bg.bhm.Get(), resp.Result.ParentSlot, *resp.Result.BlockHeight))
						}
						cnt += 1

						if bg.bhm.CanCommit(*resp.Result.BlockHeight) {
							blockCh <- resp.Result
							bg.bhm.Commit(*resp.Result.BlockHeight)
							break
						}

						time.Sleep(queryDuration)
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
