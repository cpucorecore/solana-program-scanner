package config

type TokenGetterConf struct {
	TokenServerBaseUrl string
	MaxRetries         int
}

var defaultTokenGetterConf = &TokenGetterConf{
	TokenServerBaseUrl: "http://localhost:11118/token/",
	MaxRetries:         3,
}
