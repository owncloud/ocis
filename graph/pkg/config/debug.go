package config

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"GRAPH_DEBUG_ADDR"`
	Token  string `yaml:"token" env:"GRAPH_DEBUG_TOKEN"`
	Pprof  bool   `yaml:"pprof" env:"GRAPH_DEBUG_PPROF"`
	Zpages bool   `yaml:"zpages" env:"GRAPH_DEBUG_ZPAGES"`
}
