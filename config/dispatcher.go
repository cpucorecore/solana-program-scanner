package config

type GrpcConf struct {
	ID                string
	Target            string
	SendTimeoutByMs   int
	MaxRetry          int
	RetryIntervalByMs int
}

type KafkaConf struct {
	ID                string
	Brokers           []string
	Topic             string
	SendTimeoutByMs   int
	MaxRetry          int
	RetryIntervalByMs int
}

type DispatcherConf struct {
	GrpcOn             bool
	Grpc               *GrpcConf
	KafkaOn            bool
	Kafka              *KafkaConf
	KafkaPoolUpdaterOn bool
	KafkaPoolUpdater   *KafkaConf
}

var defaultDispatcherConf = &DispatcherConf{
	GrpcOn: true,
	Grpc: &GrpcConf{
		ID:                "GRPC-0",
		Target:            "localhost:50001",
		SendTimeoutByMs:   5000,
		MaxRetry:          3,
		RetryIntervalByMs: 100,
	},
	KafkaOn: true,
	Kafka: &KafkaConf{
		ID:                "Kafka-0",
		Brokers:           []string{"localhost:9092", "localhost:9093"},
		Topic:             "block",
		SendTimeoutByMs:   5000,
		MaxRetry:          3,
		RetryIntervalByMs: 100,
	},
	KafkaPoolUpdater: &KafkaConf{
		ID:                "KafkaPoolUpdater-0",
		Brokers:           []string{"localhost:9092", "localhost:9093"},
		Topic:             "sol-pool",
		SendTimeoutByMs:   5000,
		MaxRetry:          3,
		RetryIntervalByMs: 100,
	},
}
