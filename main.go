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

func main() {
	flag.StringVar(&conf.Solana.RpcEndpoint, "solana-rpc-endpoint", DefaultSolanaRpcEndpoint, "solana rpc endpoint")

	var ReqInterval int
	flag.IntVar(&ReqInterval, "req-interval", DefaultSolanaRpcReqIntervalMillisecond, "rpc request interval in millisecond")
	conf.Solana.RpcReqInterval = time.Duration(ReqInterval) * time.Millisecond

	flag.IntVar(&conf.Solana.BlockGetterWorkerNumber, "block-getter-worker-number", DefaultSolanaBlockGetterWorkerNumber, "block getter worker number")
	flag.Uint64Var(&conf.Solana.StartSlot, "start-slot", DefaultSolanaStartSlot, "start slot")

	taskCh := make(chan uint64, 1000)

	blockCh := make(chan *rpc.GetBlock, 100)

	txRawCh := make(chan string, 100)
	ixRawCh := make(chan string, 100)
	ixIndexCh := make(chan bson.M, 100)
	ixCh := make(chan bson.M, 100)

	ctx := context.Background()
	var wg sync.WaitGroup

	mongo := NewMongoAttendant(txRawCh, ixRawCh, ixIndexCh, ixCh)
	mongo.startServe(ctx, &wg)

	blockProcessor := NewBlockProcessorAdmin(blockCh, txRawCh, ixRawCh, ixIndexCh, ixCh)
	wg.Add(1)
	go blockProcessor.run(ctx, &wg)

	fc := NewFlowController(5, 5, time.Second, time.Second)
	go fc.startLog(time.Second * 1)
	bhm := block_height_manager.NewBlockHeightManager()
	blockGetter := NewBlockGetter(conf.Solana.BlockGetterWorkerNumber, bhm, fc)
	startBlockHeight := blockGetter.getBlockHeightBySlot(conf.Solana.StartSlot)
	bhm.Init(startBlockHeight - 1)

	wg.Add(1)
	go blockGetter.run(ctx, &wg, taskCh, blockCh)

	taskDispatch := NewBlockTaskDispatch(fc)
	wg.Add(1)
	go taskDispatch.keepDispatchingTask(ctx, &wg, conf.Solana.StartSlot, 20, taskCh)

	Logger.Info("wait all goroutine done")
	wg.Wait()
	fc.stopLog()
	Logger.Info("all goroutine done")
}
