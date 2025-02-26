package config

import "time"

// Store configures the store to use
type Store struct {
	Store        string        `yaml:"store" env:"OCIS_PERSISTENT_STORE;COLLABORATION_STORE" desc:"The type of the store. Supported values are: 'memory', 'nats-js-kv', 'redis-sentinel', 'noop'. See the text description for details." introductionVersion:"7.0.0"`
	Nodes        []string      `yaml:"nodes" env:"OCIS_PERSISTENT_STORE_NODES;COLLABORATION_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"7.0.0"`
	Database     string        `yaml:"database" env:"COLLABORATION_STORE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"7.0.0"`
	Table        string        `yaml:"table" env:"COLLABORATION_STORE_TABLE" desc:"The database table the store should use." introductionVersion:"7.0.0"`
	TTL          time.Duration `yaml:"ttl" env:"OCIS_PERSISTENT_STORE_TTL;COLLABORATION_STORE_TTL" desc:"Time to live for events in the store. Defaults to '30m' (30 minutes). See the Environment Variable Types description for more details." introductionVersion:"7.0.0"`
	AuthUsername string        `yaml:"username" env:"OCIS_PERSISTENT_STORE_AUTH_USERNAME;COLLABORATION_STORE_AUTH_USERNAME" desc:"The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"7.0.0"`
	AuthPassword string        `yaml:"password" env:"OCIS_PERSISTENT_STORE_AUTH_PASSWORD;COLLABORATION_STORE_AUTH_PASSWORD" desc:"The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"7.0.0"`
}
