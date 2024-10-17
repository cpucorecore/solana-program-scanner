package main

import (
	"context"
	"flag"
	"sync"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"
	"go.mongodb.org/mongo-driver/v2/bson"

	"solana-program-scanner/block_height_manager"
)

func parseFlag() {
	flag.StringVar(&conf.Solana.RpcEndpoint, "solana-rpc-endpoint", DefaultSolanaRpcEndpoint, "solana rpc endpoint")

	var ReqInterval int
	flag.IntVar(&ReqInterval, "req-interval", DefaultSolanaRpcReqIntervalMillisecond, "rpc request interval in millisecond")
	conf.Solana.RpcReqInterval = time.Duration(ReqInterval) * time.Millisecond

	flag.IntVar(&conf.Solana.BlockGetterWorkerNumber, "block-getter-worker-number", DefaultSolanaBlockGetterWorkerNumber, "block getter worker number")
	flag.Uint64Var(&conf.Solana.StartSlot, "start-slot", DefaultSolanaStartSlot, "start slot")
}

func main() {
	parseFlag()

	ctx := context.Background()
	var wg sync.WaitGroup

	ixCh := make(chan bson.M, 100)
	mongo := NewMongoAttendant(ixCh)

	mongo.startServe(ctx, &wg)

	blockCh := make(chan *rpc.GetBlock, 100)
	blockProcessor := NewBlockProcessorAdmin(blockCh, ixCh)
	wg.Add(1)
	go blockProcessor.run(ctx, &wg)

	fc := NewFlowController(5, 5, time.Second, time.Second)
	go fc.startLog(time.Second * 5)

	bhm := block_height_manager.NewBlockHeightManager()
	blockGetter := NewBlockGetter(conf.Solana.BlockGetterWorkerNumber, bhm, fc)
	startBlockHeight := blockGetter.getBlockHeightBySlot(conf.Solana.StartSlot)
	bhm.Init(startBlockHeight - 1)

	taskCh := make(chan uint64, 1000)
	wg.Add(1)
	go blockGetter.run(ctx, &wg, taskCh, blockCh)

	taskDispatcher := NewBlockTaskDispatcher(fc)
	wg.Add(1)
	go taskDispatcher.keepDispatchingTask(ctx, &wg, conf.Solana.StartSlot, 5, taskCh)

	Logger.Info("wait all goroutine done")
	wg.Wait()
	fc.stopLog()
	Logger.Info("all goroutine done")
}
