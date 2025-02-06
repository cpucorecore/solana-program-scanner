package utils

import (
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"

	"solana-program-scanner/log"
)

func ToDecimal(base string, exp int32) (decimal.Decimal, bool) {
	value := new(big.Int)
	_, ok := value.SetString(base, 10)
	if !ok {
		log.Logger.Warn(fmt.Sprintf("ToDecimal(%s, %d) failed", base, exp))
		return decimal.Zero, false
	}
	return decimal.NewFromBigInt(value, exp), true
}
