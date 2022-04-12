package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"GLAUTH_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"GLAUTH_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"GLAUTH_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"GLAUTH_DEBUG_ZPAGES"`
}
