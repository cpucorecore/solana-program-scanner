package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"
	sg "github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"

	"solana-program-scanner/idls/raydium_amm"
)

type ParserTxRaydiumAmm struct {
	txCh         chan *OrmTx
	getterMarket GetterMarket
	cacheMarket  Cache[string, *OrmMarket]
}

func NewParserTxRaydiumAmm(
	txCh chan *OrmTx,
	getterMarket GetterMarket,
	cacheMarket Cache[string, *OrmMarket],
) *ParserTxRaydiumAmm {
	return &ParserTxRaydiumAmm{
		txCh:         txCh,
		getterMarket: getterMarket,
		cacheMarket:  cacheMarket,
	}
}

func (pt *ParserTxRaydiumAmm) ParseIxSwapBaseIn(
	blockHeight int64,
	blockTime int64,
	txIndex int,
	txHash string,
	ix *raydium_amm.Instruction,
) {
	ixSwapBaseIn, ok := ix.Impl.(*raydium_amm.SwapBaseIn)
	if !ok {
		Logger.Fatal("")
	}

	marketAddress := ixSwapBaseIn.AccountMetaSlice[1].PublicKey.String()
	market, err := pt.cacheMarket.Get(marketAddress)
	if err != nil {
		Logger.Fatal(fmt.Sprintf("market cache get err:%v", err)) // TODO check
	}

	if market == nil {
		market, err = pt.getterMarket.getMarket(marketAddress) // TODO save to db table:market/pair
		if err != nil {
			Logger.Fatal(fmt.Sprintf("get market err:%v", err))
		}

		err = pt.cacheMarket.Set(market.Address, market)
		if err != nil {
			Logger.Warn(fmt.Sprintf("market cache set err:%v", err))
		}
	}

	ormTx := &OrmTx{
		TxHash:        txHash,
		Event:         0, // TODO
		Token0Amount:  strconv.FormatUint(*ixSwapBaseIn.AmountIn, 10),
		Token1Amount:  strconv.FormatUint(*ixSwapBaseIn.MinimumAmountOut, 10), // TODO check 0
		Maker:         ixSwapBaseIn.AccountMetaSlice[1].PublicKey.String(),    // TODO
		Token0Address: market.BaseMint,
		Token1Address: market.QuoteMint,
		Block:         blockHeight,
		BlockAt:       time.Unix(blockTime, 0),
		Index:         txIndex,
	}

	pt.txCh <- ormTx

	Logger.Info(fmt.Sprintf("tx:%v", ormTx))
}

func (pt *ParserTxRaydiumAmm) ParseTx(
	blockHeight int64,
	blockTime int64,
	txIndex int,
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

	switch ix.TypeID.Uint8() {
	case raydium_amm.Instruction_SwapBaseIn:
		pt.ParseIxSwapBaseIn(blockHeight, blockTime, txIndex, tx.Transaction.Signatures[0], ix)
	case raydium_amm.Instruction_SwapBaseOut:
	}
}

func (pt *ParserTxRaydiumAmm) Done() {
	close(pt.txCh)
}
