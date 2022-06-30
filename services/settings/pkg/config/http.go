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
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allowed_credentials"`
}
