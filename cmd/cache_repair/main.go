package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"xorm.io/xorm"

	"solana-program-scanner/config"
	"solana-program-scanner/log"
	spsredis "solana-program-scanner/redis"
	"solana-program-scanner/types"
	"solana-program-scanner/types/orms"
)

const (
	RedisKeyHashTableToken = "smt:token:900"
)

func getAllTokenAddressesFromRedis() []string {
	redisClient := redis.NewClient(&redis.Options{Addr: config.G.Redis.Addr})
	defer redisClient.Close()

	tokenAddresses, err := redisClient.HKeys(context.Background(), RedisKeyHashTableToken).Result()
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("HKeys err: %v", err))
	}
	log.Logger.Info(fmt.Sprintf("tokenAddresses len:%d", len(tokenAddresses)))
	return tokenAddresses
}

type TokenProcessor struct {
	cache    spsredis.Client
	cache2   spsredis.Cache2
	dbEngine *xorm.Engine
	_pair    *orms.Pair
	_token   *orms.Token
}

func NewTokenProcessor() *TokenProcessor {
	dbEngine, err := xorm.NewEngine("postgres", config.G.DBCommiter.PostgresTokenPair.Datasource())
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("connect db err:%v", err))
	}

	return &TokenProcessor{
		cache:    spsredis.New(config.G.Redis.Addr),
		cache2:   spsredis.NewCache2(config.G.Redis.Addr),
		dbEngine: dbEngine,
		_pair:    &orms.Pair{},
		_token:   &orms.Token{},
	}
}

func (p *TokenProcessor) deleteMarket(marketAddress string) {
	err := p.cache2.DelPair2(marketAddress)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Redis delete pair(%s) err:%v", marketAddress, err))
	}

	affected, err := p.dbEngine.Where("address = ?", marketAddress).Delete(p._pair)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("DB Delete pair(%s) err: %v", marketAddress, err))
	} else {
		log.Logger.Info(fmt.Sprintf("%s:%d delete", marketAddress, affected))
	}
}

func (p *TokenProcessor) deleteToken(tokenAddress string) {
	err := p.cache2.DelToken2(tokenAddress)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Redis delete token(%s) err:%v", tokenAddress, err))
	}

	affected, err := p.dbEngine.Where("address = ?", tokenAddress).Delete(p._token)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("DB Delete token(%s) err: %v", tokenAddress, err))
	} else {
		log.Logger.Info(fmt.Sprintf("%s:%d delete", tokenAddress, affected))
	}
}

func isWrongMarket(market *types.Market) bool {
	return !types.IsTokenSol(market.QuoteMint) && !types.IsTokenSol(market.BaseMint)
}

func (p *TokenProcessor) processToken(tokenAddress string) {
	token, err := p.cache2.GetToken2(tokenAddress)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("GetToken2(%s) err: %v", tokenAddress, err))
		return
	}

	tokenUpdated := false
	newMarketAddresses := make([]string, 0, len(token.MarketAddresses))
	for _, marketAddress := range token.MarketAddresses {
		market, ok := p.cache.GetMarket(marketAddress)
		if !ok {
			log.Logger.Error(fmt.Sprintf("GetMarket(%s) err: %v", marketAddress, err))
			continue
		}

		if isWrongMarket(market) {
			log.Logger.Info(fmt.Sprintf("%s:%s can be delete", tokenAddress, marketAddress))
			p.deleteMarket(marketAddress)
			tokenUpdated = true
		} else {
			newMarketAddresses = append(newMarketAddresses, market.Address)
		}
	}

	if len(newMarketAddresses) == 0 {
		log.Logger.Info(fmt.Sprintf("token:%s have no wSOL market, can be delete", tokenAddress))
		p.deleteToken(tokenAddress)
	} else {
		if tokenUpdated {
			token.MarketAddresses = newMarketAddresses
			p.cache2.SetToken2(token)
		}
	}
}

func (p *TokenProcessor) Shutdown() {
	p.cache.Close()
	p.cache2.Close()
	p.dbEngine.Close()
}

func main() {
	config.LoadConfig("./config.json")
	log.InitLoggerForTest()

	tokenProcessor := NewTokenProcessor()
	defer tokenProcessor.Shutdown()

	tokenAddresses := getAllTokenAddressesFromRedis()
	for index, tokenAddress := range tokenAddresses {
		tokenProcessor.processToken(tokenAddress)
		if index%10000 == 0 {
			log.Logger.Info(fmt.Sprintf("token index:%d", index))
		}
	}
}
