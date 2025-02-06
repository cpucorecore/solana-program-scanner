package config

type LogConf struct {
	Async                      bool
	AsyncBufferSizeByByte      int
	AsyncFlushIntervalBySecond int
}

var defaultLogConf = &LogConf{
	Async:                      false,
	AsyncBufferSizeByByte:      1024 * 1024 * 10,
	AsyncFlushIntervalBySecond: 1,
}
