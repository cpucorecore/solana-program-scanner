package orms

import "time"

type BlockErr struct {
	Slot      uint64
	ErrCode   int
	CreatedAt time.Time `xorm:"created"`
}

func (o *BlockErr) TableName() string {
	return "block_err"
}
