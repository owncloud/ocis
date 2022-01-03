package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"ACCOUNTS_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"ACCOUNTS_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"ACCOUNTS_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"ACCOUNTS_DEBUG_ZPAGES"`
}
