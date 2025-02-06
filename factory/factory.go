package factory

import (
	"fmt"
	"sync"
	"time"
	"xorm.io/xorm"

	"solana-program-scanner/block_task_dispatcher"
	"solana-program-scanner/cache"
	"solana-program-scanner/config"
	"solana-program-scanner/dispatcher"
	"solana-program-scanner/dispatcher/grpc"
	"solana-program-scanner/dispatcher/kafka"
	"solana-program-scanner/getter/block_getter"
	"solana-program-scanner/getter/market_getter"
	"solana-program-scanner/getter/rate_limiter"
	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/msg_broker"
	"solana-program-scanner/parser"
	"solana-program-scanner/price_service"
	"solana-program-scanner/sequencer"
	"solana-program-scanner/types"
)

type Factory struct {
	cache          cache.Cache
	rpcRateLimiter rate_limiter.RpcRateLimiter

	slotBroker        msg_broker.MsgBroker[uint64]
	blockParsedBroker msg_broker.MsgBroker[*parser.Block]
	marketBroker      msg_broker.MsgBroker[*types.Market]

	taskDispatcher         *block_task_dispatcher.TaskDispatcher
	blockGetterManager     *block_getter.BlockGetterManager
	blockParser            *parser.BlockParser
	grpcClient             *grpc.Client
	kafkaClient            *kafka.Client
	kafkaPoolUpdaterClient *kafka.Client
	dispatcher             *dispatcher.BlockDispatcher
}

func NewPostgresWithConnectPool() {
	txEngine, err := xorm.NewEngine("postgres", "TODO datasource")
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("txEngine xorm.NewEngine err:%v", err))
	}
	txEngine.SetMaxIdleConns(10)
	txEngine.SetMaxOpenConns(20)
	txEngine.SetConnMaxLifetime(time.Minute * 10)
}

func (f *Factory) MakeBrokers() {
	if config.G.MsgBroker.SlotMqOn {
		f.slotBroker = msg_broker.NewMsgBrokerMq[uint64](config.G.MsgBroker.SlotMq)
	} else {
		f.slotBroker = msg_broker.NewMsgBrokerLocal[uint64](config.G.MsgBroker.SlotLocalCapacity)
	}
	f.blockParsedBroker = msg_broker.NewMsgBrokerLocal[*parser.Block](config.G.MsgBroker.BlockParsedCapacity)
	f.marketBroker = msg_broker.NewMsgBrokerLocal[*types.Market](config.G.MsgBroker.MarketCapacity)
}

func (f *Factory) NewCache() {
	f.cache = cache.New(config.G.Redis.Addr, config.G.Cache)
}

func (f *Factory) newFlowController() {
	f.rpcRateLimiter = rate_limiter.New(config.G.BlockGetterManager.RpcRateLimiter)
}

func (f *Factory) newBlockTaskDispatcher() {
	f.taskDispatcher = block_task_dispatcher.New(config.G.BlockTaskDispatcher, config.G.Solana.RPCEndpointHTTP, f.slotBroker)
}

func (f *Factory) newBlockParser() {
	marketGetter := market_getter.New(config.G.Solana.RPCEndpointHTTP, config.G.MarketGetter)
	marketGetter.SubGetMarketErr(f.rpcRateLimiter)

	priceService := price_service.NewPriceServiceHermes(config.G.PriceService, config.G.Redis.Addr)
	raydiumAmmParser := parser.NewRaydiumAmmParser(
		f.cache,
		marketGetter,
	)

	f.blockParser = parser.NewBlockParser(priceService, raydiumAmmParser)
}

func (f *Factory) newBlockGetterManager() {
	blockSequencer := sequencer.NewBlockSequencer(config.G.BlockSequencer, f.blockParsedBroker)
	config := &block_getter.BlockGetterManagerConfig{
		Conf:            config.G.BlockGetterManager,
		RPCEndpointHTTP: config.G.Solana.RPCEndpointHTTP,
		StartSlot:       config.G.BlockTaskDispatcher.StartSlot,
		BlockParser:     f.blockParser,
		BlockSequencer:  blockSequencer,
		RpcRateLimiter:  f.rpcRateLimiter,
		SlotBroker:      f.slotBroker,
	}
	f.blockGetterManager = block_getter.NewBlockGetterManager(config)
}

func (f *Factory) newDispatcher() {
	f.dispatcher = dispatcher.NewBlockDispatcher(f.blockParsedBroker)
	if config.G.Dispatcher.GrpcOn {
		f.grpcClient = grpc.NewClient(config.G.Dispatcher.Grpc)
		f.dispatcher.RegisterDispatcher(f.grpcClient)
	}
	if config.G.Dispatcher.KafkaOn {
		f.kafkaClient = kafka.NewClient(config.G.Dispatcher.Kafka).WithMsgPackager(kafka.ConvertToKafkaReq)
		f.dispatcher.RegisterDispatcher(f.kafkaClient)
	}
	if config.G.Dispatcher.KafkaPoolUpdaterOn {
		f.kafkaPoolUpdaterClient = kafka.NewClient(config.G.Dispatcher.KafkaPoolUpdater).WithMsgPackager(kafka.ConvertToPoolReq)
		f.dispatcher.RegisterDispatcher(f.kafkaPoolUpdaterClient)
	}
}

func (f *Factory) Assemble() *Factory {
	log.InitLogger()

	f.MakeBrokers()

	f.NewCache()
	f.newFlowController()

	f.newBlockTaskDispatcher()
	f.newBlockParser()
	f.newBlockGetterManager()
	f.newDispatcher()

	return f
}

func (f *Factory) Run() {
	monitor.Start()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go f.dispatcher.Run(wg)
	wg.Add(1)
	go f.blockGetterManager.Run(wg)
	wg.Add(1)
	go f.taskDispatcher.Run(wg)

	wg.Wait()
	f.Stop()
}

func (f *Factory) Stop() {
	// TODO release resources
}
