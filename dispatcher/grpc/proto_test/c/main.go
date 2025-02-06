package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"solana-program-scanner/dispatcher/grpc/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewSolClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	stream, err := client.SendBlockTxs(ctx)
	if err != nil {
		log.Fatalf("Failed to create stream: %v", err)
	}

	testBlocks := []uint64{12345, 12346, 12347}
	for _, blockNum := range testBlocks {
		block := &proto.Block{
			Block:   blockNum,
			BlockAt: time.Now().Unix(),
		}

		txs := []*proto.Tx{
			{
				Maker:         "0x123abc...",
				Token0Address: "0xtoken0...",
				Token1Address: "0xtoken1...",
				Token0Amount:  "100.5",
				Token1Amount:  "200.75",
				PriceUsd:      "1.5",
				AmountUsd:     "150.75",
				Block:         blockNum,
				BlockIndex:    0,
				Event:         "SWAP",
				TxHash:        "0xtx" + string(blockNum),
				TxIndex:       0,
				BlockAt:       time.Now().Format(time.RFC3339),
			},
			{
				Maker:         "0x456def...",
				Token0Address: "0xtoken2...",
				Token1Address: "0xtoken3...",
				Token0Amount:  "300.25",
				Token1Amount:  "400.50",
				PriceUsd:      "2.0",
				AmountUsd:     "600.50",
				Block:         blockNum,
				BlockIndex:    1,
				Event:         "LIQUIDITY",
				TxHash:        "0xtx" + string(blockNum) + "_1",
				TxIndex:       1,
				BlockAt:       time.Now().Format(time.RFC3339),
			},
		}

		err = stream.Send(&proto.SendBlockTxsReq{
			Block: block,
			Txs:   txs,
		})
		if err != nil {
			log.Fatalf("Failed to send block %d: %v", blockNum, err)
		}
		log.Printf("Sent block %d with %d transactions", blockNum, len(txs))

		time.Sleep(time.Millisecond * 500)
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error closing stream: %v", err)
	}

	log.Println("Successfully sent all block transactions")
}
