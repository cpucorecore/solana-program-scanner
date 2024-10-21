package main

import (
	"context"
	"flag"
)

func parseFlag() { // TODO
	flag.StringVar(&gCfg.Solana.RpcEndpoint, "solana-rpc-endpoint", DefaultSolanaRpcEndpoint, "solana rpc endpoint")
	flag.IntVar(&gCfg.GetterBlock.WorkerNumber, "getter-block-worker-number", DefaultGetterBlockWorkerNumber, "getter block worker number")
	flag.Uint64Var(&gCfg.GetterBlock.StartSlot, "getter-block-start-slot", DefaultGetterBlockStartSlot, "start slot")
}

func main() {
	parseFlag()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := Factory{}

	wg, fc := f.assemble().runProducts(ctx)
	f.servicePrice.Run(ctx)
	wg.Wait()
	fc.stopLog()

	f.Shutdown()
	Logger.Sync()
}
