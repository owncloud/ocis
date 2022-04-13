package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"WEB_HTTP_ADDR"`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"WEB_HTTP_ROOT"`
	CacheTTL  int    `yaml:"cache_ttl" env:"WEB_CACHE_TTL"`
}
