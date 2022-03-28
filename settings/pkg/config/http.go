package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"SETTINGS_HTTP_ADDR"`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"SETTINGS_HTTP_ROOT"`
	CacheTTL  int    `yaml:"cache_ttl" env:"SETTINGS_CACHE_TTL"`
	CORS      CORS   `yaml:"cors"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allowed_credentials"`
}
