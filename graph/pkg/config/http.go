package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"GRAPH_HTTP_ADDR"`
	Namespace string
	Root      string `ocisConfig:"root" env:"GRAPH_HTTP_ROOT"`
}
