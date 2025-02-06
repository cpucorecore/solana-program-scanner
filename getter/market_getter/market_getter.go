package market_getter

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/mr-tron/base58"

	"solana-program-scanner/config"
	"solana-program-scanner/getter/pubsub"
	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/solana"
	"solana-program-scanner/types"
)

type Getter interface {
	GetMarket(addr string) (market *types.Market, err error)
	MustGetMarket(addr string) *types.Market
	SubGetMarketErr(sub pubsub.GetMarketErrSub)
}

type getter struct {
	id                 string
	rpcUrl             string
	ctx                context.Context
	getAccountInterval time.Duration
	getAccountTimeout  time.Duration
	rpcClient          *rpc.RpcClient
	getMarketErrSubs   []pubsub.GetMarketErrSub
	errCnt             int
}

func New(rpcUrl string, conf *config.MarketGetterConf) Getter {
	return &getter{
		id:                 "MG",
		rpcUrl:             rpcUrl,
		ctx:                context.Background(),
		getAccountInterval: time.Millisecond * time.Duration(conf.GetAccountIntervalByMs),
		getAccountTimeout:  time.Millisecond * time.Duration(conf.GetAccountTimeoutByMs),
		rpcClient:          solana.NewCommonClient(rpcUrl),
	}
}

const (
	AccountDataLengthMarket = 752
	BaseDecimalStartIndex   = 32
	QuoteDecimalStartIndex  = 40
	BaseVaultStartIndex     = 336
	QuoteVaultStartIndex    = 368
	BaseMintStartIndex      = 400
	QuoteMintStartIndex     = 432
)

var (
	getAccountInfoConfig = rpc.GetAccountInfoConfig{
		Commitment: rpc.CommitmentConfirmed,
		Encoding:   rpc.AccountEncodingBase64,
		DataSlice: &rpc.DataSlice{
			Offset: 0,
			Length: AccountDataLengthMarket * 3,
		},
	}
)

func (g *getter) resetRpcClient() {
	log.Logger.Warn(fmt.Sprintf("%s resetRpcClient", g.id))
	rpcClient := rpc.NewRpcClient(g.rpcUrl)
	g.rpcClient = &rpcClient
}

func (g *getter) dealErr() {
	g.errCnt++
	if g.errCnt >= 5 {
		g.resetRpcClient()
		g.errCnt = 0
	}
}

func (g *getter) mustGetAccount(address string) *rpc.AccountInfo {
	ctx, cancel := context.WithTimeout(g.ctx, g.getAccountTimeout)
	defer cancel()

	for {
		now := time.Now()
		resp, err := g.rpcClient.GetAccountInfoWithConfig(ctx, address, getAccountInfoConfig)
		monitor.GetAccountDuration.Observe(float64(time.Since(now).Milliseconds()))
		if err != nil {
			g.dealErr()
			g.pubGetMarketErr(address, types.RpcErr)
			log.Logger.Error(fmt.Sprintf("%s GetAccount:%s err:%s", g.id, address, err.Error()))
			continue
		}

		if resp.Error != nil {
			g.dealErr()
			g.pubGetMarketErr(address, resp.Error.Code)
			log.Logger.Error(fmt.Sprintf("%s GetAccount:%s JsonRpc err:%s", g.id, address, resp.Error))
			continue
		}

		g.errCnt = 0
		return &resp.Result.Value
	}
}

func toUint64(bytesLE []byte) (value uint64) {
	for i := len(bytesLE) - 1; i >= 0; i-- {
		value = value<<8 + uint64(bytesLE[i])
	}
	return
}

const (
	TypeAssertionData  = "Data"
	TypeAssertionData0 = "Data[0]"
)

func debugAccount(ai *rpc.AccountInfo) {
	marshal, err := json.Marshal(ai)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("json Marshal %v err:%s", ai, err.Error()))
	}
	log.Logger.Debug(fmt.Sprintf("accountInfo: [%s]", string(marshal)))
}

func (g *getter) GetMarket(addr string) (*types.Market, error) {
	account := g.mustGetAccount(addr)

	dataArray, ok := account.Data.([]any)
	if !ok {
		debugAccount(account)
		return nil, fmt.Errorf("market:%s type assertion failed on '%s'", addr, TypeAssertionData)
	}

	accountDataStr, ok := dataArray[0].(string)
	if !ok {
		debugAccount(account)
		return nil, fmt.Errorf("market:%s type assertion failed on '%s'", addr, TypeAssertionData0)
	}

	accountDataBytes, err := base64.StdEncoding.DecodeString(accountDataStr)
	if err != nil {
		return nil, fmt.Errorf("market:%s base64 decode:[%s] err:%v", addr, accountDataStr, err)
	}

	if len(accountDataBytes) != AccountDataLengthMarket {
		return nil, fmt.Errorf("data length:%d != expected length:%d, data:[%s]",
			len(accountDataBytes),
			AccountDataLengthMarket,
			hex.EncodeToString(accountDataBytes))
	}

	baseDecimalBytes := accountDataBytes[BaseDecimalStartIndex : BaseDecimalStartIndex+8]
	quoteDecimalBytes := accountDataBytes[QuoteDecimalStartIndex : QuoteDecimalStartIndex+8]
	baseVaultBytes := accountDataBytes[BaseVaultStartIndex : BaseVaultStartIndex+32]
	quoteVaultBytes := accountDataBytes[QuoteVaultStartIndex : QuoteVaultStartIndex+32]
	baseMintBytes := accountDataBytes[BaseMintStartIndex : BaseMintStartIndex+32]
	quoteMintBytes := accountDataBytes[QuoteMintStartIndex : QuoteMintStartIndex+32]

	market := &types.Market{
		Address:      addr,
		BaseDecimal:  int32(toUint64(baseDecimalBytes)),
		QuoteDecimal: int32(toUint64(quoteDecimalBytes)),
		BaseVault:    base58.Encode(baseVaultBytes),
		QuoteVault:   base58.Encode(quoteVaultBytes),
		BaseMint:     base58.Encode(baseMintBytes),
		QuoteMint:    base58.Encode(quoteMintBytes),
	}

	return market, nil
}

func (g *getter) MustGetMarket(addr string) *types.Market {
	cnt := 0
	for {
		market, err := g.GetMarket(addr)
		if err != nil {
			cnt++
			log.Logger.Error(fmt.Sprintf("%s GetMarket:%s err:%s cnt:%d", g.id, addr, err.Error(), cnt))
			time.Sleep(g.getAccountInterval)
			continue
		}
		return market
	}
}

func (g *getter) SubGetMarketErr(sub pubsub.GetMarketErrSub) {
	g.getMarketErrSubs = append(g.getMarketErrSubs, sub)
}

func (g *getter) pubGetMarketErr(marketAddr string, errCode int) {
	for _, sub := range g.getMarketErrSubs {
		sub.PubGetMarketErr(marketAddr, errCode)
	}
}
