package parser

import (
	"encoding/json"
	"fmt"
	"strconv"

	sg "github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"

	"solana-program-scanner/cache"
	"solana-program-scanner/getter/market_getter"
	"solana-program-scanner/idls/raydium_amm"
	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/parser/filters"
	"solana-program-scanner/types"
	"solana-program-scanner/utils"
)

type RaydiumAmmParser struct {
	marketCache  cache.MarketCache
	marketGetter market_getter.Getter
}

func NewRaydiumAmmParser(
	marketCache cache.MarketCache,
	marketGetter market_getter.Getter,
) *RaydiumAmmParser {
	return &RaydiumAmmParser{
		marketCache:  marketCache,
		marketGetter: marketGetter,
	}
}

func (p *RaydiumAmmParser) getMarket(addr string) (market *types.Market, ok bool) {
	market, ok = p.marketCache.GetMarket(addr)
	if ok {
		return market, true
	}

	market = p.marketGetter.MustGetMarket(addr)
	p.marketCache.SetMarket(addr, market)

	return market, false
}

type RaydiumInstructionParsed struct {
	Signer         string
	MarketAddress  string
	Index          int
	TransferDetail *TransferDetail
}

func logInstruction(txHash string, errMsg string, ix *RaydiumInstructionWrap) {
	bytes, _ := json.Marshal(ix)
	log.Logger.Error("raydium_amm DecodeInstruction err",
		zap.String("tx_hash", txHash),
		zap.String("err_msg", errMsg),
		zap.String("instruction", string(bytes)),
	)
}

func (p *RaydiumAmmParser) parseRaydiumInstructions(
	txHash string,
	accountBook map[string]TxAccount,
	raydiumInstructions *RaydiumInstructions,
	price decimal.Decimal,
) []*RaydiumInstructionParsed {
	monitor.IxSrcCounter.WithLabelValues("direct").Add(float64(len(raydiumInstructions.Direct)))
	monitor.IxSrcCounter.WithLabelValues("inner").Add(float64(len(raydiumInstructions.Inner)))

	ixs := append(raydiumInstructions.Direct, raydiumInstructions.Inner...)
	ixsParsed := make([]*RaydiumInstructionParsed, 0, len(ixs))

	for _, ix := range ixs {
		accounts := createInstructionAccountMetas(accountBook, ix.Raydium.Accounts)
		ixDecoded, err := raydium_amm.DecodeInstruction(accounts, ix.Raydium.DataBytes)
		if err != nil {
			logInstruction(txHash, err.Error(), ix)
			continue
		}

		ixTypeID := ixDecoded.TypeID.Uint8()
		monitor.IxTypeCounter.WithLabelValues(strconv.Itoa(int(ixTypeID))).Inc()

		if filters.FilterIxByIxType(ixTypeID) {
			continue
		}

		signer, ok := getSigner(ixTypeID, ix.Raydium.Accounts)
		if !ok {
			logInstruction(txHash, "getSigner failed", ix)
			continue
		}

		market, _ := p.getMarket(ix.Raydium.Accounts[1])
		if filters.FilterIxByMarket(market) {
			continue
		}

		if ixDecoded.TypeID.Uint8() != raydium_amm.Instruction_Withdraw && filters.FilterIxByTransfer0Destination(ix.Transfers[0].Parsed.Info.Destination, market) {
			logInstruction(txHash, "FilterIxByTransfer0Destination", ix)
			continue
		}

		amt0 := ix.Transfers[0].Parsed.Info.Amount
		amt1 := ix.Transfers[1].Parsed.Info.Amount
		dst0 := ix.Transfers[0].Parsed.Info.Destination

		var transferDetail *TransferDetail
		if ixDecoded.TypeID.Uint8() == raydium_amm.Instruction_SwapBaseIn || ixDecoded.TypeID.Uint8() == raydium_amm.Instruction_SwapBaseOut {
			transferDetail = ParseEventBuyOrSell(amt0, amt1, dst0, market, price)
		} else if ixDecoded.TypeID.Uint8() == raydium_amm.Instruction_Deposit {
			transferDetail = parseEventPoolAdd(amt0, amt1, dst0, market)
		} else if ixDecoded.TypeID.Uint8() == raydium_amm.Instruction_Withdraw {
			src0 := ix.Transfers[0].Parsed.Info.Source
			transferDetail = ParseEventPoolRemove(amt0, amt1, src0, market)
		} else {
			continue
		}

		ixsParsed = append(ixsParsed, &RaydiumInstructionParsed{
			Signer:         signer,
			MarketAddress:  market.Address,
			Index:          ix.Raydium.Index,
			TransferDetail: transferDetail,
		})
	}

	return ixsParsed
}

