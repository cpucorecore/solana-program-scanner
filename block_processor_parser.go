package main

import (
	"github.com/blocto/solana-go-sdk/rpc"
)

type BlockProcessorParser struct {
	parserTx *ParserTx
}

func (bpp *BlockProcessorParser) id() string {
	return "BlockProcessorParser"
}

const (
	RadiumAmmAddressMainnet = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
)

func (bpp *BlockProcessorParser) process(block *rpc.GetBlock) error {
	for index, tx := range block.Transactions {
		if tx.Meta.Err != nil {
			continue
		}

		bpp.parserTx.ParseTx(*block.BlockHeight, *block.BlockTime, uint64(index), &tx) // TODO go
	}

	return nil
}

func (bpp *BlockProcessorParser) done() {
	bpp.parserTx.Done()
}

var _ ObjProcessor[*rpc.GetBlock] = &BlockProcessorParser{}

func NewBlockProcessorParser(parserTx *ParserTx) *BlockProcessorParser {
	return &BlockProcessorParser{
		parserTx: parserTx,
	}
}
