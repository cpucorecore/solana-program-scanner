package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"solana-program-scanner/log"
	cache2 "solana-program-scanner/redis/types"
)

const (
	KeyToken = "smt:token:900"
	KeyPair  = "smt:pair:900"
)

type Cache2 interface {
	SetToken2(token *cache2.Token) error
	GetToken2(tokenAddr string) (*cache2.Token, error)
	DelToken2(tokenAddr string) error
	SetPair2(pair *cache2.Pair) error
	GetPair2(addr string) (*cache2.Pair, bool)
	DelPair2(addr string) error
	Publish(pair *cache2.Pair) error
	Close()
}

func NewCache2(redisAddr string) Cache2 {
	return &client{
		ctx: context.Background(),
		redisClient: redis.NewClient(&redis.Options{
			Addr: redisAddr,
		}),
	}
}

func (c *client) SetToken2(token *cache2.Token) error {
	tokenBytes, _ := json.Marshal(token)
	_, err := c.redisClient.HSet(c.ctx, KeyToken, token.Address, string(tokenBytes)).Result()
	return err
}

func (c *client) GetToken2(tokenAddr string) (*cache2.Token, error) {
	result, err := c.redisClient.HGet(c.ctx, KeyToken, tokenAddr).Result()
	if err != nil {
		return nil, err
	}

	var token cache2.Token
	err = json.Unmarshal([]byte(result), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (c *client) DelToken2(addr string) error {
	return c.redisClient.HDel(c.ctx, KeyToken, addr).Err()
}

func (c *client) SetPair2(pair *cache2.Pair) error {
	pairBytes, _ := json.Marshal(pair)
	_, err := c.redisClient.HSet(c.ctx, KeyPair, pair.Address, string(pairBytes)).Result()
	return err
}

func (c *client) GetPair2(addr string) (*cache2.Pair, bool) {
	result, err := c.redisClient.HGet(c.ctx, KeyPair, addr).Result()
	if err != nil {
		return nil, false
	}

	var pair cache2.Pair
	err = json.Unmarshal([]byte(result), &pair)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("redis json.Unmarshal[%s] err:%v", result, err))
		return nil, false
	}

	return &pair, true
}

func (c *client) DelPair2(addr string) error {
	return c.redisClient.HDel(c.ctx, KeyPair, addr).Err()
}

func (c *client) Publish(pair *cache2.Pair) error {
	pairBytes, err := json.Marshal(pair)
	if err != nil {
		return err
	}

	_, err = c.redisClient.Publish(c.ctx, "smt:subscribe:900:pool:"+pair.Address, string(pairBytes)).Result()
	if err != nil {
		return err
	}

	return nil
}
