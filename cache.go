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
	local          map[string]*OrmMarket
	redisKeyPrefix string
}

var _ Cache[string, *OrmMarket] = &CacheMarketRedis{}

const (
	RedisKeyPrefixMarket = "m:"
)

func NewCacheRedisMarket(redisCli *redis.Client) Cache[string, *OrmMarket] {
	return &CacheMarketRedis{
		cli:            redisCli,
		ctx:            context.Background(),
		local:          make(map[string]*OrmMarket),
		redisKeyPrefix: RedisKeyPrefixMarket,
	}
}

func (c *CacheMarketRedis) Get(k string) (*OrmMarket, error) {
	if market, ok := c.local[k]; ok {
		return market, nil
	}

	data, err := c.cli.Get(c.ctx, RedisKeyPrefixMarket+k).Result()
	if err != nil {
		return nil, nil
	}

	var m OrmMarket
	err = json.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil, err
	}

	c.local[k] = &m

	return &m, nil
}

func (c *CacheMarketRedis) Set(k string, v *OrmMarket) error {
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
