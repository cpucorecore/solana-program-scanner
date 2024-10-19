package main

import (
	"encoding/json"
	"go.uber.org/zap"

	"github.com/blocto/solana-go-sdk/rpc"
)

type ParserTx struct {
	parserTxRaydiumAmm *ParserTxRaydiumAmm
}

func NewParserTx(parserTxRaydiumAmm *ParserTxRaydiumAmm) *ParserTx {
	return &ParserTx{
		parserTxRaydiumAmm: parserTxRaydiumAmm,
	}
}

func (pt *ParserTx) ParseTx(blockHeight int64, blockTime int64, txIndex int, tx *rpc.GetBlockTransaction) {
	var ixf rpc.InstructionFull
	for _, instruction := range tx.Transaction.Message.Instructions {
		ixJson, err := json.Marshal(instruction)
		if err != nil {
			Logger.Fatal("json encode error", zap.Error(err))
		}

		err = json.Unmarshal(ixJson, &ixf)
		if err != nil {
			Logger.Fatal("json decode error", zap.String("json", string(ixJson)), zap.Error(err))
		}

		if ixf.ProgramId == RadiumAmmAddressMainnet {
			pt.parserTxRaydiumAmm.ParseTx(blockHeight, blockTime, txIndex, tx, &ixf)
		}
	}
}

func (pt *ParserTx) Done() {
	pt.parserTxRaydiumAmm.Done()
}
