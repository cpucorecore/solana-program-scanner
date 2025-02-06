package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	KeyPrice = "smt:oracle:price:900:So11111111111111111111111111111111111111112"
)

type CachePrice interface {
	SetPrice(price string) error
}

func NewCachePrice(redisAddr string) CachePrice {
	return &client{
		ctx: context.Background(),
		redisClient: redis.NewClient(&redis.Options{
			Addr: redisAddr,
		}),
	}
}

type TimePrice struct {
	UnixTimestamp int64
	Price         string
}

var timePrice TimePrice

func (c *client) SetPrice(price string) error {
	timePrice.UnixTimestamp = time.Now().Unix()
	timePrice.Price = price
	timePriceBytes, _ := json.Marshal(timePrice)
	_, err := c.redisClient.Set(c.ctx, KeyPrice, string(timePriceBytes), 0).Result()
	return err
}
