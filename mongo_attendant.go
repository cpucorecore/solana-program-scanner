package main

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	Database     = "raydium_amm"
	CollectionIx = "ixes"
)

type MongoAttendant struct {
	cli  *mongo.Client
	ixCh chan bson.M
}

func NewMongoAttendant(ixCh chan bson.M, client *mongo.Client) *MongoAttendant {
	return &MongoAttendant{
		cli:  client,
		ixCh: ixCh,
	}
}

func (ma *MongoAttendant) startServe(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go ma.serveIx(ctx, wg)
}

func (ma *MongoAttendant) serveIx(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case ix := <-ma.ixCh:
			if ix == nil {
				Logger.Info("ixCh @ done")
				return
			}

			_, err := ma.cli.Database(Database).Collection(CollectionIx).InsertOne(ctx, ix)
			if err != nil {
				Logger.Error(fmt.Sprintf("insert ix err:%s", err.Error()))
			}
		}
	}
}
