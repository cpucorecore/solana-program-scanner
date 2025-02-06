package config

type PriceServiceConf struct {
	QueryUrl              string
	QueryIntervalBySecond int
	BaseUrl               string
	Id                    string
}

var defaultPriceServiceConf = &PriceServiceConf{
	QueryUrl:              "",
	QueryIntervalBySecond: 10,
	BaseUrl:               "",
	Id:                    "",
}
