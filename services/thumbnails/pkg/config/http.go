package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"THUMBNAILS_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Root      string `yaml:"root" env:"THUMBNAILS_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	Namespace string `yaml:"-"`
}
