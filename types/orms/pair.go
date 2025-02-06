package orms

import "time"

type Pair struct {
	Address   string
	Name      string
	Token0    string
	Token1    string
	Reserve0  string
	Reserve1  string
	ChainId   int
	CreatedAt time.Time `xorm:"created"`
}

func (p *Pair) TableName() string {
	return "pair"
}
