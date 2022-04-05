package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"ACCOUNTS_HTTP_ADDR" desc:"The address of the http service."`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"ACCOUNTS_HTTP_ROOT" desc:"The root path of the http service."`
	CacheTTL  int    `yaml:"cache_ttl" env:"ACCOUNTS_CACHE_TTL" desc:"The cache time for the static assets."`
	CORS      CORS   `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allowed_credentials"`
}
