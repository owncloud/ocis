package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"GRAPH_EXPLORER_HTTP_ADDR"`
	Root      string `yaml:"root" env:"GRAPH_EXPLORER_HTTP_ROOT"`
	Namespace string `yaml:"-"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allowed_credentials"`
}
