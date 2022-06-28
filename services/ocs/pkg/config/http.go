package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"OCS_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Root      string `yaml:"root" env:"OCS_HTTP_ROOT" desc:"The root path of the HTTP service."`
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
