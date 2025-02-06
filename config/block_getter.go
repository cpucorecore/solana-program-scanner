package config

type BlockGetterConf struct {
	EnableSequencer     bool
	BlockBufferSize     int
	GetBlockTimeoutByMs int
}

type RpcRateLimiterConf struct {
	ErrWaitUnitByMs      int
	ErrWaitUnitByMsQuick int
}

type BlockGetterManagerConf struct {
	GetterNumber   int
	BlockGetter    *BlockGetterConf
	RpcRateLimiter *RpcRateLimiterConf
}

var (
	defaultRpcRateLimiterConf = &RpcRateLimiterConf{
		ErrWaitUnitByMs:      500,
		ErrWaitUnitByMsQuick: 50,
	}

	defaultBlockGetterManagerConf = &BlockGetterManagerConf{
		GetterNumber: 20,
		BlockGetter: &BlockGetterConf{
			EnableSequencer:     false,
			BlockBufferSize:     10,
			GetBlockTimeoutByMs: 5000,
		},
		RpcRateLimiter: defaultRpcRateLimiterConf,
	}
)
