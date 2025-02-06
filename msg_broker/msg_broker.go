package msg_broker

import "sync"

type Runnable interface {
	Run(wg *sync.WaitGroup)
}

type MsgProducer[T any] interface {
	Produce(msg T)
	ProducerChan() chan<- T
	Close()
}

type MsgConsumer[T any] interface {
	ConsumerChan() <-chan T
}

type MsgBroker[T any] interface {
	MsgProducer[T]
	MsgConsumer[T]
}
