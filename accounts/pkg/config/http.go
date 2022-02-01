package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `env:"ACCOUNTS_HTTP_ADDR"`
	Namespace string
	Root      string `env:"ACCOUNTS_HTTP_ROOT"`
	CacheTTL  int    `env:"ACCOUNTS_CACHE_TTL"`
	CORS      CORS
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}
