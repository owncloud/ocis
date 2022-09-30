package config

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allowed_origins" env:"OCIS_CORS_ALLOW_ORIGINS;WEBDAV_CORS_ALLOW_ORIGINS" desc:"Set the allowed CORS origins"`
	AllowedMethods   []string `yaml:"allowed_methods" env:"OCIS_CORS_ALLOW_METHODS;WEBDAV_CORS_ALLOW_METHODS" desc:"Set the allowed CORS methods"`
	AllowedHeaders   []string `yaml:"allowed_headers" env:"OCIS_CORS_ALLOW_HEADERS;WEBDAV_CORS_ALLOW_HEADERS" desc:"Set the allowed CORS headers"`
	AllowCredentials bool     `yaml:"allowed_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;WEBDAV_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"WEBDAV_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"WEBDAV_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	CORS      CORS   `yaml:"cors"`
}
