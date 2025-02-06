package config

type MonitorConf struct {
	ListenHost string
	ListenPort int
}

var defaultMonitorConf = &MonitorConf{
	ListenHost: "0.0.0.0",
	ListenPort: 18080,
}
