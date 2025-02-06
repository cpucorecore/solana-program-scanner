package config

type RedisConf struct {
	Addr     string
	Username string
	Password string
}

var defaultRedisConf = &RedisConf{
	Addr:     "localhost:6379",
	Username: "",
	Password: "",
}
