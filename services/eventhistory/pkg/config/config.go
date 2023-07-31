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
	Addr      string                 `ocisConfig:"addr" env:"EVENTHISTORY_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string                 `ocisConfig:"-" yaml:"-"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
}

// Store configures the store to use
type Store struct {
	Store    string        `yaml:"store" env:"OCIS_PERSISTENT_STORE;EVENTHISTORY_STORE" desc:"The type of the store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	Nodes    []string      `yaml:"nodes" env:"OCIS_PERSISTENT_STORE_NODES;EVENTHISTORY_STORE_NODES" desc:"A comma separated list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store."`
	Database string        `yaml:"database" env:"EVENTHISTORY_STORE_DATABASE" desc:"The database name the configured store should use."`
	Table    string        `yaml:"table" env:"EVENTHISTORY_STORE_TABLE" desc:"The database table the store should use."`
	TTL      time.Duration `yaml:"ttl" env:"OCIS_PERSISTENT_STORE_TTL;EVENTHISTORY_STORE_TTL" desc:"Time to live for events in the store. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '336h' (2 weeks)."`
	Size     int           `yaml:"size" env:"OCIS_PERSISTENT_STORE_SIZE;EVENTHISTORY_STORE_SIZE" desc:"The maximum quantity of items in the store. Only applies when store type 'ocmem' is configured. Defaults to 512."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;EVENTHISTORY_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster              string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;EVENTHISTORY_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;EVENTHISTORY_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE;EVENTHISTORY_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided NOTIFICATIONS_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;EVENTHISTORY_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.."`
}
