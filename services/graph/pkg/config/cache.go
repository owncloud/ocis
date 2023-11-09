package config

import "time"

// Cache defines the available configuration for a cache store
type Cache struct {
	Store    string        `yaml:"store" env:"OCIS_CACHE_STORE;GRAPH_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	Nodes    []string      `yaml:"nodes" env:"OCIS_CACHE_STORE_NODES;GRAPH_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details."`
	Database string        `yaml:"database" env:"GRAPH_CACHE_STORE_DATABASE" desc:"The database name the configured store should use."`
	Table    string        `yaml:"table" env:"GRAPH_CACHE_STORE_TABLE" desc:"The database table the store should use."`
	TTL      time.Duration `yaml:"ttl" env:"OCIS_CACHE_TTL;GRAPH_CACHE_TTL" desc:"Time to live for cache records in the graph. Defaults to '336h' (2 weeks). See the Environment Variable Types description for more details."`
	Size     int           `yaml:"size" env:"OCIS_CACHE_SIZE;GRAPH_CACHE_SIZE" desc:"The maximum quantity of items in the store. Only applies when store type 'ocmem' is configured. Defaults to 512 which is derived from the ocmem package though not exclicitely set as default."`
}
