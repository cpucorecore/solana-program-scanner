package msg_broker

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"

	"solana-program-scanner/config"
	"solana-program-scanner/log"
	"solana-program-scanner/utils"
)

type MqMsgPublisher[T any] struct {
	queueName   string
	conf        *config.PublisherConf
	inputBuffer chan T
	inspector   *MqInspector
	conn        *rabbitmq.Conn
	publisher   *rabbitmq.Publisher
	mu          sync.Mutex
	done        bool
	rateLimiter *rateLimiter
}

func NewMqMsgPublisher[T any](mqUrl, queueName string, conf *config.PublisherConf) *MqMsgPublisher[T] {
	inspector := NewMqInspector(mqUrl, queueName, time.Second*time.Duration(conf.InspectIntervalBySecond))

	conn, err := rabbitmq.NewConn(mqUrl, rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		log.Logger.Fatal("connect to mq err", zap.Error(err))
	}

	publisher, err := rabbitmq.NewPublisher(conn, rabbitmq.WithPublisherOptionsLogging)
	if err != nil {
		log.Logger.Fatal("create publisher err", zap.Error(err))
	}

	return &MqMsgPublisher[T]{
		queueName:   queueName,
		conf:        conf,
		inputBuffer: make(chan T, conf.BufferSize),
		inspector:   inspector,
		conn:        conn,
		publisher:   publisher,
		rateLimiter: &rateLimiter{errWaitUnit: time.Millisecond * time.Duration(conf.WaitUnitOnErrByMs)},
	}
}

func (p *MqMsgPublisher[T]) Run() {
	for {
		if p.inspector.Messages() < p.conf.MqMessagesMin {
			if p.publishN(p.conf.MqMessagesMax) {
				log.Logger.Info("no more msg to publish, do setDone()")
				p.setDone()
				return
			}
		}
	}
}

func (p *MqMsgPublisher[T]) WaitDone() {
	log.Logger.Info("MqMsgPublisher WaitDone begin")
	for !p.isDone() {
		log.Logger.Info("MqMsgPublisher isDone=false ")
		time.Sleep(time.Second)
	}
	log.Logger.Info("MqMsgPublisher WaitDone success")
}

func (p *MqMsgPublisher[T]) setDone() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.done = true
}

func (p *MqMsgPublisher[T]) isDone() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.done
}

func (p *MqMsgPublisher[T]) Publish(msg T) {
	p.inputBuffer <- msg
}

func (p *MqMsgPublisher[T]) mustPublish(data []byte) {
	var err error
	for {
		err = p.publisher.Publish(
			data,
			[]string{p.queueName},
			rabbitmq.WithPublishOptionsContentType("text/plain"),
		)

		if err != nil {
			log.Logger.Error("MqMsgPublisher publish data fail", zap.ByteString("data", data), zap.Error(err))
			p.rateLimiter.OnErr()
			continue
		}

		p.rateLimiter.OnDone()
		return
	}
}

func (p *MqMsgPublisher[T]) publishN(n int) (done bool) {
	log.Logger.Info("MqMsgPublisher publish messages begin", zap.Int("messages", n))

	for i := 1; i <= n; i++ {
		select {
		case msg, ok := <-p.inputBuffer:
			if !ok {
				log.Logger.Info("MqMsgPublisher inputBuffer @ done", zap.Int("messages", i))
				return true
			}

			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Logger.Error("MqMsgPublisher marshal msg fail", zap.Any("msg", msg), zap.Error(err))
				continue
			}

			p.mustPublish(msgBytes)
		}
	}
	log.Logger.Info("MqMsgPublisher publish finish", zap.Int("messages", n))
	return false
}

func (p *MqMsgPublisher[T]) getChan() chan T {
	return p.inputBuffer
}

type rateLimiter struct {
	errWaitUnit time.Duration
	errCounter  utils.MutexCounter
}

func (r *rateLimiter) OnErr() {
	errCnt := r.errCounter.Up()
	time.Sleep(r.errWaitUnit * time.Duration(errCnt))
}

func (r *rateLimiter) OnDone() {
	r.errCounter.Reset()
}
