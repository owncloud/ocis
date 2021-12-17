package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"GRAPH_EXPLORER_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"GRAPH_EXPLORER_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"GRAPH_EXPLORER_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"GRAPH_EXPLORER_DEBUG_ZPAGES"`
}
