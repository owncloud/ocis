package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"NOTIFICATIONS_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"NOTIFICATIONS_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"NOTIFICATIONS_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"NOTIFICATIONS_DEBUG_ZPAGES"`
}
