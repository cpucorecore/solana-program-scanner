package monitor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"solana-program-scanner/config"
)

var (
	BlockTxCount = prometheus.NewGauge(prometheus.GaugeOpts{Name: "block_tx_count"})
	CurrentSlot  = prometheus.NewGauge(prometheus.GaugeOpts{Name: "current_slot"})
	NewestSlot   = prometheus.NewGauge(prometheus.GaugeOpts{Name: "newest_slot"})

	GetBlockDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_block_duration_ms",
		Help:       "get_block duration in milliseconds",
		MaxAge:     time.Minute,
		AgeBuckets: 10,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	GetBlockFailedDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_block_failed_duration_ms",
		Help:       "get_block failed duration in milliseconds",
		MaxAge:     time.Minute,
		AgeBuckets: 10,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	GetBlockErrCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "get_block_err_counter",
			Help: "count the get_block err",
		},
		[]string{"err_code"},
	)

	ParseBlockDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "parse_block_duration_ms",
		Help:       "parse block duration in milliseconds",
		MaxAge:     time.Minute,
		AgeBuckets: 10,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	SendBlockGrpcDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "send_block_grpc_duration_ms",
		Help:       "send block grpc duration in milliseconds",
		MaxAge:     time.Minute,
		AgeBuckets: 10,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	SendBlockKafkaDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "send_block_kafka_duration_ms",
		Help:       "send block kafka duration in milliseconds",
		MaxAge:     time.Minute,
		AgeBuckets: 10,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	GetTokenDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_token_duration_ms",
		Help:       "get_token duration in milliseconds",
		MaxAge:     time.Minute,
		AgeBuckets: 10,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	DBBatchCommitDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "db_batch_commit_duration_ms",
			Help:       "database batch_commit duration in milliseconds",
			MaxAge:     time.Minute,
			AgeBuckets: 10,
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{"committer_id"},
	)

	DBCommitOneByOneDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "db_one_by_one_commit_duration_ms",
			Help:       "database OneByOne commit duration in milliseconds",
			MaxAge:     time.Minute,
			AgeBuckets: 10,
			Objectives: map[float64]float64{
				0.5:  0.05,
				0.9:  0.01,
				0.99: 0.001,
			},
		},
		[]string{"committer_id"},
	)

	DBCommitErrCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_commit_err",
			Help: "database commit err count",
		},
		[]string{"committer_id", "err_code"},
	)

	IxTypeCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ix_type_counter",
			Help: "count the ix type",
		},
		[]string{"ix_type"},
	)

	IxSrcCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ix_src_counter",
			Help: "count the ix source",
		},
		[]string{"ix_src"},
	)

	GetAccountDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "get_account_duration_ms",
		Help:       "get account duration in milliseconds",
		MaxAge:     time.Minute,
		AgeBuckets: 10,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})

	BlockDelay = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "block_delay_ms",
		Help:       "block delay in milliseconds",
		MaxAge:     time.Minute,
		AgeBuckets: 10,
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.9:  0.01,
			0.99: 0.001,
		},
	})
)

func init() {
	prometheus.MustRegister(BlockTxCount)
	prometheus.MustRegister(CurrentSlot)
	prometheus.MustRegister(NewestSlot)

	prometheus.MustRegister(GetBlockDuration)
	prometheus.MustRegister(GetBlockFailedDuration)
	prometheus.MustRegister(GetBlockErrCounter)

	prometheus.MustRegister(ParseBlockDuration)
	prometheus.MustRegister(SendBlockGrpcDuration)
	prometheus.MustRegister(SendBlockKafkaDuration)

	prometheus.MustRegister(GetTokenDuration)

	prometheus.MustRegister(DBBatchCommitDuration)
	prometheus.MustRegister(DBCommitOneByOneDuration)
	prometheus.MustRegister(DBCommitErrCounter)

	prometheus.MustRegister(IxTypeCounter)
	prometheus.MustRegister(IxSrcCounter)
	prometheus.MustRegister(GetAccountDuration)
	prometheus.MustRegister(BlockDelay)
}

func Start() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(fmt.Sprintf("%s:%d", config.G.Monitor.ListenHost, config.G.Monitor.ListenPort), nil)
	}()
}
