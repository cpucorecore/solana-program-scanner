package block_getter

import (
	"solana-program-scanner/config"
	"solana-program-scanner/getter/rate_limiter"
	"solana-program-scanner/msg_broker"
	"solana-program-scanner/parser"
	"solana-program-scanner/sequencer"
)

type BlockGetterManagerConfig struct {
	Conf            *config.BlockGetterManagerConf
	RPCEndpointHTTP string
	StartSlot       uint64
	BlockParser     *parser.BlockParser
	BlockSequencer  sequencer.BlockSequencer
	RpcRateLimiter  rate_limiter.RpcRateLimiter
	SlotBroker      msg_broker.MsgConsumer[uint64]
	Getters         []*BlockGetter
}
