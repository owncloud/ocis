package config

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"WEB_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"WEB_HTTP_ROOT" desc:"The root path of the HTTP service."`
	CacheTTL  int    `yaml:"cache_ttl" env:"WEB_CACHE_TTL" desc:"Cache policy in seconds for ownCloud Web assets."`
}
