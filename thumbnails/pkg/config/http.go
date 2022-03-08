package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"THUMBNAILS_HTTP_ADDR"`
	Root      string `ocisConfig:"root" env:"THUMBNAILS_HTTP_ROOT"`
	Namespace string
}
