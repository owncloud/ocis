package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `env:"ACCOUNTS_HTTP_ADDR" desc:"The address of the http service."`
	Namespace string
	Root      string `env:"ACCOUNTS_HTTP_ROOT" desc:"The root path of the http service."`
	CacheTTL  int    `env:"ACCOUNTS_CACHE_TTL" desc:"The cache time for the static assets."`
	CORS      CORS
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}
