package main

import (
	"context"
	"flag"
	"time"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := Factory{}
	wg, fc := f.assemble().runProducts(ctx)
	wg.Wait()
	fc.stopLog()
}
