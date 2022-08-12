package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"EXPERIMENTAL_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"EXPERIMENTAL_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
}
