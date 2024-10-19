package main

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"solana-program-scanner/block_height_manager"
	"sync"
	"time"
	"xorm.io/xorm"
)

type Factory struct {
	redisOptions *redis.Options   // redis options
	redisCli     *redis.Client    // redis client
	postgresCli  *xorm.Engine     // postgres client
	mongoCli     *mongo.Client    // mongo client
	solanaClis   []*rpc.RpcClient // solana rpc clients

	// pipelines
	blockGetTaskCh chan uint64
	blockCh        chan *rpc.GetBlock
	txCh           chan *OrmTx
	ixCh           chan bson.M

	flowController      FlowController
	blockHeightManager  block_height_manager.BlockHeightManager
	mongoAttendant      *MongoAttendant
	postgresAttendant   *PostgresAttendant
	blockTaskDispatcher *BlockTaskDispatcher
	blockGetter         *BlockGetter
	blockProcessorAdmin BlockProcessorAdmin
	marketCache         Cache[string, *OrmMarket]
	marketGetter        MarketGetter
}

func (f *Factory) assembleRedisOptions(addr string, username string, password string) *redis.Options {
	f.redisOptions = &redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
	}

	return f.redisOptions
}

func (f *Factory) assembleRedisClient(options *redis.Options) {
	f.redisCli = redis.NewClient(options)
}

const (
	DbDriverPostgres = "postgres"
)

func (f *Factory) assemblePostgresClient(dataSource string) {
	engine, err := xorm.NewEngine(DbDriverPostgres, dataSource)
	if err != nil {
		f.destroy()
		Logger.Fatal(fmt.Sprintf("xorm.NewEngine err:%v", err))
	}

	f.postgresCli = engine
}

func (f *Factory) assembleMongoClient(uri string) {
	connect, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		f.destroy()
		Logger.Fatal(fmt.Sprintf("mongo.Connect err:%v", err))
	}

	f.mongoCli = connect
}

func (f *Factory) assembleSolanaClient() *rpc.RpcClient {
	r := rpc.NewRpcClient(conf.Solana.RpcEndpoint)
	f.solanaClis = append(f.solanaClis, &r)
	return &r
}

func (f *Factory) assemblePipelines() {
	f.blockGetTaskCh = make(chan uint64)
	f.blockCh = make(chan *rpc.GetBlock)
	f.txCh = make(chan *OrmTx, 1000)
	f.ixCh = make(chan bson.M, 100)
}

func (f *Factory) assembleFlowController() {
	f.flowController = NewFlowController(5, 5, time.Second, time.Second)
}

func (f *Factory) assembleBlockHeightManager() {
	f.blockHeightManager = block_height_manager.NewBlockHeightManager()
}

func (f *Factory) assembleMongoAttendant() {
	f.mongoAttendant = NewMongoAttendant(f.ixCh, f.mongoCli)
}

func (f *Factory) assemblePostgresAttendant() {
	f.postgresAttendant = NewPostgresAttendant(f.postgresCli)
}

func (f *Factory) assembleBlockTaskDispatcher() {
	f.blockTaskDispatcher = NewBlockTaskDispatcher(f.flowController)
}

func (f *Factory) assembleBlockGetter(blockGetterWorkerNumber int) {
	f.blockGetter = NewBlockGetter(blockGetterWorkerNumber, f.blockHeightManager, f.flowController)
}

func (f *Factory) assembleBlockProcessor() {
	f.blockProcessorAdmin = NewBlockProcessorAdmin(f.blockCh, f.txCh, f.ixCh, f.blockGetter)
}

func (f *Factory) assembleMarketCache() {
	f.marketCache = NewCacheRedisMarket(f.redisCli)
}

func (f *Factory) assembleMarketGetter() {
	f.marketGetter = f.blockGetter
}

func (f *Factory) assemble() *Factory {
	redisOptions := f.assembleRedisOptions("localhost:6379", "", "") //TODO config
	f.assembleRedisClient(redisOptions)
	const DataSource = "postgres://postgres:12345678@localhost:5432/postgres?sslmode=disable"
	f.assemblePostgresClient(DataSource) // TODO config
	MongoDataSource := "mongodb://localhost:27017"
	f.assembleMongoClient(MongoDataSource) // TODO config
	f.assembleSolanaClient()
	f.assemblePipelines()
	f.assembleFlowController()
	f.assembleBlockHeightManager()
	f.assembleMongoAttendant()
	f.assemblePostgresAttendant()
	f.assembleBlockTaskDispatcher()
	f.assembleBlockGetter(3) // TODO config
	f.assembleBlockProcessor()
	f.assembleMarketCache()
	f.assembleMarketGetter()

	startBlockHeight := f.blockGetter.getBlockHeightBySlot(conf.Solana.StartSlot)
	f.blockHeightManager.Init(startBlockHeight - 1)

	return f
}

func (f *Factory) runProducts(ctx context.Context) (*sync.WaitGroup, FlowController) {
	go f.flowController.startLog(time.Second * 5)

	var wg sync.WaitGroup

	wg.Add(1)
	go f.postgresAttendant.serveTx(ctx, &wg, f.txCh)

	wg.Add(1)
	go f.blockProcessorAdmin.run(ctx, &wg)

	wg.Add(1)
	go f.blockGetter.run(ctx, &wg, f.blockGetTaskCh, f.blockCh)

	wg.Add(1)
	go f.blockTaskDispatcher.keepDispatchingTask(ctx, &wg, conf.Solana.StartSlot, 3, f.blockGetTaskCh) // TODO config

	return &wg, f.flowController
}

func (f *Factory) destroy() {
	panic("impl")
}
