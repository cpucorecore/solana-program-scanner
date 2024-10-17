package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoAttendant struct {
	cli  *mongo.Client
	ixCh chan bson.M
}

const (
	MongoDataSource = "mongodb://localhost:27017"
	Database        = "raydium_amm"
	CollectionIx    = "ixes"
)

func NewMongoAttendant(ixCh chan bson.M) *MongoAttendant {
	client, err := mongo.Connect(options.Client().ApplyURI(MongoDataSource))
	if err != nil {
		log.Fatal(err)
	}

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
