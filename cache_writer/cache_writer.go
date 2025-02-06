package cache_writer

import (
	"fmt"

	"go.uber.org/zap"

	"solana-program-scanner/log"
	"solana-program-scanner/redis"
	"solana-program-scanner/redis/types"
)

type CacheWriter struct {
	cache2 redis.Cache2
}

func NewCacheWriter(cache2 redis.Cache2) *CacheWriter {
	return &CacheWriter{cache2}
}

func (w *CacheWriter) WritePair(pair *types.Pair) {
	err := w.cache2.SetPair2(pair)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("SetPair2 err:%s", err.Error()), zap.Any("pair", pair))
	}
}

func (w *CacheWriter) GetToken(tokenAddr string) (*types.Token, error) {
	return w.cache2.GetToken2(tokenAddr)
}

func (w *CacheWriter) WriteToken(token *types.Token) {
	err := w.cache2.SetToken2(token)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("SetToken2 err:%s", err.Error()), zap.Any("token", token))
	}
}
