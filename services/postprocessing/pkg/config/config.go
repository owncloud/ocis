package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	Store          Store          `yaml:"store"`
	Postprocessing Postprocessing `yaml:"postprocessing"`

	Context context.Context `yaml:"-"`
}

// Postprocessing defines the config options for the postprocessing service.
type Postprocessing struct {
	Events          Events        `yaml:"events"`
	Steps           []string      `yaml:"steps" env:"POSTPROCESSING_STEPS" desc:"A comma separated list of postprocessing steps, processed in order of their appearance. Currently supported values by the system are: 'virusscan', 'policies' and 'delay'. Custom steps are allowed. See the documentation for instructions."`
	Delayprocessing time.Duration `yaml:"delayprocessing" env:"POSTPROCESSING_DELAY" desc:"After uploading a file but before making it available for download, a delay step can be added. Intended for developing purposes only. The duration can be set as number followed by a unit identifier like s, m or h. If a duration is set but the keyword 'delay' is not explicitely added to 'POSTPROCESSING_STEPS', the delay step will be processed as last step. In such a case, a log entry will be written on service startup to remind the admin about that situation."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;POSTPROCESSING_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster  string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;POSTPROCESSING_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`

	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;POSTPROCESSING_EVENTS_TLS_INSECURE" desc:"Whether the ocis server should skip the client certificate verification during the TLS handshake."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE;POSTPROCESSING_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided POSTPROCESSING_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;POSTPROCESSING_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services."`
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `yaml:"addr" env:"POSTPROCESSING_DEBUG_ADDR" desc:"Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed."`
	Token  string `yaml:"token" env:"POSTPROCESSING_DEBUG_TOKEN" desc:"Token to secure the metrics endpoint."`
	Pprof  bool   `yaml:"pprof" env:"POSTPROCESSING_DEBUG_PPROF" desc:"Enables pprof, which can be used for profiling."`
	Zpages bool   `yaml:"zpages" env:"POSTPROCESSING_DEBUG_ZPAGES" desc:"Enables zpages, which can be used for collecting and viewing in-memory traces."`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;POSTPROCESSING_TRACING_ENABLED" desc:"Activates tracing."`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;POSTPROCESSING_TRACING_TYPE" desc:"The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger' and '' as of now."`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;POSTPROCESSING_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent."`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;POSTPROCESSING_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset."`
}

// Store configures the store to use
type Store struct {
	Store    string        `yaml:"store" env:"OCIS_PERSISTENT_STORE;POSTPROCESSING_STORE" desc:"The type of the store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	Nodes    []string      `yaml:"nodes" env:"OCIS_PERSISTENT_STORE_NODES;POSTPROCESSING_STORE_NODES" desc:"A comma separated list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store."`
	Database string        `yaml:"database" env:"POSTPROCESSING_STORE_DATABASE" desc:"The database name the configured store should use."`
	Table    string        `yaml:"table" env:"POSTPROCESSING_STORE_TABLE" desc:"The database table the store should use."`
	TTL      time.Duration `yaml:"ttl" env:"OCIS_PERSISTENT_STORE_TTL;POSTPROCESSING_STORE_TTL" desc:"Time to live for events in the store. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '336h' (2 weeks)."`
	Size     int           `yaml:"size" env:"OCIS_PERSISTENT_STORE_SIZE;POSTPROCESSING_STORE_SIZE" desc:"The maximum quantity of items in the store. Only applies when store type 'ocmem' is configured. Defaults to 512."`
}
