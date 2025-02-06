package price_service

import (
	"github.com/shopspring/decimal"
	"sync"
)

type MutexPrice struct {
	rwMutex sync.RWMutex
	price   decimal.Decimal
}

func NewMutexPrice() *MutexPrice {
	return &MutexPrice{}
}

func (p *MutexPrice) Get() decimal.Decimal {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	return p.price
}

func (p *MutexPrice) Set(price decimal.Decimal) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.price = price
}
