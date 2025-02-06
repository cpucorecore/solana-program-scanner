package kafka

import (
	"encoding/json"
	"solana-program-scanner/parser"
)

type SendBlockPoolReq struct {
	Block    Block                 `json:"block"`
	PoolInfo []*parser.PoolBalance `json:"txs"`
}

func ConvertToPoolReq(block *parser.Block) ([]byte, error) {
	var req SendBlockPoolReq
	req.Block = Block{
		Block:   block.Block.Slot,
		BlockAt: block.Block.BlockAt,
	}

	req.PoolInfo = block.PoolBalances

	return json.Marshal(req)
}
