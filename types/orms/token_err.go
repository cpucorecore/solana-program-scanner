package orms

import "time"

type TokenErr struct {
	Address   string `json:"address"`
	Reason    int
	CreatedAt time.Time `xorm:"created"`
}

func (e *TokenErr) TableName() string {
	return "token_err"
}
