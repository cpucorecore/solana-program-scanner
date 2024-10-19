package main

import (
	"context"
	"fmt"
	"sync"
	"xorm.io/xorm"

	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"solana-program-scanner/block_height_manager"
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
	ormTxCh        chan *OrmTx

	flowController      FlowController
	blockHeightManager  block_height_manager.BlockHeightManager
	mongoAttendant      *MongoAttendant
	postgresAttendant   *PostgresAttendant
	blockTaskDispatcher *BlockTaskDispatcher
	blockGetter         *GetterBlock
	cacheMarket         Cache[string, *OrmMarket]
	getterMarket        GetterMarket
	parserTxRaydiumAmm  *ParserTxRaydiumAmm
	parserTx            *ParserTx
	blockHandler        *BlockHandler
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

func (f *Factory) assemblePostgresClient(dataSource string) {
	engine, err := xorm.NewEngine(PostgresDbDriver, dataSource)
	if err != nil {
		Logger.Fatal(fmt.Sprintf("xorm.NewEngine err:%v", err))
	}

	f.postgresCli = engine
}

func (f *Factory) assembleMongoClient(uri string) {
	connect, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		f.postgresCli.Close()
		Logger.Fatal(fmt.Sprintf("mongo.Connect err:%v", err))
	}

	f.mongoCli = connect
}

func (f *Factory) assembleSolanaClient() *rpc.RpcClient {
	r := rpc.NewRpcClient(gc.Solana.RpcEndpoint)
	f.solanaClis = append(f.solanaClis, &r)
	return &r
}

func (f *Factory) assemblePipelines() {
	f.blockGetTaskCh = make(chan uint64)
	f.blockCh = make(chan *rpc.GetBlock)
	f.ormTxCh = make(chan *OrmTx, 1000)
}

func (f *Factory) assembleFlowController() {
	f.flowController = NewFlowController(
		gc.FlowController.TpsMax,
		gc.FlowController.TpsCountWindow,
		gc.FlowController.TpsWaitUnit,
		gc.FlowController.ErrWaitUnit,
	)
}

func (f *Factory) assembleBlockHeightManager() {
	f.blockHeightManager = block_height_manager.NewBlockHeightManager()
}

func (f *Factory) assembleMongoAttendant() {
}

func (f *Factory) assemblePostgresAttendant() {
	f.postgresAttendant = NewPostgresAttendant(f.postgresCli)
}

func (f *Factory) assembleBlockTaskDispatcher() {
	f.blockTaskDispatcher = NewBlockTaskDispatcher(f.flowController)
}

func (f *Factory) assembleGetterBlock(blockGetterWorkerNumber int) {
	f.blockGetter = NewGetterBlock(blockGetterWorkerNumber, f.blockHeightManager, f.flowController)
}

func (f *Factory) assembleCacheMarket() {
	f.cacheMarket = NewCacheRedisMarket(f.redisCli)
}

func (f *Factory) assembleMarketGetter() {
	f.getterMarket = f.blockGetter
}

func (f *Factory) assembleParserTxRaydiumAmm() {
	f.parserTxRaydiumAmm = NewParserTxRaydiumAmm(f.ormTxCh, f.getterMarket, f.cacheMarket)
}

func (f *Factory) assembleParserTx() {
	f.parserTx = NewParserTx(f.parserTxRaydiumAmm)
}

func (f *Factory) assembleBlockHandler() {
	f.blockHandler = NewBlockHandler(f.blockCh)

	bpf := NewBlockProcessorFile(DefaultBlocksFilePath)
	bpp := NewBlockProcessorParser(f.parserTx)

	f.blockHandler.registerProcessor(bpf)
	f.blockHandler.registerProcessor(bpp)
}

func (f *Factory) assemble() *Factory {
	redisOptions := f.assembleRedisOptions(gc.Redis.Addr, gc.Redis.Username, gc.Redis.Password)
	f.assembleRedisClient(redisOptions)
	f.assemblePostgresClient(gc.Postgres.Datasource())
	f.assembleMongoClient(gc.Mongo.Datasource())
	f.assembleSolanaClient()
	f.assemblePipelines()
	f.assembleFlowController()
	f.assembleBlockHeightManager()
	f.assembleMongoAttendant()
	f.assemblePostgresAttendant()
	f.assembleBlockTaskDispatcher()
	f.assembleGetterBlock(gc.GetterBlock.WorkerNumber)
	f.assembleCacheMarket()
	f.assembleMarketGetter()
	f.assembleParserTxRaydiumAmm()
	f.assembleParserTx()
	f.assembleBlockHandler()

	startBlockHeight := f.blockGetter.getBlockHeight(gc.GetterBlock.StartSlot)
	f.blockHeightManager.Init(startBlockHeight - 1)

	return f
}

func (f *Factory) runProducts(ctx context.Context) (*sync.WaitGroup, FlowController) {
	go f.flowController.startLog(gc.FlowController.LogInterval)

	var wg sync.WaitGroup

	wg.Add(1)
	go f.postgresAttendant.serveTx(ctx, &wg, f.ormTxCh)

	wg.Add(1)
	go f.blockHandler.keepHandling(ctx, &wg)

	wg.Add(1)
	go f.blockGetter.keepBlockGetting(ctx, &wg, f.blockGetTaskCh, f.blockCh)

	wg.Add(1)
	go f.blockTaskDispatcher.keepDispatchingTask(
		ctx,
		&wg,
		gc.GetterBlock.StartSlot,
		gc.GetterBlock.SlotCount,
		f.blockGetTaskCh,
	)

	return &wg, f.flowController
}

func (f *Factory) Shutdown() {
	f.mongoCli.Disconnect(context.Background())
	f.postgresCli.Close()
	f.redisCli.Close()
}
