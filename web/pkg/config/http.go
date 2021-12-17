package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"WEB_HTTP_ADDR"`
	Namespace string
	Root      string `ocisConfig:"root" env:"WEB_HTTP_ROOT"`
	CacheTTL  int    `ocisConfig:"cache_ttl" env:"WEB_CACHE_TTL"`
}
