package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

func SubscribeBlock(endpoint string) error {
	client, err := ws.Connect(context.Background(), endpoint)
	if err != nil {
		return fmt.Errorf("failed to connect websocket: %v", err)
	}
	defer client.Close()
	log.Printf("Subscribed to %s", endpoint)

	r := false
	m := uint64(1)
	opt := &ws.BlockSubscribeOpts{
		Commitment:                     rpc.CommitmentFinalized,
		TransactionDetails:             "jsonParsed",
		Rewards:                        &r,
		MaxSupportedTransactionVersion: &m,
	}
	sub, err := client.BlockSubscribe(
		ws.BlockSubscribeFilterAll(""),
		opt,
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe blocks: %v", err)
	}
	defer sub.Unsubscribe()
	log.Printf("Subscribed to %s", endpoint)

	for {
		block, err := sub.Recv()
		if err != nil {
			return fmt.Errorf("failed to receive block: %v", err)
		}

		processBlock(block)
	}
}

func processBlock(block *ws.BlockResult) {
	log.Printf(fmt.Sprintf("New block received:%v", block.Value.Block.BlockTime))
	marshal, err := json.Marshal(block.Value.Block)
	if err != nil {
		log.Printf("failed to marshal block: %v", err)
	}
	log.Printf(string(marshal))
}

func main() {
	endpoint := "wss://wider-intensive-hill.solana-mainnet.quiknode.pro/952048cc248a07c391d2128da1fdaeb6349571d9"
	if err := SubscribeBlock(endpoint); err != nil {
		log.Fatal(err)
	}
}
