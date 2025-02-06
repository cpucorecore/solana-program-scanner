package config

const (
	DefaultCacheExpirationByHour      = 24 * 7
	DefaultCacheCleanupIntervalByHour = 24
)

type MemoryCacheConf struct {
	ExpirationByHour      int
	CleanupIntervalByHour int
}

var defaultMemoryCacheConf = &MemoryCacheConf{
	ExpirationByHour:      DefaultCacheExpirationByHour,
	CleanupIntervalByHour: DefaultCacheCleanupIntervalByHour,
}
