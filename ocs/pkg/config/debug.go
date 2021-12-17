package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"OCS_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"OCS_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"OCS_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"OCS_DEBUG_ZPAGES"`
}
