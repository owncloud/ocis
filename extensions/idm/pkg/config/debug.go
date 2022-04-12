package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"IDM_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"IDM_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"IDM_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"IDM_DEBUG_ZPAGES"`
}
