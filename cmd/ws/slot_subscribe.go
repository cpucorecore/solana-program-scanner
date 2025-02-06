package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type SlotSubscribeRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type SlotNotification struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Result struct {
			Parent uint64 `json:"parent"`
			Root   uint64 `json:"root"`
			Slot   uint64 `json:"slot"`
		} `json:"result"`
		Subscription int `json:"subscription"`
	} `json:"params"`
}

type SlotSubscriber struct {
	endpoint     string
	conn         *websocket.Conn
	done         chan struct{}
	maxRetries   int
	retryDelay   time.Duration
	pingInterval time.Duration
	pongWait     time.Duration
}

func NewSlotSubscriber(endpoint string) *SlotSubscriber {
	return &SlotSubscriber{
		endpoint:     endpoint,
		done:         make(chan struct{}),
		maxRetries:   5,
		retryDelay:   time.Second * 5,
		pingInterval: time.Second * 30,
		pongWait:     time.Second * 10,
	}
}

func (s *SlotSubscriber) connect() error {
	var err error
	for i := 0; i < s.maxRetries; i++ {
		s.conn, _, err = websocket.DefaultDialer.Dial(s.endpoint, nil)
		if err == nil {
			return nil
		}
		time.Sleep(s.retryDelay)
	}
	return fmt.Errorf("failed to connect after %d retries: %v", s.maxRetries, err)
}

func (s *SlotSubscriber) Subscribe(ctx context.Context) (<-chan uint64, error) {
	if err := s.connect(); err != nil {
		return nil, err
	}

	go s.startHeartbeat(ctx)

	request := SlotSubscribeRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "slotSubscribe",
	}

	if err := s.conn.WriteJSON(request); err != nil {
		return nil, fmt.Errorf("write error: %v", err)
	}

	slotChan := make(chan uint64)
	go s.readPump(ctx, slotChan)

	return slotChan, nil
}

func (s *SlotSubscriber) startHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(s.pingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(s.pongWait)); err != nil {
				log.Printf("ping error: %v", err)
				s.reconnect(ctx)
				return
			}
		case <-ctx.Done():
			return
		case <-s.done:
			return
		}
	}
}

func (s *SlotSubscriber) reconnect(ctx context.Context) {
	s.conn.Close()
	for {
		if err := s.connect(); err != nil {
			log.Printf("重连失败: %v, 将在 %v 后重试", err, s.retryDelay)
			select {
			case <-ctx.Done():
				return
			case <-s.done:
				return
			case <-time.After(s.retryDelay):
				continue
			}
		}
		break
	}

	request := SlotSubscribeRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "slotSubscribe",
	}

	if err := s.conn.WriteJSON(request); err != nil {
		log.Printf("重新订阅失败: %v, 将重试重连", err)
		go s.reconnect(ctx)
		return
	}
}

func (s *SlotSubscriber) readPump(ctx context.Context, slotChan chan<- uint64) {
	defer func() {
		s.conn.Close()
		close(slotChan)
	}()

	s.conn.SetPongHandler(func(string) error {
		return s.conn.SetReadDeadline(time.Now().Add(s.pongWait))
	})

	for {
		var notification SlotNotification
		err := s.conn.ReadJSON(&notification)
		if err != nil {
			log.Printf("读取错误: %v", err)
			s.reconnect(ctx)
			continue
		}

		log.Printf("notification:%v", notification)
		if notification.Method == "slotNotification" {
			select {
			case slotChan <- notification.Params.Result.Root:
			case <-ctx.Done():
				return
			case <-s.done:
				return
			}
		}
	}
}

func (s *SlotSubscriber) Close() {
	close(s.done)
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	subscriber := NewSlotSubscriber("wss://api.mainnet-beta.solana.com")
	defer subscriber.Close()

	for {
		slotChan, err := subscriber.Subscribe(ctx)
		if err != nil {
			log.Printf("订阅失败: %v, 将在 5 秒后重试", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for slot := range slotChan {
			log.Printf("新区块: %d", slot)
		}

		select {
		case <-ctx.Done():
			return
		default:
			log.Printf("连接断开，准备重新订阅...")
			time.Sleep(5 * time.Second)
		}
	}
}
