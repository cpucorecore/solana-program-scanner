package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type BlockSubscriber struct {
	conn *websocket.Conn
}

// BlockSubscription 区块订阅请求
type BlockSubscription struct {
	Jsonrpc string   `json:"jsonrpc"`
	ID      int      `json:"id"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
}

// BlockResponse 区块订阅响应
type BlockResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Result struct {
			Context struct {
				Slot uint64 `json:"slot"`
			} `json:"context"`
			Value struct {
				BlockHeight uint64 `json:"blockHeight"`
				BlockTime   int64  `json:"blockTime"`
				Blockhash   string `json:"blockhash"`
				ParentSlot  uint64 `json:"parentSlot"`
			} `json:"value"`
		} `json:"result"`
		Subscription int `json:"subscription"`
	} `json:"params"`
}

// NewBlockSubscriber 创建区块订阅器
func NewBlockSubscriber(wsEndpoint string) (*BlockSubscriber, error) {
	// 建立 WebSocket 连接
	log.Printf("connecting to %s", wsEndpoint)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %v", err)
	} else {
		log.Println("Connected to WebSocket")
	}

	return &BlockSubscriber{
		conn: conn,
	}, nil
}

type SubscriptionConfig struct {
	Commitment         string `json:"commitment"`
	Encoding           string `json:"encoding"`
	ShowRewards        bool   `json:"showRewards"`
	TransactionDetails string `json:"transactionDetails"`
}

// Subscribe 订阅区块
func (bs *BlockSubscriber) Subscribe(ctx context.Context) error {
	// 发送订阅请求
	//subscription := BlockSubscription{
	//	Jsonrpc: "2.0",
	//	ID:      1,
	//	Method:  "blockSubscribe",
	//	Params:  []string{"all"},
	//}

	type BlockSubscription struct {
		Jsonrpc string        `json:"jsonrpc"`
		ID      int           `json:"id"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"` // 改为 interface{} 切片
	}

	//if err := bs.conn.WriteJSON(subscription); err != nil {
	//	return fmt.Errorf("failed to send subscription request: %v", err)
	//} else {
	//	log.Printf("subscribed to finalized block")
	//}

	config := SubscriptionConfig{
		Commitment:         "confirmed",
		Encoding:           "jsonParsed",
		ShowRewards:        false,
		TransactionDetails: "full",
	}

	// 发送订阅请求
	s := BlockSubscription{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "blockSubscribe",
		Params:  []interface{}{"all", config}, // 包含 "all" 和配置对象
	}

	if err := bs.conn.WriteJSON(s); err != nil {
		return fmt.Errorf("failed to send subscription request: %v", err)
	} else {
		log.Printf("subscribed to finalized block")
	}

	// 处理订阅消息
	for {
		select {
		case <-ctx.Done():
			return bs.conn.Close()
		default:
			// 读取消息
			_, message, err := bs.conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}
			log.Printf("recv: %s", message)

			// 解析消息
			var response BlockResponse
			if err := json.Unmarshal(message, &response); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			// 处理区块信息
			bs.processBlock(&response)
		}
	}
}

// processBlock 处理区块信息
func (bs *BlockSubscriber) processBlock(response *BlockResponse) {
	// 只处理区块更新消息
	if response.Method != "blockNotification" {
		return
	}

	block := response.Params.Result.Value
	fmt.Printf("\nNew Block Received:\n")
	fmt.Printf("Slot: %d\n", response.Params.Result.Context.Slot)
	fmt.Printf("Block Height: %d\n", block.BlockHeight)
	fmt.Printf("Block Time: %d\n", block.BlockTime)
	fmt.Printf("Block Hash: %s\n", block.Blockhash)
	fmt.Printf("Parent Slot: %d\n", block.ParentSlot)
}

// Close 关闭连接
func (bs *BlockSubscriber) Close() error {
	if bs.conn != nil {
		return bs.conn.Close()
	}
	return nil
}

func main() {
	// 创建区块订阅器
	subscriber, err := NewBlockSubscriber("wss://wider-intensive-hill.solana-mainnet.quiknode.pro/952048cc248a07c391d2128da1fdaeb6349571d9")
	if err != nil {
		log.Fatal(err)
	}
	defer subscriber.Close()

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 开始订阅区块
	if err := subscriber.Subscribe(ctx); err != nil {
		log.Fatal(err)
	}
}
