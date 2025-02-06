package market_detail_getter

import (
	"fmt"
	"sync"
	"xorm.io/xorm"

	"solana-program-scanner/cache_writer"
	"solana-program-scanner/config"
	"solana-program-scanner/db_commiter"
	"solana-program-scanner/log"
	"solana-program-scanner/msg_broker"
	cache2 "solana-program-scanner/redis/types"
	"solana-program-scanner/types"
	"solana-program-scanner/types/orms"
	"solana-program-scanner/utils"
)

type MarketDetailGetterService struct {
	cacheWriter        *cache_writer.CacheWriter
	marketDetailGetter *MarketDetailGetter
	msgBrokerConf      *config.MsgBrokerConf
	marketBroker       msg_broker.MsgBroker[*types.Market]
	pairBroker         msg_broker.MsgBroker[*orms.Pair]
	tokenBroker        msg_broker.MsgBroker[*orms.Token]
	dbCommiter         *db_commiter.DBCommitter
}

func NewMarketDetailGetterService2(
	cacheWriter *cache_writer.CacheWriter,
	marketDetailGetter *MarketDetailGetter,
	msgBrokerConf *config.MsgBrokerConf,
) *MarketDetailGetterService {
	ormEngine, err := xorm.NewEngine("postgres", config.G.DBCommiter.PostgresTokenPair.Datasource())
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("NewEngine with Datasource(%s) err:%v", config.G.DBCommiter.PostgresTokenPair.Datasource(), err))
	}

	marketBroker := msg_broker.NewMsgBrokerLocal[*types.Market](config.G.MsgBroker.MarketCapacity)
	pairBroker := msg_broker.NewMsgBrokerLocal[*orms.Pair](config.G.MsgBroker.PairCapacity)
	tokenBroker := msg_broker.NewMsgBrokerLocal[*orms.Token](config.G.MsgBroker.TokenCapacity)
	dbCommiter := &db_commiter.DBCommitter{
		PairCommiter:  db_commiter.NewBatchCommitterBuilder[*orms.Pair]("pair", config.G.DBCommiter.MarketConf, pairBroker).Build(ormEngine),
		TokenCommiter: db_commiter.NewBatchCommitterBuilder[*orms.Token]("token", config.G.DBCommiter.TokenConf, tokenBroker).Build(ormEngine),
	}

	return &MarketDetailGetterService{
		cacheWriter:        cacheWriter,
		marketDetailGetter: marketDetailGetter,
		msgBrokerConf:      msgBrokerConf,
		marketBroker:       marketBroker,
		pairBroker:         pairBroker,
		tokenBroker:        tokenBroker,
		dbCommiter:         dbCommiter,
	}
}

func (g *MarketDetailGetterService) cacheToken0(marketAddr string, token *cache2.Token) {
	token0Old, err := g.cacheWriter.GetToken(token.Address)
	if err != nil {
		g.cacheWriter.WriteToken(token)
	} else {
		if !utils.StringSliceContains(token0Old.MarketAddresses, marketAddr) {
			token0Old.MarketAddresses = append(token0Old.MarketAddresses, marketAddr)
			g.cacheWriter.WriteToken(token0Old)
		}
	}
}

func (g *MarketDetailGetterService) ProcessMarket(market *types.Market) *cache2.Pair {
	return g.processMarket(market)
}

func (g *MarketDetailGetterService) processMarket(market *types.Market) *cache2.Pair {
	token0, marketDetail := g.marketDetailGetter.GetMarketDetail(market)

	var totalSupply string
	if token0 != nil {
		token := &orms.Token{
			Address:     token0.Address,
			Name:        token0.Name,
			Decimal:     uint64(token0.Decimals),
			ChainId:     config.G.Solana.ChainId,
			TotalSupply: token0.Supply,
			Symbol:      token0.Symbol,
			Creator:     token0.UpdateAuthority,
			Website:     token0.Uri,
		}
		g.tokenBroker.Produce(token)

		supplyDecimal, ok := utils.ToDecimal(token0.Supply, int32(-token0.Decimals))
		if ok {
			totalSupply = supplyDecimal.String()
		}

		cacheToken := &cache2.Token{
			Address:         token0.Address,
			Name:            token0.Name,
			Decimal:         token0.Decimals,
			ChainId:         config.G.Solana.ChainId,
			TotalSupply:     totalSupply,
			Symbol:          token0.Symbol,
			Website:         token0.Uri,
			MarketAddresses: []string{market.Address},
		}
		g.cacheToken0(market.Address, cacheToken)
	}

	pair := &orms.Pair{
		Address:  market.Address,
		Name:     marketDetail.Name,
		Token0:   marketDetail.Token0,
		Token1:   types.SOL,
		Reserve0: "0",
		Reserve1: "0",
		ChainId:  config.G.Solana.ChainId,
	}

	g.pairBroker.Produce(pair)
	g.cacheWriter.WritePair(marketDetail)

	return marketDetail
}

func (g *MarketDetailGetterService) run(wg *sync.WaitGroup, id int) {
	defer wg.Done()

	marketCh := g.marketBroker.ConsumerChan()
	for {
		market, ok := <-marketCh
		if !ok {
			log.Logger.Info(fmt.Sprintf("MDG:%d marketCh @ done", id))
			return
		}

		g.processMarket(market)
	}
}

func (g *MarketDetailGetterService) RunN(wg *sync.WaitGroup, n int) {
	defer wg.Done()

	wg_ := &sync.WaitGroup{}
	wg_.Add(n)
	for i := 0; i < n; i++ {
		go g.run(wg_, i)
	}
	wg_.Wait()
	g.CloseDBInput()
}

func (g *MarketDetailGetterService) StartDBCommiter() *sync.WaitGroup {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go g.dbCommiter.Run(wg)
	return wg
}

func (g *MarketDetailGetterService) CloseDBInput() {
	g.tokenBroker.Close()
	g.pairBroker.Close()
}
