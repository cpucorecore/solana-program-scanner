package grpc

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"solana-program-scanner/config"
	"solana-program-scanner/dispatcher/grpc/proto"
	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/parser"
)

type Client struct {
	ID            string
	conf          *config.GrpcConf
	sendTimeout   time.Duration
	retryInterval time.Duration
	conn          *grpc.ClientConn
	stream        grpc.ClientStreamingClient[proto.SendBlockTxsReq, emptypb.Empty]
}

func NewClient(conf *config.GrpcConf) *Client {
	client := &Client{
		ID:            conf.ID,
		conf:          conf,
		sendTimeout:   time.Millisecond * time.Duration(conf.SendTimeoutByMs),
		retryInterval: time.Millisecond * time.Duration(conf.RetryIntervalByMs),
	}

	err := client.connect()
	if err != nil {
		log.Logger.Fatal("grpc connect err", zap.Error(err))
	}

	return client
}

func (c *Client) connect() error {
	conn, err := grpc.NewClient(c.conf.Target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("grpc dial err: %v", err)
	}

	client := proto.NewSolClient(conn)
	stream, err := client.SendBlockTxs(context.Background())
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create stream: %v", err)
	}

	c.conn = conn
	c.stream = stream
	return nil
}

func (c *Client) reconnect() error {
	c.Close()
	return c.connect()
}

func (c *Client) Close() {
	if c.stream != nil {
		c.stream.CloseSend()
		c.stream = nil
	}
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func ConvertToGrpcReq(block *parser.Block) *proto.SendBlockTxsReq {
	var req proto.SendBlockTxsReq
	req.Block = &proto.Block{
		Block:   block.Block.Slot,
		BlockAt: block.Block.BlockAt,
	}

	txsLen := 0
	for _, tx := range block.Txs {
		txsLen += len(tx.Ixs)
	}
	txs := make([]*proto.Tx, 0, txsLen)

	for _, tx := range block.Txs {
		for _, ix := range tx.Ixs {
			txs = append(txs, &proto.Tx{
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
	return &req
}

func (c *Client) connected() bool {
	return c.conn != nil && c.conn.GetState() == connectivity.Ready && c.stream != nil
}

func (c *Client) Send(block *parser.Block) error {
	req := ConvertToGrpcReq(block)
	now := time.Now()
	err := c.stream.Send(req)
	if err != nil {
		return err
	}
	monitor.SendBlockGrpcDuration.Observe(float64(time.Since(now).Milliseconds()))
	return nil
}

func (c *Client) SendWithRetry(block *parser.Block) error {
	retryCnt := 0
	for retryCnt < c.conf.MaxRetry {
		err := c.Send(block)
		if err == nil {
			return nil
		}

		log.Logger.Error(fmt.Sprintf("send block:%d err: %v", block.Block.Slot, err))
		if err == io.EOF {
			log.Logger.Warn("EOF error, to reconnect")
			if reconnectErr := c.reconnect(); reconnectErr != nil {
				log.Logger.Error(fmt.Sprintf("reconnect err: %v", reconnectErr))
				return reconnectErr
			}
		}

		retryCnt++
		log.Logger.Error(fmt.Sprintf("block:%d retry %d/%d", block.Block.Slot, retryCnt, c.conf.MaxRetry))
		time.Sleep(c.retryInterval * time.Duration(retryCnt))
	}

	return fmt.Errorf("failed to send block:%d after %d retries", block.Block.Slot, c.conf.MaxRetry)
}

func (c *Client) Id() string {
	return c.ID
}

func (c *Client) Dispatch(block *parser.Block) error {
	if !c.connected() {
		if err := c.reconnect(); err != nil {
			return err
		}
	}
	return c.SendWithRetry(block)
}

func (c *Client) Stop() {
	c.Close()
}
