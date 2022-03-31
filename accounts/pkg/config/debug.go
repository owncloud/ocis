package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"ACCOUNTS_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"ACCOUNTS_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"ACCOUNTS_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"ACCOUNTS_DEBUG_ZPAGES"`
}
