package main

import (
	"fmt"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"
	sg "github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"

	"solana-program-scanner/idls/raydium_amm"
)

type ParserTxRaydiumAmm struct {
	txCh         chan *OrmTx
	marketCh     chan *OrmMarket
	getterMarket GetterMarket
	cacheMarket  Cache[string, *OrmMarket]
}

func NewParserTxRaydiumAmm(
	txCh chan *OrmTx,
	marketCh chan *OrmMarket,
	getterMarket GetterMarket,
	cacheMarket Cache[string, *OrmMarket],
) *ParserTxRaydiumAmm {
	return &ParserTxRaydiumAmm{
		txCh:         txCh,
		marketCh:     marketCh,
		getterMarket: getterMarket,
		cacheMarket:  cacheMarket,
	}
}

func (pt *ParserTxRaydiumAmm) getMarket(marketAddress string) *OrmMarket {
	market, err := pt.cacheMarket.Get(marketAddress)
	if err != nil {
		Logger.Fatal(fmt.Sprintf("market cache get err:%v", err)) // TODO check
	}

	if market == nil {
		market, err = pt.getterMarket.getMarket(marketAddress)
		if err != nil {
			Logger.Fatal(fmt.Sprintf("get market err:%v", err))
		}

		pt.marketCh <- market

		err = pt.cacheMarket.Set(market.Address, market)
		if err != nil {
			Logger.Warn(fmt.Sprintf("market cache set err:%v", err))
		}
	}

	return market
}

func (pt *ParserTxRaydiumAmm) ParseIxSwapBaseIn(
	blockHeight int64,
	blockTime int64,
	txIndex uint64,
	txHash string,
	ix *raydium_amm.Instruction,
	ixesInnerParsed []rpc.InstructionInnerParsed,
) {
	swapBaseIn, ok := ix.Impl.(*raydium_amm.SwapBaseIn)
	if !ok {
		Logger.Fatal("type assertion (*raydium_amm.SwapBaseIn) failed")
	}

	market := pt.getMarket(swapBaseIn.AccountMetaSlice[1].PublicKey.String())

	var signer string
	if len(swapBaseIn.AccountMetaSlice) == 18 {
		signer = swapBaseIn.AccountMetaSlice[17].PublicKey.String()
	} else if len(swapBaseIn.AccountMetaSlice) == 17 {
		signer = swapBaseIn.AccountMetaSlice[16].PublicKey.String()
	} else {
		Logger.Fatal(fmt.Sprintf("wrong account number:%d, txHash:%s", len(swapBaseIn.AccountMetaSlice), txHash))
	}

	ormTx := &OrmTx{
		TxHash: txHash,
		Event:  0, // TODO
		//Token0Amount:  strconv.FormatUint(*ixSwapBaseIn.AmountIn, 10),
		//Token1Amount:  strconv.FormatUint(*ixSwapBaseIn.MinimumAmountOut, 10),
		Token0Amount:  ixesInnerParsed[0].Parsed.Info.Amount,
		Token1Amount:  ixesInnerParsed[1].Parsed.Info.Amount,
		Maker:         signer,
		Token0Address: market.BaseMint,
		Token1Address: market.QuoteMint,
		Block:         blockHeight,
		BlockAt:       time.Unix(blockTime, 0),
		Index:         txIndex,
	}

	pt.txCh <- ormTx
}

func (pt *ParserTxRaydiumAmm) ParseIxSwapBaseOut(
	blockHeight int64,
	blockTime int64,
	txIndex uint64,
	txHash string,
	ix *raydium_amm.Instruction,
	ixesInnerParsed []rpc.InstructionInnerParsed,
) {
	swapBaseOut, ok := ix.Impl.(*raydium_amm.SwapBaseOut)
	if !ok {
		Logger.Fatal("type assertion (*raydium_amm.SwapBaseOut) failed")
	}

	var signer string
	if len(swapBaseOut.AccountMetaSlice) == 18 {
		signer = swapBaseOut.AccountMetaSlice[17].PublicKey.String()
	} else if len(swapBaseOut.AccountMetaSlice) == 17 {
		signer = swapBaseOut.AccountMetaSlice[16].PublicKey.String()
	} else {
		Logger.Fatal(fmt.Sprintf("wrong account number:%d, txHash:%s", len(swapBaseOut.AccountMetaSlice), txHash))
	}

	market := pt.getMarket(swapBaseOut.AccountMetaSlice[1].PublicKey.String())

	ormTx := &OrmTx{
		TxHash: txHash,
		Event:  0, // TODO
		//Token0Amount:  strconv.FormatUint(*swapBaseOut.MaxAmountIn, 10),
		//Token1Amount:  strconv.FormatUint(*swapBaseOut.AmountOut, 10),
		Token0Amount:  ixesInnerParsed[0].Parsed.Info.Amount,
		Token1Amount:  ixesInnerParsed[1].Parsed.Info.Amount,
		Maker:         signer,
		Token0Address: market.BaseMint,
		Token1Address: market.QuoteMint,
		Block:         blockHeight,
		BlockAt:       time.Unix(blockTime, 0),
		Index:         txIndex,
	}

	pt.txCh <- ormTx
}

func (pt *ParserTxRaydiumAmm) ParseTx(
	blockHeight int64,
	blockTime int64,
	txIndex uint64,
	ixIndex uint64,
	tx *rpc.GetBlockTransaction,
	ixf *rpc.InstructionFull,
) {
	accountBook := make(map[string]rpc.AccountKey, len(tx.Transaction.Message.AccountKeys))
	for _, accountKey := range tx.Transaction.Message.AccountKeys {
		accountBook[accountKey.Pubkey] = accountKey
	}

	var accounts []*sg.AccountMeta
	for _, account := range ixf.Accounts {
		accounts = append(accounts, &sg.AccountMeta{
			PublicKey:  sg.MustPublicKeyFromBase58(account),
			IsWritable: accountBook[account].Writable,
			IsSigner:   accountBook[account].Signer,
		})
	}

	ixData, err := base58.Decode(ixf.Data)
	if err != nil {
		Logger.Fatal(fmt.Sprintf("ix data:'%s' base58 decode err:%v", ixf.Data, err))
	}

	ix, err := raydium_amm.DecodeInstruction(accounts, ixData)
	if err != nil {
		Logger.Fatal(fmt.Sprintf("raydium_amm.DecodeInstruction err:%v", err))
	}

	var ixesInnerParsed []rpc.InstructionInnerParsed
	for _, txMetaIxes := range tx.Meta.InnerInstructions {
		if txMetaIxes.Index == ixIndex {
			ixesInnerParsed = txMetaIxes.Instructions
			break
		}
	}

	if ixesInnerParsed == nil {
		Logger.Fatal("empty InnerInstructions") // TODO detail
	}

	switch ix.TypeID.Uint8() {
	case raydium_amm.Instruction_SwapBaseIn:
		pt.ParseIxSwapBaseIn(blockHeight, blockTime, txIndex, tx.Transaction.Signatures[0], ix, ixesInnerParsed)
	case raydium_amm.Instruction_SwapBaseOut:
		pt.ParseIxSwapBaseOut(blockHeight, blockTime, txIndex, tx.Transaction.Signatures[0], ix, ixesInnerParsed)
	}
}

func (pt *ParserTxRaydiumAmm) Done() {
	close(pt.marketCh)
	close(pt.txCh)
}
