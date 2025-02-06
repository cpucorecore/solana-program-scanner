package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"solana-program-scanner/config"
	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/parser"
)

type Client struct {
	ID            string
	conf          *config.KafkaConf
	sendTimeout   time.Duration
	asyncProducer sarama.AsyncProducer
	msgPackager   func(block *parser.Block) ([]byte, error)
}

func NewClient(conf *config.KafkaConf) *Client {
	client := &Client{
		ID:          conf.ID,
		conf:        conf,
		sendTimeout: time.Millisecond * time.Duration(conf.SendTimeoutByMs),
	}

	config := sarama.NewConfig()
	config.Net.TLS.Enable = false
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 100 * time.Millisecond
	config.Producer.Retry.Max = 10

	asyncProducer, err := sarama.NewAsyncProducer(conf.Brokers, config)
	if err != nil {
		log.Logger.Fatal("kafka NewAsyncProducer err", zap.Error(err))
	}
	client.asyncProducer = asyncProducer
	client.processErrors()

	return client
}

func (c *Client) WithMsgPackager(msgPackager func(block *parser.Block) ([]byte, error)) *Client {
	c.msgPackager = msgPackager
	return c
}

func (c *Client) Close() {
	c.asyncProducer.Close()
}

func (c *Client) processErrors() {
	errCh := c.asyncProducer.Errors()
	go func() {
		for {
			err, ok := <-errCh
			if !ok {
				log.Logger.Info("kafka asyncProducer error @ done", zap.Error(err))
				return
			}
			log.Logger.Info("kafka asyncProducer error", zap.Error(err))
		}
	}()
}

func (c *Client) Send(block *parser.Block) error {
	data, err := c.msgPackager(block)
	if err != nil {
		return fmt.Errorf("ConvertToKafkaReq error: %v", err)
	}

	now := time.Now()
	c.asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: c.conf.Topic,
		Value: sarama.ByteEncoder(data),
	}
	monitor.SendBlockKafkaDuration.Observe(float64(time.Since(now).Milliseconds()))

	return nil
}

func (c *Client) Id() string {
	return c.ID
}

func (c *Client) Dispatch(block *parser.Block) error {
	return c.Send(block)
}

func (c *Client) Stop() {
	c.Close()
}
