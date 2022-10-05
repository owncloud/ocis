package config

// CacheStore defines the available configuration for the cache store
type CacheStore struct {
	Type    string `yaml:"type" env:"OCIS_CACHE_STORE_TYPE;GRAPH_CACHE_STORE_TYPE" desc:"The type of the cache store. Valid options are \"noop\", \"ocmem\", \"etcd\" and \"memory\""`
	Address string `yaml:"address" env:"OCIS_CACHE_STORE_ADDRESS;GRAPH_CACHE_STORE_ADDRESS" desc:"A comma-separated list of addresses to connect to. Only valid if the above setting is set to \"etcd\""`
	Size    int    `yaml:"size" env:"OCIS_CACHE_STORE_SIZE;GRAPH_CACHE_STORE_SIZE" desc:"Maximum number of items per table in the ocmem cache store. Other cache stores will ignore the option and can grow indefinitely."`
}
