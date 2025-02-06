package market_detail_getter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"solana-program-scanner/cache"
	"solana-program-scanner/config"
	"solana-program-scanner/log"
	"solana-program-scanner/monitor"
	"solana-program-scanner/types"
)

const (
	ErrCodeTokenNotFound = 2
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrMaxRetry      = errors.New("err max retry")
)

type TokenGetter struct {
	conf    *config.TokenGetterConf
	cache   cache.TokenCache
	httpCli *http.Client
}

func NewTokenGetter(
	conf *config.TokenGetterConf,
	cache cache.TokenCache,
) *TokenGetter {
	return &TokenGetter{
		conf:    conf,
		cache:   cache,
		httpCli: &http.Client{}, // TODO config the http client
	}
}

func (g *TokenGetter) getToken(mint string) (*types.Token, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", g.conf.TokenServerBaseUrl, mint), nil) // TODO
	if err != nil {
		return nil, err
	}

	log.Logger.Debug(fmt.Sprintf("getToken[%s] begin", mint))
	now := time.Now()
	resp, err := g.httpCli.Do(req)
	log.Logger.Debug(fmt.Sprintf("getToken[%s] end", mint))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResp types.TokenResp
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return nil, err
	}
	monitor.GetTokenDuration.Observe(float64(time.Since(now).Milliseconds()))

	if tokenResp.Status != 0 {
		if ErrCodeTokenNotFound == tokenResp.ErrCode {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf(tokenResp.ErrMsg)
	}

	return &tokenResp.Token, nil
}

func (g *TokenGetter) GetTokenWithRetry(mint string) (*types.Token, error) {
	retryCnt := 0
	for {
		token, err := g.getToken(mint)
		if err == nil {
			return token, nil
		}

		if errors.Is(err, ErrTokenNotFound) {
			log.Logger.Warn(fmt.Sprintf("getToken[%s] err: account not found", mint))
			return nil, ErrTokenNotFound
		}

		retryCnt++
		log.Logger.Error(fmt.Sprintf("getToken[%s] failed %d times with err:%v", mint, retryCnt, err))
		if g.conf.MaxRetries != 0 && retryCnt > g.conf.MaxRetries {
			log.Logger.Error(fmt.Sprintf("getToken[%s] failed with MaxRetryCnt", mint))
			return nil, ErrMaxRetry
		}
		time.Sleep(time.Duration(retryCnt) * time.Second) // TODO config
	}
}
