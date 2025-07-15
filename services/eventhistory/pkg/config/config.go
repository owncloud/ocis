package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"go-micro.dev/v4/client"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	GRPC          GRPCConfig            `yaml:"grpc"`
	GRPCClientTLS *shared.GRPCClientTLS `yaml:"grpc_client_tls"`
	GrpcClient    client.Client         `yaml:"-"`

	Events Events `yaml:"events"`
	Store  Store  `yaml:"store"`

	Context context.Context `yaml:"-"`
}

// GRPCConfig defines the available grpc configuration.
type GRPCConfig struct {
	Addr      string                 `ocisConfig:"addr" env:"EVENTHISTORY_GRPC_ADDR" desc:"The bind address of the GRPC service." introductionVersion:"pre5.0"`
	Namespace string                 `ocisConfig:"-" yaml:"-"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
}

// Store configures the store to use
type Store struct {
	Store        string        `yaml:"store" env:"OCIS_PERSISTENT_STORE;EVENTHISTORY_STORE" desc:"The type of the store. Supported values are: 'memory', 'nats-js-kv', 'redis-sentinel', 'noop'. See the text description for details." introductionVersion:"pre5.0"`
	Nodes        []string      `yaml:"nodes" env:"OCIS_PERSISTENT_STORE_NODES;EVENTHISTORY_STORE_NODES" desc:"A list of nodes to access the configured store. This has no effect when 'memory' store is configured. Note that the behaviour how nodes are used is dependent on the library of the configured store. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	Database     string        `yaml:"database" env:"EVENTHISTORY_STORE_DATABASE" desc:"The database name the configured store should use." introductionVersion:"pre5.0"`
	Table        string        `yaml:"table" env:"EVENTHISTORY_STORE_TABLE" desc:"The database table the store should use." introductionVersion:"pre5.0"`
	TTL          time.Duration `yaml:"ttl" env:"OCIS_PERSISTENT_STORE_TTL;EVENTHISTORY_STORE_TTL" desc:"Time to live for events in the store. Defaults to '336h' (2 weeks). See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AuthUsername string        `yaml:"username" env:"OCIS_PERSISTENT_STORE_AUTH_USERNAME;EVENTHISTORY_STORE_AUTH_USERNAME" desc:"The username to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
	AuthPassword string        `yaml:"password" env:"OCIS_PERSISTENT_STORE_AUTH_PASSWORD;EVENTHISTORY_STORE_AUTH_PASSWORD" desc:"The password to authenticate with the store. Only applies when store type 'nats-js-kv' is configured." introductionVersion:"5.0"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;EVENTHISTORY_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture." introductionVersion:"pre5.0"`
	Cluster              string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;EVENTHISTORY_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system." introductionVersion:"pre5.0"`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;EVENTHISTORY_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates." introductionVersion:"pre5.0"`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE;EVENTHISTORY_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. Will be seen as empty if NOTIFICATIONS_EVENTS_TLS_INSECURE is provided." introductionVersion:"pre5.0"`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;EVENTHISTORY_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"pre5.0"`
	AuthUsername         string `yaml:"username" env:"OCIS_EVENTS_AUTH_USERNAME;EVENTHISTORY_EVENTS_AUTH_USERNAME" desc:"The username to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"5.0"`
	AuthPassword         string `yaml:"password" env:"OCIS_EVENTS_AUTH_PASSWORD;EVENTHISTORY_EVENTS_AUTH_PASSWORD" desc:"The password to authenticate with the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"5.0"`
}
