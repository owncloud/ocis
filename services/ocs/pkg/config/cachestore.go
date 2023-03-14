package config

import "time"

// CacheStore defines the available configuration for the cache store
type CacheStore struct {
	Type      string        `yaml:"type" env:"OCIS_CACHE_STORE_TYPE;OCS_CACHE_STORE_TYPE" desc:"The type of the cache store. Supported values are: 'mem', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	Addresses []string      `yaml:"addresses" env:"OCIS_CACHE_STORE_ADDRESSES;OCS_CACHE_STORE_ADDRESSES" desc:"A comma separated list of addresses to access the configured store. This has no effect when 'in-memory' stores are configured. Note that the behaviour how addresses are used is dependent on the library of the configured store."`
	Database  string        `yaml:"database" env:"OCS_CACHE_STORE_DATABASE" desc:"The database name the configured store should use."`
	Table     string        `yaml:"table" env:"OCS_CACHE_STORE_TABLE" desc:"The database table the store should use."`
	TTL       time.Duration `yaml:"ttl" env:"OCIS_CACHE_STORE_TTL;OCS_CACHE_STORE_TTL" desc:"Time to live for events in the store. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '336h' (2 weeks)."`
	Size      int           `yaml:"size" env:"OCIS_CACHE_STORE_SIZE;OCS_CACHE_STORE_SIZE" desc:"The maximum quantity of items in the store. Only applies when store type 'ocmem' is configured. Defaults to 512."`
}
