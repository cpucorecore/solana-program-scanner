package msg_broker

import (
	"encoding/json"
	"time"

	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"

	"solana-program-scanner/log"
)

type MqMsgDeliver[T any] struct {
	outputBuffer chan T
	conn         *rabbitmq.Conn
	consumer     *rabbitmq.Consumer
}

func NewMqMsgDeliver[T any](mqUrl string, queueName string, bufferSize int) *MqMsgDeliver[T] {
	//conn, err := rabbitmq.NewConn(mqUrl, rabbitmq.WithConnectionOptionsLogging)
	conn, err := rabbitmq.NewConn(
		mqUrl,
		rabbitmq.WithConnectionOptionsLogging,
		rabbitmq.WithConnectionOptionsReconnectInterval(time.Second*2),
	)
	if err != nil {
		log.Logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
	}

	//consumer, err := rabbitmq.NewConsumer(conn, queueName)
	consumer, err := rabbitmq.NewConsumer(
		conn,
		queueName,
		rabbitmq.WithConsumerOptionsQOSPrefetch(10),
	)
	if err != nil {
		log.Logger.Fatal("Failed to create consumer", zap.Error(err))
	}

	return &MqMsgDeliver[T]{
		outputBuffer: make(chan T, bufferSize),
		conn:         conn,
		consumer:     consumer,
	}
}

func (d *MqMsgDeliver[T]) Run() {
	err := d.consumer.Run(func(delivery rabbitmq.Delivery) rabbitmq.Action {
		var value T
		err := json.Unmarshal(delivery.Body, &value)
		if err != nil {
			log.Logger.Error("failed to unmarshal msg",
				zap.String("msg", string(delivery.Body)),
				zap.Error(err),
			)
		} else {
			d.outputBuffer <- value
		}
		return rabbitmq.Ack
	})
	if err != nil {
		log.Logger.Fatal("failed to consume msg", zap.Error(err))
	}
}

func (d *MqMsgDeliver[T]) StopDeliver() {
	d.consumer.Close()
	d.conn.Close()
	close(d.outputBuffer)
}

func (d *MqMsgDeliver[T]) getChan() chan T {
	return d.outputBuffer
}
