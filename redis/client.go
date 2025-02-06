package redis

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"solana-program-scanner/log"
	"solana-program-scanner/types"
)

const (
	keyBase          = "smt:900:"
	keySlotPublished = keyBase + "slot_published"
	keySlot          = keyBase + "s:"
	keyMarket        = keyBase + "m:"
	keyToken         = keyBase + "t:"
)

type Client interface {
	GetSlotPublished() (uint64, bool)
	SetSlotPublished(uint64)

	GetBlockHeight(slot uint64) (int64, bool)
	SetBlockHeight(slot uint64, height int64)

	GetMarket(string) (*types.Market, bool)
	SetMarket(string, *types.Market)

	GetToken(string) (*types.Token, bool)
	SetToken(string, *types.Token)

	Close()
}

type client struct {
	ctx         context.Context
	redisClient *redis.Client
}

func New(redisAddr string) Client {
	return &client{
		ctx: context.Background(),
		redisClient: redis.NewClient(&redis.Options{
			Addr: redisAddr,
		}),
	}
}

func (c *client) GetSlotPublished() (uint64, bool) {
	v, err := c.redisClient.Get(c.ctx, keySlotPublished).Uint64()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Logger.Error("redis get err", zap.Error(err))
		}
		return 0, false
	}
	return v, true
}

func (c *client) SetSlotPublished(slotID uint64) {
	err := c.redisClient.Set(c.ctx, keySlotPublished, slotID, 0).Err()
	if err != nil {
		log.Logger.Error("redis set err", zap.Error(err))
	}
}

func (c *client) GetBlockHeight(slot uint64) (int64, bool) {
	v, err := c.redisClient.Get(c.ctx, keySlot+strconv.FormatUint(slot, 10)).Int64()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Logger.Error("redis get err", zap.Error(err))
		}
		return 0, false
	}
	return v, true
}

func (c *client) SetBlockHeight(slot uint64, height int64) {
	err := c.redisClient.Set(c.ctx, keySlot+strconv.FormatUint(slot, 10), height, 0).Err()
	if err != nil {
		log.Logger.Error("redis set block height err", zap.Error(err))
	}
}

func (c *client) GetMarket(addr string) (*types.Market, bool) {
	market := &types.Market{}
	err := c.getJSON(keyMarket+addr, market)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Logger.Error("redis get err", zap.Error(err))
		}
		return nil, false
	}
	return market, true
}

func (c *client) SetMarket(addr string, market *types.Market) {
	err := c.setJSON(keyMarket+addr, market)
	if err != nil {
		log.Logger.Error("redis set market err", zap.Error(err), zap.String("addr", addr))
	}
}

func (c *client) GetToken(addr string) (*types.Token, bool) {
	token := &types.Token{}
	err := c.getJSON(keyToken+addr, token)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Logger.Error("redis get err", zap.Error(err))
		}
		return nil, false
	}
	return token, true
}

func (c *client) SetToken(addr string, token *types.Token) {
	err := c.setJSON(keyToken+addr, token)
	if err != nil {
		log.Logger.Error("redis set token err", zap.Error(err), zap.String("addr", addr))
	}
}

func (c *client) getJSON(key string, value interface{}) error {
	bytes, err := c.redisClient.Get(c.ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, value)
}

func (c *client) setJSON(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redisClient.Set(c.ctx, key, bytes, 0).Err()
}

func (c *client) Close() {
	c.redisClient.Close()
}
