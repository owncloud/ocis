package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"STORE_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	Token  string `yaml:"token" env:"STORE_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	Pprof  bool   `yaml:"pprof" env:"STORE_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
	Zpages bool   `yaml:"zpages" env:"STORE_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces." introductionVersion:"pre5.0" deprecationVersion:"5.0" removalVersion:"7.0.0" deprecationInfo:"The store service is optional and will be removed."`
}
