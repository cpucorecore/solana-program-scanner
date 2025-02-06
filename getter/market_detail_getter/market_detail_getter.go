package market_detail_getter

import (
	"fmt"

	"solana-program-scanner/config"
	"solana-program-scanner/log"
	cache2 "solana-program-scanner/redis/types"
	"solana-program-scanner/types"
)

type MarketDetailGetter struct {
	tokenGetter          *TokenGetter
	accountBalanceGetter *AccountBalanceGetter
}

func NewMarketDetailGetter(tokenGetter *TokenGetter, accountBalanceGetter *AccountBalanceGetter) *MarketDetailGetter {
	return &MarketDetailGetter{
		tokenGetter:          tokenGetter,
		accountBalanceGetter: accountBalanceGetter,
	}
}

type MarketDetail struct {
	token0        *types.Token
	token1        *types.Token
	vault0Balance string
	vault1Balance string
}

func makeMarketName(token0Symbol string) string {
	return fmt.Sprintf("%s/SOL", token0Symbol)
}

func (g *MarketDetailGetter) GetMarketBalances(marketAddr, vault0Addr, vault1Addr string) (vault0Balance, vault1Balance *Balance) {
	vault0Balance, err0 := g.accountBalanceGetter.GetBalanceWithRetry(vault0Addr)
	if err0 != nil {
		log.Logger.Error(fmt.Sprintf("get balance[%s-%s] err:%v", marketAddr, vault0Addr, err0))
	}

	vault1Balance, err1 := g.accountBalanceGetter.GetBalanceWithRetry(vault1Addr)
	if err1 != nil {
		log.Logger.Error(fmt.Sprintf("get balance[%s-%s] err:%v", marketAddr, vault1Addr, err1))
	}

	return
}

func arrangeTokensVaults(market *types.Market) (token0Addr string, token1Addr string, vault0Addr string, vault1Addr string) {
	token0Addr = market.BaseMint
	token1Addr = market.QuoteMint
	vault0Addr = market.BaseVault
	vault1Addr = market.QuoteVault

	if types.IsTokenSol(market.BaseMint) {
		token0Addr = market.QuoteMint
		token1Addr = market.BaseMint
		vault0Addr = market.QuoteVault
		vault1Addr = market.BaseVault
	}

	return
}

func (g *MarketDetailGetter) GetMarketDetail(market *types.Market) (*types.Token, *cache2.Pair) {
	token0Addr, _, vault0Addr, vault1Addr := arrangeTokensVaults(market)

	var marketName string
	token0, err := g.tokenGetter.GetTokenWithRetry(token0Addr)
	if err != nil {
		marketName = makeMarketName(token0Addr)
	} else {
		marketName = makeMarketName(token0.Symbol)
	}

	cachePair := &cache2.Pair{
		Address: market.Address,
		Name:    marketName,
		Token0:  token0Addr,
		Token1:  types.SOL,
		Reserve0: &cache2.Amount{
			Slot:  0,
			Value: "0",
		},
		Reserve1: &cache2.Amount{
			Slot:  0,
			Value: "0",
		},
		ChainId: config.G.Solana.ChainId,
	}

	vault0Balance, vault1Balance := g.GetMarketBalances(market.Address, vault0Addr, vault1Addr)
	if vault0Balance != nil {
		cachePair.Reserve0 = &cache2.Amount{
			Slot:  vault0Balance.Slot,
			Value: vault0Balance.Amount,
		}
	}
	if vault1Balance != nil {
		cachePair.Reserve1 = &cache2.Amount{
			Slot:  vault1Balance.Slot,
			Value: vault1Balance.Amount,
		}
	}

	return token0, cachePair
}