func getSigner(ixType uint8, accounts []string) (string, bool) {
	switch ixType {
	case raydium_amm.Instruction_SwapBaseIn, raydium_amm.Instruction_SwapBaseOut:
		if len(accounts) == 17 {
			return accounts[16], true
		}

		if len(accounts) == 18 {
			return accounts[17], true
		}

		return "", false

	case raydium_amm.Instruction_Deposit:
		if len(accounts) != 14 {
			return "", false
		}
		return accounts[12], true

	case raydium_amm.Instruction_Withdraw:
		if len(accounts) != 22 {
			return "", false
		}
		return accounts[18], true

	default:
		return "", false
	}
}

type TransferDetail struct {
	Event         types.Event
	Token0Address string
	Token0Vault   string
	Token0Amount  decimal.Decimal
	Token1Address string
	Token1Vault   string
	Token1Amount  decimal.Decimal
	AmountUsd     decimal.Decimal
	PriceUsd      decimal.Decimal
}

/*
ParseEventBuyOrSell

Token1 always sol
[Instruction_SwapBaseIn] --> [event=Buy] https://solscan.io/tx/3Qw8rKGpxUsHxU9p1WpqbxERqmePFsd3k89CyELADdM6UXSLtzyczBHqFLjHV4VZMnccq1jX8ujd3QNraxzwvov7
[Instruction_SwapBaseIn] --> [event=Sell] https://solscan.io/tx/5EKhSGjhPo3UXDbSBEXhazqXE1bRg6eFp9oZgA8jUgmnBM5sSeXBPqnxQCXSkTJEMvYVWx2t42TrwQqMpdX64kqE
[Instruction_SwapBaseOut] --> [event=Buy] https://solscan.io/tx/FsEoYDxhekaRxjbEr3qJfsER32sp8yNEm6adupeDnXFE6V5y2jqhLpbjHSWdxVqf7ZmxTbZhYMYUaPbnAT1K3TG
[Instruction_SwapBaseOut] --> [event=Sell] https://solscan.io/tx/5Bce5HWPAGgZPto9nyhCxBBpuQE8apfoqsAjvdB1mScLmRQk5Towio2TMdxdJx3YaGWbnivTm6t3ooMiPwZtjimH
*/
func ParseEventBuyOrSell(
	amount0 string,
	amount1 string,
	destination0 string,
	market *types.Market,
	solPriceUSD decimal.Decimal,
) *TransferDetail {
	var d = &TransferDetail{}

	if destination0 == market.BaseVault {
		if types.IsTokenSol(market.BaseMint) {
			d.Event = types.Buy
			d.Token0Address = market.QuoteMint
			d.Token0Vault = market.QuoteVault
			d.Token0Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
			d.Token1Address = market.BaseMint
			d.Token1Vault = market.BaseVault
			d.Token1Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
		} else {
			d.Event = types.Sell
			d.Token0Address = market.BaseMint
			d.Token0Vault = market.BaseVault
			d.Token0Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
			d.Token1Address = market.QuoteMint
			d.Token1Vault = market.QuoteVault
			d.Token1Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
		}
	} else if destination0 == market.QuoteVault {
		if types.IsTokenSol(market.QuoteMint) {
			d.Event = types.Buy
			d.Token0Address = market.BaseMint
			d.Token0Vault = market.BaseVault
			d.Token0Amount, _ = utils.ToDecimal(amount1, -market.BaseDecimal)
			d.Token1Address = market.QuoteMint
			d.Token1Vault = market.QuoteVault
			d.Token1Amount, _ = utils.ToDecimal(amount0, -market.QuoteDecimal)
		} else {
			d.Event = types.Sell
			d.Token0Address = market.QuoteMint
			d.Token0Vault = market.QuoteVault
			d.Token0Amount, _ = utils.ToDecimal(amount0, -market.QuoteDecimal)
			d.Token1Address = market.BaseMint
			d.Token1Vault = market.BaseVault
			d.Token1Amount, _ = utils.ToDecimal(amount1, -market.BaseDecimal)
		}
	} else {
		log.Logger.Fatal(fmt.Sprintf("wrong destination0")) // TODO check
	}

	d.PriceUsd, d.AmountUsd = CalcPriceAndAmount(solPriceUSD, d.Token1Amount, d.Token0Amount)

	return d
}

