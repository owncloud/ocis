package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"SETTINGS_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"SETTINGS_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	CacheTTL  int    `yaml:"cache_ttl" env:"SETTINGS_CACHE_TTL" desc:"Browser cache control max-age value in seconds for settings Web UI assets."`
	CORS      CORS   `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allowed_origins" env:"OCIS_CORS_ALLOW_ORIGINS;SETTINGS_CORS_ALLOW_ORIGINS" desc:"Set the allowed CORS origins"`
	AllowedMethods   []string `yaml:"allowed_methods" env:"OCIS_CORS_ALLOW_METHODS;SETTINGS_CORS_ALLOW_METHODS" desc:"Set the allowed CORS methods"`
	AllowedHeaders   []string `yaml:"allowed_headers" env:"OCIS_CORS_ALLOW_HEADERS;SETTINGS_CORS_ALLOW_HEADERS" desc:"Set the allowed CORS headers"`
	AllowCredentials bool     `yaml:"allowed_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;SETTINGS_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS"`
}
