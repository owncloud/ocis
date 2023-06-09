package config

import "github.com/owncloud/ocis/v2/ocis-pkg/shared"

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string                `yaml:"addr" env:"WEB_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	TLS       shared.HTTPServiceTLS `yaml:"tls"`
	Namespace string                `yaml:"-"`
	Root      string                `yaml:"root" env:"WEB_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	CacheTTL  int                   `yaml:"cache_ttl" env:"WEB_CACHE_TTL" desc:"Cache policy in seconds for ownCloud Web assets."`
	CORS      CORS                  `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allow_origins" env:"OCIS_CORS_ALLOW_ORIGINS;WEB_CORS_ALLOW_ORIGINS" desc:"A comma-separated list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin"`
	AllowedMethods   []string `yaml:"allow_methods" env:"OCIS_CORS_ALLOW_METHODS;WEB_CORS_ALLOW_METHODS" desc:"A comma-separated list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method"`
	AllowedHeaders   []string `yaml:"allow_headers" env:"OCIS_CORS_ALLOW_HEADERS;WEB_CORS_ALLOW_HEADERS" desc:"A blank or comma-separated list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers."`
	AllowCredentials bool     `yaml:"allow_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;WEB_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS. See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials."`
}
