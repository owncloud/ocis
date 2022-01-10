package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"GATEWAY_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"GATEWAY_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"GATEWAY_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"GATEWAY_DEBUG_ZPAGES"`
}
