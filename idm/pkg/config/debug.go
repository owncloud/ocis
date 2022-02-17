package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"IDM_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"IDM_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"IDM_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"IDM_DEBUG_ZPAGES"`
}
