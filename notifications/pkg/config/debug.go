package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"NOTIFICATIONS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"NOTIFICATIONS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"NOTIFICATIONS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"NOTIFICATIONS_DEBUG_ZPAGES"`
}
