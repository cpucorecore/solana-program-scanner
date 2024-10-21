package main

import (
	"fmt"
	"time"
)

const (
	DefaultAsyncLog       = false
	DefaultBlocksFilePath = "blocks.json"

	DefaultSolanaRpcEndpoint = "https://api.mainnet-beta.solana.com"

	DefaultRedisAddr     = "127.0.0.1:6379"
	DefaultRedisUsername = ""
	DefaultRedisPassword = ""

	DefaultPostgresHost     = "localhost"
	DefaultPostgresPort     = 5432
	DefaultPostgresUsername = "postgres"
	DefaultPostgresPassword = "12345678"
	DefaultPostgresDbName   = "postgres"
	PostgresDbDriver        = "postgres"

	DefaultMongoHost = "localhost"
	DefaultMongoPort = 27017

	DefaultGetterBlockWorkerNumber = 2
	DefaultGetterBlockStartSlot    = uint64(295503385)
	DefaultGetterBlockSlotCount    = uint64(1000)

	DefaultFlowControllerTpsMax            = 0
	DefaultFlowControllerTpsTpsCountWindow = 5
	DefaultFlowControllerTpsWaitUnit       = time.Second
	DefaultFlowControllerErrWaitUnit       = time.Second
	DefaultFlowControllerLogInterval       = time.Second * 10
)

type SolanaConf struct {
	RpcEndpoint string
}

type RedisConf struct {
	Addr     string
	Username string
	Password string
}

type PostgresConf struct {
	Host     string
	Port     int
	Username string
	Password string
	DbName   string
}

func (pc *PostgresConf) Datasource() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pc.Username,
		pc.Password,
		pc.Host,
		pc.Port,
		pc.DbName,
	)
}

type MongoConf struct {
	Host string
	Port int
}

func (mc *MongoConf) Datasource() string {
	return fmt.Sprintf("mongodb://%s:%d", mc.Host, mc.Port)
}

type GetterBlockConf struct {
	WorkerNumber int
	StartSlot    uint64
	SlotCount    uint64
}

type FlowControllerConf struct {
	TpsMax         uint
	TpsCountWindow int
	TpsWaitUnit    time.Duration
	ErrWaitUnit    time.Duration
	LogInterval    time.Duration
}

type Config struct {
	AsyncLog       bool
	Solana         *SolanaConf
	Redis          *RedisConf
	Postgres       *PostgresConf
	Mongo          *MongoConf
	GetterBlock    *GetterBlockConf
	FlowController *FlowControllerConf
}

var gCfg = &Config{
	AsyncLog: DefaultAsyncLog,
	Solana: &SolanaConf{
		RpcEndpoint: DefaultSolanaRpcEndpoint,
	},
	Redis: &RedisConf{
		Addr:     DefaultRedisAddr,
		Username: DefaultRedisUsername,
		Password: DefaultRedisPassword,
	},
	Postgres: &PostgresConf{
		Host:     DefaultPostgresHost,
		Port:     DefaultPostgresPort,
		Username: DefaultPostgresUsername,
		Password: DefaultPostgresPassword,
		DbName:   DefaultPostgresDbName,
	},
	Mongo: &MongoConf{
		Host: DefaultMongoHost,
		Port: DefaultMongoPort,
	},
	GetterBlock: &GetterBlockConf{
		WorkerNumber: DefaultGetterBlockWorkerNumber,
		StartSlot:    DefaultGetterBlockStartSlot,
		SlotCount:    DefaultGetterBlockSlotCount,
	},
	FlowController: &FlowControllerConf{
		TpsMax:         DefaultFlowControllerTpsMax,
		TpsCountWindow: DefaultFlowControllerTpsTpsCountWindow,
		TpsWaitUnit:    DefaultFlowControllerTpsWaitUnit,
		ErrWaitUnit:    DefaultFlowControllerErrWaitUnit,
		LogInterval:    DefaultFlowControllerLogInterval,
	},
}
