package main

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	cm := NewCacheRedisMarket(client)

	mockMarket := OrmMarket{
		Address:      "mock",
		BaseDecimal:  0,
		QuoteDecimal: 0,
		BaseMint:     "mock",
		QuoteMint:    "mock",
	}

	err := cm.Set(mockMarket.Address, &mockMarket)
	require.Equal(t, nil, err)

	get, err := cm.Get(mockMarket.Address)
	require.Equal(t, nil, err)
	require.Equal(t, mockMarket.QuoteMint, get.QuoteMint)
}
