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
	DbDriver                 = "postgres"
	DbErrCodeUniqueConstrain = "23505"
)

type PostgresAttendant struct {
	dataSource string
	cli        *xorm.Engine
}

func NewPostgresAttendant(dataSource string) *PostgresAttendant {
	engine, err := xorm.NewEngine(DbDriver, dataSource)
	if err != nil {
		Logger.Fatal(fmt.Sprintf("db NewEngine err:%v", err))
	}

	return &PostgresAttendant{
		dataSource: dataSource,
		cli:        engine,
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
