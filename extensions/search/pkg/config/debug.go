package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"SEARCH_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"SEARCH_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"SEARCH_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"SEARCH_DEBUG_ZPAGES"`
}
