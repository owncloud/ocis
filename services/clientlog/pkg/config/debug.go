package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"CLIENTLOG_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"5.0"`
	Token  string `yaml:"token" env:"CLIENTLOG_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"5.0"`
	Pprof  bool   `yaml:"pprof" env:"CLIENTLOG_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"5.0"`
	Zpages bool   `yaml:"zpages" env:"CLIENTLOG_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"5.0"`
}
