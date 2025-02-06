package cache

import (
	"solana-program-scanner/types"
)

type BlockCache interface {
	GetBlockHeight(slot uint64) (int64, bool)
	SetBlockHeight(slot uint64, height int64)
	BlockExist(slot uint64) bool
}

type MarketCache interface {
	GetMarket(string) (*types.Market, bool)
	SetMarket(string, *types.Market)
}

type TokenCache interface {
	GetToken(string) (*types.Token, bool)
	SetToken(string, *types.Token)
	TokenExist(k string) bool
}

type Cache interface {
	BlockCache
	MarketCache
	TokenCache
}
