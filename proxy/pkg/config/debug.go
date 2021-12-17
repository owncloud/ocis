package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"PROXY_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"PROXY_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"PROXY_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"PROXY_DEBUG_ZPAGES"`
}
