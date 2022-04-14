package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"SETTINGS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"SETTINGS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"SETTINGS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"SETTINGS_DEBUG_ZPAGES"`
}
