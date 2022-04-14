package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"PROXY_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"PROXY_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"PROXY_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"PROXY_DEBUG_ZPAGES"`
}
