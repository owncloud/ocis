package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"THUMBNAILS_HTTP_ADDR"`
	Root      string `yaml:"root" env:"THUMBNAILS_HTTP_ROOT"`
	Namespace string `yaml:"-"`
}
