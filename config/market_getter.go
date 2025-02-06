package config

type MarketGetterConf struct {
	GetAccountIntervalByMs int
	GetAccountTimeoutByMs  int
}

var defaultMarketGetterConf = &MarketGetterConf{
	GetAccountIntervalByMs: 50,
	GetAccountTimeoutByMs:  1000,
}
