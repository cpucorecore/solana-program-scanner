package db_commiter

import (
	"xorm.io/xorm"

	"solana-program-scanner/config"
	"solana-program-scanner/msg_broker"
	"solana-program-scanner/types/orms"
)

type BatchCommitterBuilder[T orms.Committable] struct {
	name        string
	conf        *config.BatchCommitConf
	msgConsumer msg_broker.MsgConsumer[T]
}

func NewBatchCommitterBuilder[T orms.Committable](
	name string,
	conf *config.BatchCommitConf,
	msgConsumer msg_broker.MsgConsumer[T],
) *BatchCommitterBuilder[T] {
	return &BatchCommitterBuilder[T]{
		name:        name,
		conf:        conf,
		msgConsumer: msgConsumer,
	}
}

func (b *BatchCommitterBuilder[T]) WithName(name string) *BatchCommitterBuilder[T] {
	b.name = name
	return b
}

func (b *BatchCommitterBuilder[T]) WithBatchSize(size int) *BatchCommitterBuilder[T] {
	b.conf.BatchSize = size
	return b
}

func (b *BatchCommitterBuilder[T]) WithFlushTimeout(timeout int) *BatchCommitterBuilder[T] {
	b.conf.FlushTimeoutByMs = timeout
	return b
}

func (b *BatchCommitterBuilder[T]) WithBufferSize(size int) *BatchCommitterBuilder[T] {
	b.conf.BufferSize = size
	return b
}

func (b *BatchCommitterBuilder[T]) WithMsgConsumer(dbObjConsumer msg_broker.MsgConsumer[T]) *BatchCommitterBuilder[T] {
	b.msgConsumer = dbObjConsumer
	return b
}

func (b *BatchCommitterBuilder[T]) Build(db *xorm.Engine) *BatchCommitter[T] {
	return NewBatchCommitter[T](b.name, db, b.conf, b.msgConsumer)
}