func parseEventPoolAdd(
	amount0 string,
	amount1 string,
	destination0 string,
	market *types.Market,
) *TransferDetail {
	var d = &TransferDetail{}
	d.Event = types.Add
	d.AmountUsd = decimal.Zero
	d.PriceUsd = decimal.Zero

	if destination0 == market.BaseVault {
		if types.IsTokenSol(market.BaseMint) {
			d.Token0Address = market.QuoteMint
			d.Token0Vault = market.QuoteVault
			d.Token0Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
			d.Token1Address = market.BaseMint
			d.Token1Vault = market.BaseVault
			d.Token1Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
		} else {
			d.Token0Address = market.BaseMint
			d.Token0Vault = market.BaseVault
			d.Token0Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
			d.Token1Address = market.QuoteMint
			d.Token1Vault = market.QuoteVault
			d.Token1Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
		}
	} else if destination0 == market.QuoteVault {
		if types.IsTokenSol(market.QuoteMint) {
			d.Token0Address = market.BaseMint
			d.Token0Vault = market.BaseVault
			d.Token0Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
			d.Token1Address = market.QuoteMint
			d.Token1Vault = market.QuoteVault
			d.Token1Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
		} else {
			d.Token0Address = market.QuoteMint
			d.Token0Vault = market.QuoteVault
			d.Token0Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
			d.Token1Address = market.BaseMint
			d.Token1Vault = market.BaseVault
			d.Token1Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
		}
	} else {
		log.Logger.Fatal(fmt.Sprintf("wrong destination0"))
	}

	return d
}

func ParseEventPoolRemove(
	amount0 string,
	amount1 string,
	source0 string,
	market *types.Market,
) *TransferDetail {
	var d = &TransferDetail{}
	d.Event = types.Remove
	d.AmountUsd = decimal.Zero
	d.PriceUsd = decimal.Zero

	if source0 == market.BaseVault {
		if types.IsTokenSol(market.BaseMint) {
			d.Token0Address = market.QuoteMint
			d.Token0Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
			d.Token1Address = market.BaseMint
			d.Token1Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
		} else {
			d.Token0Address = market.BaseMint
			d.Token0Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
			d.Token1Address = market.QuoteMint
			d.Token1Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
		}
	} else if source0 == market.QuoteVault {
		if types.IsTokenSol(market.QuoteMint) {
			d.Token0Address = market.BaseMint
			d.Token0Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
			d.Token1Address = market.QuoteMint
			d.Token1Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
		} else {
			d.Token0Address = market.QuoteMint
			d.Token0Amount, _ = utils.ToDecimal(amount1, -market.QuoteDecimal)
			d.Token1Address = market.BaseMint
			d.Token1Amount, _ = utils.ToDecimal(amount0, -market.BaseDecimal)
		}
	} else {
		log.Logger.Fatal(fmt.Sprintf("wrong source0"))
	}

	return d
}

func createInstructionAccountMetas(transactionAccountBook map[string]TxAccount, instructionAccounts []string) []*sg.AccountMeta {
	var accountMetas []*sg.AccountMeta
	for _, account := range instructionAccounts {
		accountMetas = append(accountMetas, &sg.AccountMeta{
			PublicKey:  sg.MustPublicKeyFromBase58(account),
			IsWritable: transactionAccountBook[account].AccountKey.Writable,
			IsSigner:   transactionAccountBook[account].AccountKey.Signer,
		})
	}
	return accountMetas
}
