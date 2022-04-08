package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"SEARCH_HTTP_ADDR"`
	Namespace string `ocisConfig:"-" yaml:"-"`
	Root      string `ocisConfig:"root" env:"SEARCH_HTTP_ROOT"`
}
