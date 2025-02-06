package types

import "github.com/blocto/solana-go-sdk/rpc"

type BlockWithSlot struct {
	Slot  uint64
	Block *rpc.GetBlock
}
