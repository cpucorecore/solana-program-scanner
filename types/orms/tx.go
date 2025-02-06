package orms

import (
	"time"

	"github.com/shopspring/decimal"
	"solana-program-scanner/types"
)

type Tx struct {
	TxHash        string
	Event         types.Event
	Token0Amount  decimal.Decimal
	Token1Amount  decimal.Decimal
	Maker         string
	Token0Address string
	Token1Address string
	AmountUsd     decimal.Decimal
	PriceUsd      decimal.Decimal
	Block         uint64
	BlockAt       time.Time
	CreatedAt     time.Time `xorm:"created"`
	BlockIndex    int
	TxIndex       int
}

func (t *Tx) TableName() string {
	return "tx"
}
