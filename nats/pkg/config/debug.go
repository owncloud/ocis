package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"NATS_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"NATS_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"NATS_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"NATS_DEBUG_ZPAGES"`
}
