package parser

import (
	"fmt"
	"time"

	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/price_service"
	"solana-program-scanner/types"
)

type BlockParser struct {
	priceService price_service.PriceService
	txParser     *RaydiumAmmParser
	sem          chan struct{}
}

func NewBlockParser(priceService price_service.PriceService, txParser *RaydiumAmmParser) *BlockParser {
	return &BlockParser{
		priceService: priceService,
		txParser:     txParser,
		sem:          make(chan struct{}, 24), // TODO config
	}
}

type Tx struct {
	Hash  string
	Index int
	Ixs   []*RaydiumInstructionParsed
}

type BlockHeader struct {
	Slot    uint64
	Height  int64
	BlockAt int64
}

type Block struct {
	Block        *BlockHeader
	Txs          []*Tx
	PoolBalances []*PoolBalance
}

func (p *BlockParser) ParseBlock(b *types.BlockWithSlot) *Block {
	log.Logger.Info(fmt.Sprintf("slot:%d ParseBlock get sem", b.Slot))
	p.sem <- struct{}{}
	log.Logger.Info(fmt.Sprintf("slot:%d ParseBlock get sem ok", b.Slot))
	defer func() {
		<-p.sem
	}()

	price := p.priceService.GetPrice(b.Block.ParentSlot, *b.Block.BlockTime)
	log.Logger.Info(fmt.Sprintf("slot:%d, blockTime:%d, price:%s", b.Slot, *b.Block.BlockTime, price))

	txs := make([]*Tx, 0, len(b.Block.Transactions))
	blockPoolBalances := make([]*PoolBalance, 0, 1024)
	monitor.BlockTxCount.Set(float64(len(b.Block.Transactions)))
	now := time.Now()
	for index, tx := range b.Block.Transactions {
		if tx.Meta.Err != nil {
			continue
		}

		raydiumInstructions := GetRaydiumInstructions(&tx)
		if raydiumInstructions.Empty() {
			continue
		}

		accountBookWithIndex := CreateTransactionAccountBook(&tx)
		ixsParsed := p.txParser.parseRaydiumInstructions(tx.Transaction.Signatures[0], accountBookWithIndex, raydiumInstructions, price)
		txs = append(txs, &Tx{
			Hash:  tx.Transaction.Signatures[0], // TODO check len(tx.Transaction.Signatures[0]) > 1
			Index: index,
			Ixs:   ixsParsed,
		})

		txPoolBalances := getTxPoolBalances(&tx, accountBookWithIndex, ixsParsed)
		blockPoolBalances = append(blockPoolBalances, txPoolBalances...)
	}
	monitor.ParseBlockDuration.Observe(float64(time.Since(now).Milliseconds()))

	mergedBlockPoolBalances := mergePoolBalances(blockPoolBalances)
	return &Block{
		Block: &BlockHeader{
			Slot:    b.Slot,
			Height:  *b.Block.BlockHeight,
			BlockAt: *b.Block.BlockTime,
		},
		Txs:          txs,
		PoolBalances: mergedBlockPoolBalances,
	}
}
