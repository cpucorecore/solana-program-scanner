package db_commiter

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/lib/pq"
	"xorm.io/xorm"

	"solana-program-scanner/config"
	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/msg_broker"
	"solana-program-scanner/types/orms"
)

const (
	DbErrCodeUniqueConstrain = "23505"
)

type BatchCommitter[T orms.Committable] struct {
	identity        string
	db              *xorm.Engine
	batchCommitConf *config.BatchCommitConf
	dbObjConsumer   msg_broker.MsgConsumer[T]
}

func NewBatchCommitter[T orms.Committable](name string, db *xorm.Engine, batchCommitConf *config.BatchCommitConf, dbObjConsumer msg_broker.MsgConsumer[T]) *BatchCommitter[T] {
	return &BatchCommitter[T]{
		identity:        fmt.Sprintf("db_batch_commiter-%s", name),
		db:              db,
		batchCommitConf: batchCommitConf,
		dbObjConsumer:   dbObjConsumer,
	}
}

func (c *BatchCommitter[T]) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	ticker := time.NewTicker(time.Millisecond * time.Duration(c.batchCommitConf.FlushTimeoutByMs))
	defer ticker.Stop()

	batch := make([]T, 0, c.batchCommitConf.BatchSize)
	dbObjCh := c.dbObjConsumer.ConsumerChan()
	for {
		select {
		case data, ok := <-dbObjCh:
			if !ok {
				log.Logger.Info(fmt.Sprintf("%s dbObjCh @ done", c.identity))
				if len(batch) > 0 {
					c.flush(batch)
				}
				return
			}

			batch = append(batch, data)
			if len(batch) >= c.batchCommitConf.BatchSize {
				c.flush(batch)
				batch = make([]T, 0, c.batchCommitConf.BatchSize)
			}

		case <-ticker.C:
			if len(batch) > 0 {
				c.flush(batch)
				batch = make([]T, 0, c.batchCommitConf.BatchSize)
			}
		}
	}
}

func (c *BatchCommitter[T]) insertOneByOne(batch []T) {
	for _, bean := range batch {
		now := time.Now()
		cnt, err := c.db.InsertOne(bean)
		monitor.DBCommitOneByOneDuration.WithLabelValues(c.identity).Observe(float64(time.Since(now).Milliseconds()))
		if err != nil {
			var pgErr *pq.Error
			ok := errors.As(err, &pgErr)
			if ok {
				monitor.DBCommitErrCounter.WithLabelValues(c.identity, string(pgErr.Code)).Inc()
				if string(pgErr.Code) == DbErrCodeUniqueConstrain {
					log.Logger.Warn(fmt.Sprintf("%s insertOneByOne([%v]) UniqueConstrain err:%s, insert cnt:%d", c.identity, bean, err.Error(), cnt)) // TODO check
					continue
				}
			}
			log.Logger.Warn(fmt.Sprintf("%s insertOneByOne([%v])err:%v", c.identity, bean, err))
		}
		log.Logger.Debug(fmt.Sprintf("%s insertOneByOne([%v]) success, insert cnt:%d", c.identity, bean, cnt))
	}
}

func (c *BatchCommitter[T]) flush(batch []T) {
	now := time.Now()
	cnt, err := c.db.Insert(batch)
	duration := time.Since(now).Milliseconds()
	monitor.DBBatchCommitDuration.WithLabelValues(c.identity).Observe(float64(duration))
	if duration > (time.Second * time.Duration(15)).Milliseconds() {
		log.Logger.Debug(fmt.Sprintf("%s flush duration:%d cnt:%d", c.identity, duration, cnt))
	}

	if err != nil {
		var pgErr *pq.Error
		ok := errors.As(err, &pgErr)
		if ok {
			monitor.DBCommitErrCounter.WithLabelValues(c.identity, string(pgErr.Code)).Inc()
			if string(pgErr.Code) == DbErrCodeUniqueConstrain {
				log.Logger.Warn(fmt.Sprintf("%s flush batch[%v] UniqueConstrain err:%s, insert cnt:%d, do insertOneByOne", c.identity, batch, err.Error(), cnt)) // TODO check
				c.insertOneByOne(batch)
				return
			}
		}
		log.Logger.Warn(fmt.Sprintf("%s flush batch[%v] err:%v, insert rows:%d", c.identity, batch, err, cnt))
	}
	log.Logger.Debug(fmt.Sprintf("%s flush %d rows", c.identity, cnt))
}
