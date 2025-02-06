package config

import (
	"fmt"
	"net/url"
)

type PostgresConf struct {
	Host     string
	Port     int
	Username string
	Password string
	Db       string
}

func (pc *PostgresConf) Datasource() string {
	passwordEncoded := url.QueryEscape(pc.Password)
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		pc.Username,
		passwordEncoded,
		pc.Host,
		pc.Port,
		pc.Db,
	)
}

type BatchCommitConf struct {
	BatchSize        int
	FlushTimeoutByMs int
	BufferSize       int
	WorkerCount      int
}

type DBCommiterConf struct {
	PostgresTokenPair *PostgresConf
	MarketConf        *BatchCommitConf
	TokenConf         *BatchCommitConf
}

var (
	defaultBatchConfig = &BatchCommitConf{
		BatchSize:        100,
		FlushTimeoutByMs: 500,
		BufferSize:       1000,
		WorkerCount:      1,
	}

	defaultDBCommiterConf = &DBCommiterConf{
		PostgresTokenPair: &PostgresConf{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "postgres",
			Db:       "postgres",
		},
		MarketConf: defaultBatchConfig,
		TokenConf:  defaultBatchConfig,
	}
)
