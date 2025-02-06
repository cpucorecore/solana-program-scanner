package config

type BlockTaskDispatcherConf struct {
	Enable              bool
	StartSlot           uint64
	EndSlot             uint64
	GetSlotIntervalByMs int
}

var defaultBlockTaskDispatcherConf = &BlockTaskDispatcherConf{
	Enable:              true,
	StartSlot:           269131631,
	EndSlot:             0,
	GetSlotIntervalByMs: 200,
}
