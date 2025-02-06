package orms

import "time"

type Token struct {
	Address     string    `json:"address"`
	Name        string    `json:"name"`
	Decimal     uint64    `json:"decimal"`
	ChainId     int       `json:"chain_id"`
	TotalSupply string    `json:"total_supply"`
	Symbol      string    `json:"symbol"`
	CreatedAt   time.Time `xorm:"created"`
	Creator     string    `json:"creator"`
	Website     string    `json:"website"`
}

func (t *Token) TableName() string {
	return "token"
}
