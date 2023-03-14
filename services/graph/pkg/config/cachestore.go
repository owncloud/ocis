package config

import "time"

// CacheStore defines the available configuration for a cache store
type CacheStore struct {
	Type      string        `yaml:"type" env:"OCIS_CACHE_STORE_TYPE;GRAPH_CACHE_STORE_TYPE" desc:"The type of the cache store. Supported values are: 'mem', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	Addresses []string      `yaml:"addresses" env:"OCIS_CACHE_STORE_ADDRESSES;GRAPH_CACHE_STORE_ADDRESSES" desc:"A comma-separated list of addresses to connect to. Only valid if the above setting is set to \"etcd\""`
	Database  string        `yaml:"database" env:"GRAPH_CACHE_STORE_DATABASE" desc:"The database name the configured store should use. This has no effect when 'in-memory' stores are configured."`
	Table     string        `yaml:"table" env:"GRAPH_CACHE_STORE_TABLE" desc:"The database table the store should use. This has no effect when 'in-memory' stores are configured."`
	TTL       time.Duration `yaml:"ttl" env:"OCIS_CACHE_STORE_TTL;GRAPH_CACHE_STORE_TTL" desc:"Time to live for cache records in the graph. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '336h' (2 weeks)."`
	Size      int           `yaml:"size" env:"OCIS_CACHE_STORE_SIZE;GRAPH_CACHE_STORE_SIZE" desc:"Maximum number of items per table in the ocmem cache store. Other cache stores will ignore the option and can grow indefinitely."`
}
