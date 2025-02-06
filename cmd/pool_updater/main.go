package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"solana-program-scanner/cache"
	"solana-program-scanner/cache_writer"
	"solana-program-scanner/config"
	"solana-program-scanner/dispatcher/kafka"
	"solana-program-scanner/getter/market_detail_getter"
	"solana-program-scanner/getter/market_getter"
	"solana-program-scanner/log"
	"solana-program-scanner/parser"
	"solana-program-scanner/redis"
	"solana-program-scanner/types"
)

type PoolUpdater struct {
	marketCache               cache.MarketCache
	cache2                    redis.Cache2
	marketGetter              market_getter.Getter
	marketDetailGetterService *market_detail_getter.MarketDetailGetterService
}

func (u *PoolUpdater) processBlock(block *kafka.SendBlockPoolReq) {
	log.Logger.Info(fmt.Sprintf("slot:%d start", block.Block.Block))
	for _, pool := range block.PoolInfo {
		u.processPoolBalance(block, pool)
	}
	log.Logger.Info(fmt.Sprintf("slot:%d done", block.Block.Block))
}

func (u *PoolUpdater) getMarket(marketAddr string) *types.Market {
	market, ok := u.marketCache.GetMarket(marketAddr)
	if ok {
		return market
	}

	log.Logger.Info(fmt.Sprintf("market:%s not int cache", marketAddr))
	market = u.marketGetter.MustGetMarket(marketAddr)
	u.marketCache.SetMarket(marketAddr, market)
	return market
}

func (u *PoolUpdater) processPoolBalance(block *kafka.SendBlockPoolReq, pool *parser.PoolBalance) {
	pair, ok := u.cache2.GetPair2(pool.Market)
	if !ok {
		market := u.getMarket(pool.Market)
		pair = u.marketDetailGetterService.ProcessMarket(market)
	}

	if block.Block.Block > pair.Reserve0.Slot {
		pair.Reserve0.Slot = block.Block.Block
		pair.Reserve0.Value = pool.B0
		pair.Reserve1.Slot = block.Block.Block
		pair.Reserve1.Value = pool.B1
		u.cache2.SetPair2(pair)

		if time.Now().Unix()-block.Block.BlockAt <= 5 {
			u.cache2.Publish(pair)
		}
	}
}

const (
	GroupId = "pool-updator-group"
)

func makeKafkaConfig() *sarama.Config {
	c := sarama.NewConfig()
	c.Net.TLS.Enable = false
	c.Consumer.Return.Errors = false
	c.Consumer.Offsets.AutoCommit.Enable = true
	c.Consumer.Offsets.AutoCommit.Interval = time.Second * 1
	c.Consumer.Offsets.Initial = sarama.OffsetOldest
	c.Consumer.MaxWaitTime = time.Millisecond * 500
	c.Consumer.Group.Session.Timeout = time.Second * 30
	c.Consumer.Group.Heartbeat.Interval = time.Second * 10
	c.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	c.Consumer.Fetch.Max = 1024 * 1024
	c.ClientID = "pool-updater"
	return c
}

func newPoolUpdater() *PoolUpdater {
	cache1 := cache.New(config.G.Redis.Addr, config.G.Cache)
	cache2 := redis.NewCache2(config.G.Redis.Addr)
	cacheWriter := cache_writer.NewCacheWriter(cache2)
	tokenGetter := market_detail_getter.NewTokenGetter(config.G.TokenGetter, cache1)
	accountBalanceGetter := market_detail_getter.NewAccountBalanceGetter(config.G.Solana.RPCEndpointHTTP, config.G.MarketDetailGetter)
	marketDetailGetter := market_detail_getter.NewMarketDetailGetter(tokenGetter, accountBalanceGetter)
	marketDetailGetterService := market_detail_getter.NewMarketDetailGetterService2(
		cacheWriter,
		marketDetailGetter,
		config.G.MsgBroker,
	)

	poolUpdater := &PoolUpdater{
		marketCache:               cache1,
		cache2:                    cache2,
		marketGetter:              market_getter.New(config.G.Solana.RPCEndpointHTTP, config.G.MarketGetter.GetAccountTimeoutByMs),
		marketDetailGetterService: marketDetailGetterService,
	}

	return poolUpdater
}

func main() {
	config.LoadConfig("./config.json")
	log.InitLoggerForTest()

	kafkaConfig := makeKafkaConfig()
	group, err := sarama.NewConsumerGroup(config.G.Dispatcher.KafkaPoolUpdater.Brokers, GroupId, kafkaConfig)
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("creating consumer group err: %v", err))
	}
	defer group.Close()

	poolUpdater := newPoolUpdater()
	dbCommiterWaitGroup := poolUpdater.marketDetailGetterService.StartDBCommiter()

	handler := &ConsumerGroupHandler{poolUpdater: poolUpdater}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := <-sigChan
		log.Logger.Info(fmt.Sprintf("Received signal: %v", sig))
		cancel()
	}()

	for {
		err = group.Consume(ctx, []string{config.G.Dispatcher.KafkaPoolUpdater.Topic}, handler)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("Error from consumer: %v", err))
		}
		if ctx.Err() != nil {
			log.Logger.Error(fmt.Sprintf("Error from consumer: %v, do marketDetailGetterService.CloseDBInput()", ctx.Err()))
			poolUpdater.marketDetailGetterService.CloseDBInput()
			break
		}
	}

	dbCommiterWaitGroup.Wait()
}

type ConsumerGroupHandler struct {
	poolUpdater *PoolUpdater
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Logger.Warn(fmt.Sprintf("claim.Messages() done"))
				return nil
			}

			var block kafka.SendBlockPoolReq
			if err := json.Unmarshal(msg.Value, &block); err != nil {
				log.Logger.Fatal(fmt.Sprintf("Failed to unmarshal message: %v", err))
			}

			h.poolUpdater.processBlock(&block)
			session.MarkMessage(msg, "")
		case <-session.Context().Done():
			log.Logger.Info(fmt.Sprintf("session.Context().Done: %v", session.Context().Err()))
			return nil
		}
	}
}
