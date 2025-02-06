package cache

import (
	"time"

	"github.com/patrickmn/go-cache"

	"solana-program-scanner/config"
	"solana-program-scanner/redis"
	"solana-program-scanner/types"
)

type TwoTierCache struct {
	memory *cache.Cache
	redis  redis.Client
}

func New(redisAddr string, conf *config.MemoryCacheConf) Cache {
	return &TwoTierCache{
		memory: cache.New(time.Hour*time.Duration(conf.ExpirationByHour), time.Hour*time.Duration(conf.CleanupIntervalByHour)),
		redis:  redis.New(redisAddr),
	}
}

func (c *TwoTierCache) GetBlockHeight(slot uint64) (int64, bool) {
	return c.redis.GetBlockHeight(slot)
}

func (c *TwoTierCache) SetBlockHeight(slot uint64, height int64) {
	c.redis.SetBlockHeight(slot, height)
}

func (c *TwoTierCache) BlockExist(slot uint64) bool {
	_, ok := c.redis.GetBlockHeight(slot)
	return ok
}

func (c *TwoTierCache) GetMarket(addr string) (*types.Market, bool) {
	if market, ok := c.memory.Get(addr); ok {
		return market.(*types.Market), true
	}

	market, ok := c.redis.GetMarket(addr)
	if !ok {
		return nil, false
	}
	c.memory.SetDefault(addr, market)

	return market, true
}

func (c *TwoTierCache) SetMarket(addr string, market *types.Market) {
	c.memory.SetDefault(addr, market)
	c.redis.SetMarket(addr, market)
}

func (c *TwoTierCache) GetToken(addr string) (*types.Token, bool) {
	if token, ok := c.memory.Get(addr); ok {
		return token.(*types.Token), true
	}

	token, ok := c.redis.GetToken(addr)
	if !ok {
		return nil, false
	}
	c.memory.SetDefault(addr, token)

	return token, true
}

func (c *TwoTierCache) SetToken(addr string, token *types.Token) {
	c.memory.SetDefault(addr, token)
	c.redis.SetToken(addr, token)
}

func (c *TwoTierCache) TokenExist(addr string) bool {
	_, ok := c.redis.GetToken(addr)
	if !ok {
		return false
	}

	_, ok = c.memory.Get(addr)
	return ok
}
