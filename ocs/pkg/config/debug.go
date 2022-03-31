package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"OCS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"OCS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"OCS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"OCS_DEBUG_ZPAGES"`
}
