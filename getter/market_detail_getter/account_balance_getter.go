package market_detail_getter

import (
	"context"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"solana-program-scanner/config"
)

type AccountBalanceGetter struct {
	ctx           context.Context
	conf          *config.MarketDetailGetterConf
	retryInterval time.Duration
	rpcClient     *rpc.Client
}

func NewAccountBalanceGetter(rpcEndpoint string, conf *config.MarketDetailGetterConf) *AccountBalanceGetter {
	client := rpc.New(rpcEndpoint)
	return &AccountBalanceGetter{
		ctx:           context.Background(),
		conf:          conf,
		retryInterval: time.Millisecond * time.Duration(conf.GetBalanceRetryIntervalByMs),
		rpcClient:     client,
	}
}

type Balance struct {
	Slot   uint64
	Amount string
}

func (g *AccountBalanceGetter) GetBalance(addr string) (*Balance, error) {
	account, _ := solana.PublicKeyFromBase58(addr)
	accountBalance, err := g.rpcClient.GetTokenAccountBalance(
		g.ctx,
		account,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		return nil, err
	}
	return &Balance{
		Slot:   accountBalance.Context.Slot,
		Amount: accountBalance.Value.UiAmountString,
	}, nil
}

func (g *AccountBalanceGetter) GetBalanceWithRetry(addr string) (*Balance, error) {
	retryCnt := 0
	for {
		balance, err := g.GetBalance(addr)
		if err == nil {
			return balance, nil
		}

		retryCnt++
		if retryCnt > g.conf.GetBalanceMaxRetry {
			return nil, err
		}

		time.Sleep(g.retryInterval)
	}
}
