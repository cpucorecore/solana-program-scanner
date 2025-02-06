package price_service

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"solana-program-scanner/config"
	"solana-program-scanner/log"
)

func TestPriceServiceHermes_GetPrice(t *testing.T) {
	log.InitLoggerForTest()
	ps := &PriceServiceHermes{
		conf: &config.PriceServiceConf{
			BaseUrl: "https://hermes.pyth.network/v2/updates/price/",
			Id:      "ids%5B%5D=0xef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d",
		},
	}

	tests := []struct {
		name      string
		slot      uint64
		blockTime int64
		setup     func()
		want      decimal.Decimal
	}{
		{
			name:      "recent block time - return current price",
			slot:      1000,
			blockTime: time.Now().Unix(),
			setup: func() {
				ps.currentPrice.Set(decimal.NewFromFloat(100.5))
			},
			want: decimal.NewFromFloat(100.5),
		},
		{
			name:      "old block time with slot multiple of 200 - query historical price",
			slot:      308685200,
			blockTime: 1734703444,
			setup: func() {
				ps.historyPrice.Set(decimal.NewFromFloat(186))
			},
			want: decimal.NewFromFloat(186.69635719),
		},
		{
			name:      "old block time with slot multiple of 200 - query historical price",
			slot:      308685201,
			blockTime: 1734703444,
			setup: func() {
				ps.historyPrice.Set(decimal.NewFromFloat(186))
			},
			want: decimal.NewFromFloat(186),
		},
		{
			name:      "old block time with non-200 slot - return cached history price",
			slot:      2001,
			blockTime: time.Now().Unix() - 120,
			setup: func() {
				ps.historyPrice.Set(decimal.NewFromFloat(95.5))
			},
			want: decimal.NewFromFloat(95.5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := ps.GetPrice(tt.slot, tt.blockTime)
			assert.Equal(t, tt.want, got)
		})
	}
}
