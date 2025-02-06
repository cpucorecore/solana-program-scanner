package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/blocto/solana-go-sdk/rpc"
)

const (
	RAYDIUM_AMM_V4_ID = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
)

var (
	transactionVersion = uint8(0)
	rewards            = false
	getBlockConfig     = rpc.GetBlockConfig{
		Encoding:                       rpc.GetBlockConfigEncodingJsonParsed,
		TransactionDetails:             rpc.GetBlockConfigTransactionDetailsFull,
		Rewards:                        &rewards,
		Commitment:                     rpc.CommitmentFinalized,
		MaxSupportedTransactionVersion: &transactionVersion,
	}
	ctx = context.Background()
)

func getBlock(rpcCli *rpc.RpcClient, slot uint64) (*rpc.GetBlock, bool) {
	for {
		resp, err := rpcCli.GetBlockWithConfig(ctx, slot, getBlockConfig)
		if err != nil {
			continue
		}

		if resp.Error != nil {
			return nil, false
		}
		return resp.Result, true
	}
}

func main() {
	c := rpc.NewRpcClient("todo")

	resp, err := c.GetSlot(context.Background())
	if err != nil {
		panic(err)
	}

	newestSlot := resp.Result
	//startSlot := uint64(269132093)
	endSlot := newestSlot
	startSlot := endSlot - 1000

	for slot := startSlot; slot < endSlot; slot++ {
		block, ok := getBlock(&c, slot)
		if !ok {
			continue
		}

		for _, tx := range block.Transactions {
			if tx.Meta.Err != nil {
				continue
			}
			getDeepRaydiumCall(tx.Transaction.Signatures[0], tx.Meta.InnerInstructions)
		}
	}
}

type Instruction struct {
	Accounts    []string
	Data        string
	ProgramId   string
	StackHeight int
}

func UnmarshalInstruction(instructions any) (*Instruction, bool) {
	bs, err := json.Marshal(instructions)
	if err != nil {
		return nil, false
	}
	var tx Instruction
	err = json.Unmarshal(bs, &tx)
	if err != nil {
		return nil, false
	}
	return &tx, true
}

func getDeepRaydiumCall(signature string, innerInstructions []rpc.TransactionMetaInnerInstruction) {
	raydiumCallCount := 0
	maxStackHeight := 0

	for _, inner := range innerInstructions {
		for _, ix := range inner.Instructions {
			ixMap, ok := ix.(map[string]interface{})
			if !ok {
				log.Fatal("1")
				continue
			}
			if ixMap["parsed"] != nil {
				continue
			}

			programId, ok := ixMap["programId"].(string)
			if !ok {
				log.Fatal("2")
				continue
			}

			if programId == RAYDIUM_AMM_V4_ID {
				if stackHeight, ok := ixMap["stackHeight"].(float64); ok {
					stackHeight := int(stackHeight)
					if stackHeight >= 2 {
						fmt.Printf("stackHeight: %d, signature: %s\n", stackHeight, signature)
						raydiumCallCount++
						if stackHeight > maxStackHeight {
							maxStackHeight = stackHeight
						}
					}
				} else {
					log.Fatal("3")
				}
			}
		}
	}

	if raydiumCallCount > 0 {
		fmt.Printf("signature: %s, maxStackHeight: %d, raydiumCallCount:%d\n", signature, maxStackHeight, raydiumCallCount)
	}
}
