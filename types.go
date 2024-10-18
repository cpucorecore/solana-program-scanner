package main

import "time"

type OrmTx struct {
	TxHash        string
	Event         int8
	Token0Amount  string
	Token1Amount  string
	Maker         string
	Token0Address string
	Token1Address string
	AmountUsd     float64
	PriceUsd      float64
	Block         int64
	BlockAt       time.Time
	CreatedAt     time.Time `xorm:"created"`
	Index         int
}

func (ot *OrmTx) TableName() string {
	return "tx"
}

type Market struct {
	Address      string `json:"address"`
	BaseDecimal  uint64 `json:"base_decimal"`
	QuoteDecimal uint64 `json:"quote_decimal"`
	BaseMint     string `json:"base_mint"`
	QuoteMint    string `json:"quote_mint"`
}

func (ot *Market) OrmMarket() string {
	return "market"
}
