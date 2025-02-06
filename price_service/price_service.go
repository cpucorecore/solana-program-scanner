package price_service

import "github.com/shopspring/decimal"

type PriceService interface {
	GetPrice(slot uint64, blockTime int64) decimal.Decimal
}
