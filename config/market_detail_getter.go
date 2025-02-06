package config

type MarketDetailGetterConf struct {
	GetterNumber                int
	GetBalanceMaxRetry          int
	GetBalanceRetryIntervalByMs int
}

var defaultMarketDetailGetterConf = &MarketDetailGetterConf{
	GetterNumber:                32,
	GetBalanceMaxRetry:          3,
	GetBalanceRetryIntervalByMs: 100,
}
