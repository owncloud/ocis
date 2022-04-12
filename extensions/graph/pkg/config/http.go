package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"GRAPH_HTTP_ADDR"`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"GRAPH_HTTP_ROOT"`
}
