package main

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

const (
	wsBaseUrl = "wss://hermes.pyth.network/ws"
	priceId   = "0xef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d"
)

func main() {
	u := url.URL{Scheme: "wss", Host: "hermes.pyth.network", Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	subscribeMessage := map[string]interface{}{
		"type": "subscribe",
		"ids":  []string{priceId},
	}

	err = c.WriteJSON(subscribeMessage)
	if err != nil {
		log.Fatal("write:", err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal("read:", err)
		}
		log.Printf("recv: %s", message)
	}
}
