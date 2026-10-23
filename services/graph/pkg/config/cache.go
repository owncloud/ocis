package config

import "time"

// Cache defines the available configuration for a cache store
type Cache struct {
	Store                string        `yaml:"store" env:"OCIS_CACHE_STORE;GRAPH_CACHE_STORE" desc:"The type of the cache store. Supported values are: 'memory', 'redis-sentinel', 'nats-js-kv', 'noop'. See the text description for details." introductionVersion:"pre5.0"`
	Nodes                []string      `yaml:"nodes" env:"OCIS_CACHE_STORE_NODES;GRAPH_CACHE_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' store are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	Database             string        `yaml:"database" env:"GRAPH_CACHE_STORE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"pre5.0"`
	Table                string        `yaml:"table" env:"GRAPH_CACHE_STORE_TABLE" desc:"The database table the store should use." introductionVersion:"pre5.0"`
	TTL                  time.Duration `yaml:"ttl" env:"OCIS_CACHE_TTL;GRAPH_CACHE_TTL" desc:"Time to live for cache records in the graph. Defaults to '336h' (2 weeks). See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	DisablePersistence   bool          `yaml:"disable_persistence" env:"OCIS_CACHE_DISABLE_PERSISTENCE;GRAPH_CACHE_DISABLE_PERSISTENCE" desc:"Disables persistence of the cache. Only applies when store type 'nats-js-kv' is configured. Defaults to false." introductionVersion:"5.0"`
	AuthUsername         string        `yaml:"username" env:"OCIS_CACHE_AUTH_USERNAME;GRAPH_CACHE_AUTH_USERNAME" desc:"The username to authenticate with the cache. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
	AuthPassword         string        `yaml:"password" env:"OCIS_CACHE_AUTH_PASSWORD;GRAPH_CACHE_AUTH_PASSWORD" desc:"The password to authenticate with the cache. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
	CredsFile            string        `yaml:"creds_file" env:"OCIS_NATS_CREDS_FILE;GRAPH_CACHE_CREDS_FILE" desc:"Path to a NATS .creds file with a user JWT and an NKey seed. Use it to authenticate with the cache when NATS runs in JWT (operator) mode. When set, the username and the password are not used. The file is read again on each reconnect. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	Creds                string        `mask:"password" yaml:"creds" env:"OCIS_NATS_CREDS;GRAPH_CACHE_CREDS" desc:"Content of a NATS .creds file with a user JWT and an NKey seed. Use it to authenticate with the cache when NATS runs in JWT (operator) mode. A service-specific variable takes precedence over a global variable. In the same scope, the content takes precedence over the path of a .creds file. When set, the username and the password are not used. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	EnableTLS            bool          `yaml:"enable_tls" env:"OCIS_CACHE_ENABLE_TLS;GRAPH_CACHE_ENABLE_TLS" desc:"Activate TLS for the connection to the cache store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	TLSInsecure          bool          `yaml:"tls_insecure" env:"OCIS_CACHE_TLS_INSECURE;GRAPH_CACHE_TLS_INSECURE" desc:"Disable TLS certificate verification for the cache store connection. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	TLSRootCACertificate string        `yaml:"tls_root_ca_certificate" env:"OCIS_CACHE_TLS_ROOT_CA_CERTIFICATE;GRAPH_CACHE_TLS_ROOT_CA_CERTIFICATE" desc:"Path to the PEM-encoded root CA certificate for the cache store TLS connection. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
}
