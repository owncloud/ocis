package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"STORE_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"STORE_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"STORE_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"STORE_DEBUG_ZPAGES"`
}
