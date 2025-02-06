package parser

import (
	"math/big"

	"github.com/shopspring/decimal"

	"solana-program-scanner/log"
)

func MustToDecimal(base string, exp int32) decimal.Decimal {
	value := new(big.Int)
	_, ok := value.SetString(base, 10)
	if !ok {
		log.Logger.Fatal("big.Int SetString failed") // TODO non Fatal
	}
	return decimal.NewFromBigInt(value, exp)
}

func CalcPriceAndAmount(
	solPriceUSD decimal.Decimal,
	solAmount decimal.Decimal,
	nonSolAmount decimal.Decimal,
) (priceUSD decimal.Decimal, amountUSD decimal.Decimal) {
	amountUSD = solPriceUSD.Mul(solAmount)
	priceUSD = amountUSD.Div(nonSolAmount)
	return
}
