package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"LLM_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"8.3.0"`
	Token  string `yaml:"token" env:"LLM_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"8.3.0"`
	Pprof  bool   `yaml:"pprof" env:"LLM_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"8.3.0"`
	Zpages bool   `yaml:"zpages" env:"LLM_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"8.3.0"`
}
