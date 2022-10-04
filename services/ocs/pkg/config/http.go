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
	AllowedOrigins   []string `yaml:"allowed_origins" env:"OCIS_CORS_ALLOW_ORIGINS;OCS_CORS_ALLOW_ORIGINS" desc:"A comma-separated list of allowed CORS origins. See [Access-Control-Allow-Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin) for more details."`
	AllowedMethods   []string `yaml:"allowed_methods" env:"OCIS_CORS_ALLOW_METHODS;OCS_CORS_ALLOW_METHODS" desc:"A comma-separated list of allowed CORS methods. See [Access-Control-Request-Method](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method) for more details."`
	AllowedHeaders   []string `yaml:"allowed_headers" env:"OCIS_CORS_ALLOW_HEADERS;OCS_CORS_ALLOW_HEADERS" desc:"A comma-separated list of allowed CORS headers. See [Access-Control-Request-Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers) for more details."`
	AllowCredentials bool     `yaml:"allowe_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;OCS_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS. See [Access-Control-Allow-Credentials](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials) for more details."`
}
