package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"STORAGE_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"STORAGE_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"STORAGE_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"STORAGE_DEBUG_ZPAGES"`
}
