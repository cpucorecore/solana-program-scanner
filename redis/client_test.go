package redis_test

import (
	"github.com/test-go/testify/require"
	"solana-program-scanner/redis"
	"solana-program-scanner/types"
	"testing"
)

func TestClient_GetSlotPublished(t *testing.T) {
	redisClient := redis.New("localhost:6379")
	slotPublished, err := redisClient.GetSlotPublished()
	require.Error(t, err)
	require.Equal(t, uint64(0), slotPublished)
}

func TestClient_GetMarket(t *testing.T) {
	redisClient := redis.New("localhost:6379")
	market := &types.Market{
		ChainId:      1,
		Address:      "test",
		BaseDecimal:  1,
		QuoteDecimal: 1,
		BaseVault:    "test",
		QuoteVault:   "test",
		BaseMint:     "test",
		QuoteMint:    "test",
	}

	err := redisClient.SetMarket("test", market)
	require.NoError(t, err)

	getMarket, err := redisClient.GetMarket("test")
	require.NoError(t, err)
	require.Equal(t, market, getMarket)
}

func TestClient_GetBlockHeight(t *testing.T) {
	redisClient := redis.New("localhost:6379")
	height := redisClient.GetBlockHeight(1)
	require.Equal(t, int64(-1), height)
}
