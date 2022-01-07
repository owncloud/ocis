package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"STORE_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"STORE_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"STORE_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"STORE_DEBUG_ZPAGES"`
}
