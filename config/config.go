package config

import (
	"encoding/json"
	"os"
)

const DefaultConfigFilePath = "config.json"
const DefaultConfigFileTemplatePath = "config.json.template"

type Config struct {
	Log                 *LogConf
	Monitor             *MonitorConf
	Redis               *RedisConf
	Cache               *MemoryCacheConf
	MsgBroker           *MsgBrokerConf
	Solana              *SolanaConf
	PriceService        *PriceServiceConf
	TokenGetter         *TokenGetterConf
	BlockTaskDispatcher *BlockTaskDispatcherConf
	BlockGetterManager  *BlockGetterManagerConf
	DBCommiter          *DBCommiterConf
	Dispatcher          *DispatcherConf
	MarketGetter        *MarketGetterConf
	MarketDetailGetter  *MarketDetailGetterConf
	BlockSequencer      *BlockSequencerConf
}

var F string
var FT string
var G = &Config{
	Log:                 defaultLogConf,
	Monitor:             defaultMonitorConf,
	Redis:               defaultRedisConf,
	Cache:               defaultMemoryCacheConf,
	MsgBroker:           defaultMsgBrokerConf,
	Solana:              defaultSolanaConf,
	PriceService:        defaultPriceServiceConf,
	TokenGetter:         defaultTokenGetterConf,
	BlockTaskDispatcher: defaultBlockTaskDispatcherConf,
	BlockGetterManager:  defaultBlockGetterManagerConf,
	DBCommiter:          defaultDBCommiterConf,
	Dispatcher:          defaultDispatcherConf,
	MarketGetter:        defaultMarketGetterConf,
	MarketDetailGetter:  defaultMarketDetailGetterConf,
	BlockSequencer:      defaultBlockSequencerConf,
}

func LoadConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &G)
	if err != nil {
		return err
	}

	return nil
}

func SaveConfig(filename string) error {
	data, err := json.MarshalIndent(G, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
