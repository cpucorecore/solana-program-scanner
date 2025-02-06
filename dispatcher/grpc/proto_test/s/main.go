package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"solana-program-scanner/dispatcher/grpc/proto"
)

var (
	port = flag.Int("port", 50001, "The server port")
)

type serverSol struct {
	proto.UnimplementedSolServer
}

func (s *serverSol) SendBlockTxs(stream proto.Sol_SendBlockTxsServer) error {
	var totalTxs int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Stream ended, processed %d transactions in total", totalTxs)
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			log.Printf("Error receiving: %v", err)
			return err
		}

		block := req.GetBlock()
		txs := req.GetTxs()
		totalTxs += len(txs)

		log.Printf("Received block: %d, timestamp: %d, transactions: %d",
			block.Block, block.BlockAt, len(txs))

		for _, tx := range txs {
			txBytes, _ := json.Marshal(tx)
			log.Printf(string(txBytes))
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterSolServer(s, &serverSol{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
