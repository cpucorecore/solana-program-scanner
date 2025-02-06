package kafka

import (
	"encoding/json"
	"strconv"

	"solana-program-scanner/parser"
)

type SendBlockTxsReq struct {
	Block Block `json:"block"`
	Txs   []Tx  `json:"txs"`
}

type Block struct {
	Block   uint64 `json:"block"`
	BlockAt int64  `json:"block_at"`
}

type Tx struct {
	Maker         string `json:"maker"`
	MarketAddress string `json:"market_address"`
	Token0Address string `json:"token0_address"`
	Token1Address string `json:"token1_address"`
	Token0Amount  string `json:"token0_amount"`
	Token1Amount  string `json:"token1_amount"`
	PriceUsd      string `json:"price_usd"`
	AmountUsd     string `json:"amount_usd"`
	Block         uint64 `json:"block"`
	BlockIndex    uint64 `json:"block_index"`
	Event         string `json:"event"`
	TxHash        string `json:"tx_hash"`
	TxIndex       uint64 `json:"tx_index"`
	BlockAt       string `json:"block_at"`
}

func ConvertToKafkaReq(block *parser.Block) ([]byte, error) {
	var req SendBlockTxsReq
	req.Block = Block{
		Block:   block.Block.Slot,
		BlockAt: block.Block.BlockAt,
	}

	txsLen := 0
	for _, tx := range block.Txs {
		txsLen += len(tx.Ixs)
	}
	txs := make([]Tx, 0, txsLen)

	for _, tx := range block.Txs {
		for _, ix := range tx.Ixs {
			txs = append(txs, Tx{
				Maker:         ix.Signer,
				MarketAddress: ix.MarketAddress,
				Token0Address: ix.TransferDetail.Token0Address,
				Token1Address: ix.TransferDetail.Token1Address,
				Token0Amount:  ix.TransferDetail.Token0Amount.String(),
				Token1Amount:  ix.TransferDetail.Token1Amount.String(),
				PriceUsd:      ix.TransferDetail.PriceUsd.String(),
				AmountUsd:     ix.TransferDetail.AmountUsd.String(),
				Block:         block.Block.Slot,
				BlockIndex:    uint64(tx.Index),
				Event:         string(ix.TransferDetail.Event),
				TxHash:        tx.Hash,
				TxIndex:       uint64(ix.Index),
				BlockAt:       strconv.FormatInt(block.Block.BlockAt, 10),
			})
		}
	}
	req.Txs = txs
	return json.Marshal(req)
}
