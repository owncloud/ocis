package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"OCS_HTTP_ADDR"`
	Root      string `yaml:"root" env:"OCS_HTTP_ROOT"`
	Namespace string `yaml:"-"`
	CORS      CORS   `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allowed_credentials"`
}
