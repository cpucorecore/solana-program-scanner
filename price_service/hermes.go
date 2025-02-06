package price_service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"solana-program-scanner/config"
	"solana-program-scanner/log"
	"solana-program-scanner/redis"
	"solana-program-scanner/utils"
)

type PriceServiceHermes struct {
	conf                      *config.PriceServiceConf
	currentPriceQueryInterval time.Duration
	currentPrice              MutexPrice
	historyPrice              MutexPrice
	currentPriceCache         redis.CachePrice
}

func (ps *PriceServiceHermes) GetPrice(slot uint64, blockTime int64) decimal.Decimal {
	if time.Now().Unix()-blockTime > 60 {
		if slot%200 == 0 || ps.historyPrice.Get().IsZero() {
			price := ps.GetPriceByTime(blockTime)
			ps.historyPrice.Set(price)
			return price
		}
		return ps.historyPrice.Get()
	}

	return ps.currentPrice.Get()
}

func (ps *PriceServiceHermes) getUrlByTime(blockTime int64) string {
	blockTimeStr := strconv.FormatInt(blockTime, 10)
	return ps.conf.BaseUrl + blockTimeStr + "?" + ps.conf.Id
}

func (ps *PriceServiceHermes) GetPriceByTime(blockTime int64) decimal.Decimal {
	url := ps.getUrlByTime(blockTime)
	return ps.MustQueryPriceByUrl(url)
}

func (ps *PriceServiceHermes) start() {
	ticker := time.NewTicker(ps.currentPriceQueryInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				price, err := ps.queryPrice(ps.conf.QueryUrl)
				if err == nil {
					ps.savePrice(price)
				}
			}
		}
	}()
}

func (ps *PriceServiceHermes) InitPrice() {
	for {
		price, err := ps.queryPrice(ps.conf.QueryUrl)
		if err == nil {
			ps.savePrice(price)
			return
		}
		log.Logger.Error(fmt.Sprintf("query price(latest) err:%s", err.Error()))
		time.Sleep(1 * time.Second)
	}
}

func (ps *PriceServiceHermes) MustQueryPriceByUrl(url string) decimal.Decimal {
	for {
		price, err := ps.queryPrice(url)
		if err == nil {
			return price
		}
		log.Logger.Error(fmt.Sprintf("query price(history) url[%s] err:%s", url, err.Error()))
		time.Sleep(1 * time.Second)
	}
}

func (ps *PriceServiceHermes) savePrice(price decimal.Decimal) {
	ps.currentPrice.Set(price)
	ps.currentPriceCache.SetPrice(price.String())
}

var (
	priceZero = decimal.New(0, 0)
)

func (ps *PriceServiceHermes) queryPrice(url string) (decimal.Decimal, error) {
	now := time.Now()
	resp, err := http.Get(url)
	timeUsed := time.Since(now)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("query price with url:[%s] err: %s", url, err))
		return priceZero, err
	}

	log.Logger.Info(fmt.Sprintf("query price success, %dms used", timeUsed.Milliseconds()))

	var data ResponseData
	err = json.NewDecoder(resp.Body).Decode(&data)
	resp.Body.Close()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("json decode price data err: %v", err))
		return priceZero, err
	}

	if len(data.Parsed) == 0 {
		log.Logger.Error("query price empty data")
		return priceZero, fmt.Errorf("data.Parsed empty")
	}

	priceData := data.Parsed[0].Price
	expo := priceData.Expo
	price, ok := utils.ToDecimal(priceData.Price, int32(expo))
	if !ok {
		log.Logger.Error(fmt.Sprintf("price service ToDecimal(%s, %d) err: %v", priceData.Price, expo, err))
		return priceZero, err
	}

	log.Logger.Info(fmt.Sprintf("slot:%d, price:%s", data.Parsed[0].Metadata.Slot, price.String()))

	return price, nil
}

var _ PriceService = &PriceServiceHermes{}

func NewPriceServiceHermes(conf *config.PriceServiceConf, redisAddr string) *PriceServiceHermes {
	priceService := &PriceServiceHermes{
		conf:                      conf,
		currentPriceQueryInterval: time.Second * time.Duration(conf.QueryIntervalBySecond),
		currentPriceCache:         redis.NewCachePrice(redisAddr),
	}
	priceService.InitPrice()
	priceService.start()
	return priceService
}
