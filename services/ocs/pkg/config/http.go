package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"OCS_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Root      string `yaml:"root" env:"OCS_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	Namespace string `yaml:"-"`
	CORS      CORS   `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allowed_origins" env:"OCIS_CORS_ALLOW_ORIGINS;OCS_CORS_ALLOW_ORIGINS" desc:"Set the allowed CORS origins"`
	AllowedMethods   []string `yaml:"allowed_methods" env:"OCIS_CORS_ALLOW_METHODS;OCS_CORS_ALLOW_METHODS" desc:"Set the allowed CORS methods"`
	AllowedHeaders   []string `yaml:"allowed_headers" env:"OCIS_CORS_ALLOW_HEADERS;OCS_CORS_ALLOW_HEADERS" desc:"Set the allowed CORS headers"`
	AllowCredentials bool     `yaml:"allowed_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;OCS_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS"`
}
