package msg_broker

type msgBrokerLocal[T any] struct {
	ch chan T
}

func NewMsgBrokerLocal[T any](capacity int) MsgBroker[T] {
	return &msgBrokerLocal[T]{
		ch: make(chan T, capacity),
	}
}

func (b *msgBrokerLocal[T]) Produce(m T) {
	b.ch <- m
}

func (b *msgBrokerLocal[T]) ProducerChan() chan<- T {
	return b.ch
}

func (b *msgBrokerLocal[T]) ConsumerChan() <-chan T {
	return b.ch
}

func (b *msgBrokerLocal[T]) Close() {
	close(b.ch)
}
