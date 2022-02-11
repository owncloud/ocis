package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"NATS_HTTP_ADDR"`
	Namespace string
	Root      string `ocisConfig:"root" env:"NATS_HTTP_ROOT"`
	CacheTTL  int    `ocisConfig:"cache_ttl" env:"NATS_CACHE_TTL"`
}
