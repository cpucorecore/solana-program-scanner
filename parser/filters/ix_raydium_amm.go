package filters

import (
	"solana-program-scanner/idls/raydium_amm"
	"solana-program-scanner/types"
)

func FilterIxByMarket(market *types.Market) bool {
	if types.IsTokenSol(market.BaseMint) {
		return types.IsTokenSol(market.QuoteMint)
	} else {
		return !types.IsTokenSol(market.QuoteMint)
	}

	return false
}

func FilterIxByTransfer0Destination(destination0 string, market *types.Market) bool {
	if destination0 != market.BaseVault && destination0 != market.QuoteVault {
		return true
	}
	return false
}

func FilterIxByIxType(ixType uint8) bool {
	switch ixType {
	case raydium_amm.Instruction_SwapBaseIn,
		raydium_amm.Instruction_SwapBaseOut,
		raydium_amm.Instruction_Deposit,
		raydium_amm.Instruction_Withdraw:
		return false
	default:
		return true
	}
}
