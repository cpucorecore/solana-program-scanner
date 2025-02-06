package msg_broker

import (
	"solana-program-scanner/config"
)

type msgBrokerMq[T any] struct {
	d *MqMsgDeliver[T]
	p *MqMsgPublisher[T]
}

func NewMsgBrokerMq[T any](conf *config.MsgBrokerMqConf) MsgBroker[T] {
	d := NewMqMsgDeliver[T](conf.MqUrl, conf.QueueName, conf.Deliver.BufferSize)
	p := NewMqMsgPublisher[T](conf.MqUrl, conf.QueueName, conf.Publisher)

	broker := &msgBrokerMq[T]{
		d: d,
		p: p,
	}

	broker.start()

	return broker
}

func (b *msgBrokerMq[T]) start() {
	go b.d.Run()
	go b.p.Run()
}

func (b *msgBrokerMq[T]) Produce(msg T) {
	b.p.Publish(msg)
}

func (b *msgBrokerMq[T]) ProducerChan() chan<- T {
	return b.p.getChan()
}

func (b *msgBrokerMq[T]) ConsumerChan() <-chan T {
	return b.d.getChan()
}

func (b *msgBrokerMq[T]) Close() {
	b.d.StopDeliver()
	close(b.p.inputBuffer)
	b.p.WaitDone()
}
