package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"GRAPH_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"GRAPH_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
}
