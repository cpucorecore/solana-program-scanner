package types

import "github.com/shopspring/decimal"

type VaultInfo struct {
	Address string
	Mint    string
	Symbol  string
	Balance decimal.Decimal
}

type MarketVaultInfo struct {
	Address string
	Token0  *VaultInfo
	Token1  *VaultInfo
}
