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
