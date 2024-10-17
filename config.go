package main

import "time"

var conf = &Config{}

const (
	DefaultSolanaRpcEndpoint               = "https://api.mainnet-beta.solana.com"
	DefaultSolanaRpcReqIntervalMillisecond = 10_000 // 10 seconds
	DefaultSolanaBlockGetterWorkerNumber   = 3
	DefaultSolanaStartSlot                 = uint64(295503380)
)

type SolanaConf struct {
	RpcEndpoint             string
	RpcReqInterval          time.Duration
	RpcReqErrDelay          string
	BlockGetterWorkerNumber int
	StartSlot               uint64
}

type Config struct {
	Solana SolanaConf
}
