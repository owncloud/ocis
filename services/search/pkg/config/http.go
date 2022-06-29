package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"SEARCH_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `ocisConfig:"-" yaml:"-"`
	Root      string `ocisConfig:"root" env:"SEARCH_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
}
