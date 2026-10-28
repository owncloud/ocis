package config

import "time"

// Store configures the store to use
type Store struct {
	Store                string        `yaml:"store" env:"OCIS_PERSISTENT_STORE;COLLABORATION_STORE" desc:"The type of the store. Supported values are: 'memory', 'nats-js-kv', 'redis-sentinel', 'noop'. See the text description for details." introductionVersion:"7.0.0"`
	Nodes                []string      `yaml:"nodes" env:"OCIS_PERSISTENT_STORE_NODES;COLLABORATION_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"7.0.0"`
	Database             string        `yaml:"database" env:"COLLABORATION_STORE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"7.0.0"`
	Table                string        `yaml:"table" env:"COLLABORATION_STORE_TABLE" desc:"The database table the store should use." introductionVersion:"7.0.0"`
	TTL                  time.Duration `yaml:"ttl" env:"OCIS_PERSISTENT_STORE_TTL;COLLABORATION_STORE_TTL" desc:"Time to live for events in the store. Defaults to '30m' (30 minutes). See the Environment Variable Types description for more details." introductionVersion:"7.0.0"`
	AuthUsername         string        `yaml:"username" env:"OCIS_PERSISTENT_STORE_AUTH_USERNAME;COLLABORATION_STORE_AUTH_USERNAME" desc:"The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"7.0.0"`
	AuthPassword         string        `yaml:"password" env:"OCIS_PERSISTENT_STORE_AUTH_PASSWORD;COLLABORATION_STORE_AUTH_PASSWORD" desc:"The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"7.0.0"`
	CredsFile            string        `yaml:"creds_file" env:"OCIS_NATS_CREDS_FILE;COLLABORATION_STORE_CREDS_FILE" desc:"Path to a NATS .creds file with a user JWT and an NKey seed. Use it to authenticate with the store when NATS runs in JWT (operator) mode. When set, the username and the password are not used. The file is read again on each reconnect. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	Creds                string        `mask:"password" yaml:"creds" env:"OCIS_NATS_CREDS;COLLABORATION_STORE_CREDS" desc:"Content of a NATS .creds file with a user JWT and an NKey seed. Use it to authenticate with the store when NATS runs in JWT (operator) mode. A service-specific variable takes precedence over a global variable. In the same scope, the content takes precedence over the path of a .creds file. When set, the username and the password are not used. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	EnableTLS            bool          `yaml:"enable_tls" env:"OCIS_PERSISTENT_STORE_ENABLE_TLS;COLLABORATION_STORE_ENABLE_TLS" desc:"Activate TLS for the connection to the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	TLSInsecure          bool          `yaml:"tls_insecure" env:"OCIS_PERSISTENT_STORE_TLS_INSECURE;COLLABORATION_STORE_TLS_INSECURE" desc:"Disable TLS certificate verification for the store connection. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
	TLSRootCACertificate string        `yaml:"tls_root_ca_certificate" env:"OCIS_PERSISTENT_STORE_TLS_ROOT_CA_CERTIFICATE;COLLABORATION_STORE_TLS_ROOT_CA_CERTIFICATE" desc:"Path to the PEM-encoded root CA certificate for the store TLS connection. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"8.3.0"`
}
