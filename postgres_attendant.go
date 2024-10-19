package main

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/lib/pq"
	"xorm.io/xorm"
)

const (
	DbErrCodeUniqueConstrain = "23505"
)

type PostgresAttendant struct {
	cli *xorm.Engine
}

func NewPostgresAttendant(engine *xorm.Engine) *PostgresAttendant {
	return &PostgresAttendant{
		cli: engine,
	}
}

func (pa *PostgresAttendant) serveTx(ctx context.Context, wg *sync.WaitGroup, txCh chan *OrmTx) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case tx := <-txCh:
			if tx == nil {
				Logger.Info("txCh @ done")
				return
			}

			_, err := pa.cli.Insert(tx)
			if err != nil {
				var pgErr *pq.Error
				if errors.As(err, &pgErr) && string(pgErr.Code) == DbErrCodeUniqueConstrain {
					continue
				}

				Logger.Fatal(fmt.Sprintf("%v", err))
			}
		}
	}
}

func (pa *PostgresAttendant) serveMarket(ctx context.Context, wg *sync.WaitGroup, marketCh chan *OrmMarket) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case market := <-marketCh:
			if market == nil {
				Logger.Info("marketCh @ done")
				return
			}

			_, err := pa.cli.Insert(market)
			if err != nil {
				var pgErr *pq.Error
				if errors.As(err, &pgErr) && string(pgErr.Code) == DbErrCodeUniqueConstrain {
					continue
				}

				Logger.Fatal(fmt.Sprintf("%v", err))
			}
		}
	}
}
