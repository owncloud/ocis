package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"IDP_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"IDP_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"IDP_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"IDP_DEBUG_ZPAGES"`
}
