package msg_broker

import (
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"solana-program-scanner/log"
)

type MqInspector struct {
	mqUrl     string
	queueName string
	interval  time.Duration

	conn    *amqp.Connection
	channel *amqp.Channel

	mu        sync.Mutex
	queueInfo amqp.Queue
}

func NewMqInspector(mqUrl string, queueName string, interval time.Duration) *MqInspector {
	i := &MqInspector{
		mqUrl:     mqUrl,
		queueName: queueName,
		interval:  interval,
	}

	i.dial()
	i.connect()
	i.start()

	return i
}

func (i *MqInspector) Messages() int {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.queueInfo.Messages
}

func (i *MqInspector) start() {
	i.inspectQueue()

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-ticker.C:
				i.inspectQueue()
			}
		}
	}()
}

func (i *MqInspector) inspectQueue() {
	q, err := i.channel.QueueInspect(i.queueName)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Failed to inspect queue %s, err:%s", i.queueName, err.Error()))
		i.reconnect()
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.queueInfo = q
	log.Logger.Info(fmt.Sprintf("queue status:%v", i.queueInfo))
}

func (i *MqInspector) dial() {
	conn, err := amqp.Dial(i.mqUrl)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("MqInspector Failed to connect to mq: %v", err))
	}
	i.conn = conn
}

func (i *MqInspector) connect() {
	channel, err := i.conn.Channel()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("MqInspector Failed to connect to mq: %v", err))
	}
	i.channel = channel
}

func (i *MqInspector) reconnect() {
	if !i.channel.IsClosed() {
		i.channel.Close()
	}

	if !i.conn.IsClosed() {
		i.conn.Close()
	}

	i.dial()
	i.connect()
}
