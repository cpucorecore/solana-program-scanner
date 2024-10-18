package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
)

type Cache[K any, V any] interface {
	Get(k K) (V, error)
	Set(k K, v V) error
}

type CacheMarketRedis struct {
	cli            *redis.Client
	ctx            context.Context
	local          map[string]*Market
	redisKeyPrefix string
}

var _ Cache[string, *Market] = &CacheMarketRedis{}

const (
	RedisKeyPrefixMarket    = "m:"
	RedisKeyPrefixMarketLen = len(RedisKeyPrefixMarket)
)

func NewCacheRedisMarket(redisCli *redis.Client) Cache[string, *Market] {
	return &CacheMarketRedis{
		cli:            redisCli,
		ctx:            context.Background(),
		local:          make(map[string]*Market),
		redisKeyPrefix: RedisKeyPrefixMarket,
	}
}

func (c *CacheMarketRedis) Get(k string) (*Market, error) {
	if market, ok := c.local[k]; ok {
		return market, nil
	}

	data, err := c.cli.Get(c.ctx, RedisKeyPrefixMarket+k).Result()
	if err != nil {
		return nil, nil
	}

	var m Market
	err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}

	c.local[k] = &m

	return &m, nil
}

func (c *CacheMarketRedis) Set(k string, v *Market) error {
	c.local[k] = v

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = c.cli.Set(c.ctx, RedisKeyPrefixMarket+k, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
