package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"NATS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"NATS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"NATS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"NATS_DEBUG_ZPAGES"`
}
