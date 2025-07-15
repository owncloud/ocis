package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	Service Service         `yaml:"-"`
	Tracing *Tracing        `yaml:"tracing"`
	Log     *Log            `yaml:"log"`
	Debug   Debug           `yaml:"debug"`

	Nats Nats `ociConfig:"nats"`

	Context context.Context `yaml:"-"`
}

// Nats is the nats config
type Nats struct {
	Host                    string `yaml:"host" env:"NATS_NATS_HOST" desc:"Bind address." introductionVersion:"pre5.0"`
	Port                    int    `yaml:"port" env:"NATS_NATS_PORT" desc:"Bind port." introductionVersion:"pre5.0"`
	ClusterID               string `yaml:"clusterid" env:"NATS_NATS_CLUSTER_ID" desc:"ID of the NATS cluster." introductionVersion:"pre5.0"`
	StoreDir                string `yaml:"store_dir" env:"NATS_NATS_STORE_DIR" desc:"The directory where the filesystem storage will store NATS JetStream data. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/nats." introductionVersion:"pre5.0"`
	TLSCert                 string `yaml:"tls_cert" env:"NATS_TLS_CERT" desc:"Path/File name of the TLS server certificate (in PEM format) for the NATS listener. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/nats." introductionVersion:"pre5.0"`
	TLSKey                  string `yaml:"tls_key" env:"NATS_TLS_KEY" desc:"Path/File name for the TLS certificate key (in PEM format) for the NATS listener. If not defined, the root directory derives from $OCIS_BASE_DATA_PATH/nats." introductionVersion:"pre5.0"`
	TLSSkipVerifyClientCert bool   `yaml:"tls_skip_verify_client_cert" env:"OCIS_INSECURE;NATS_TLS_SKIP_VERIFY_CLIENT_CERT" desc:"Whether the NATS server should skip the client certificate verification during the TLS handshake." introductionVersion:"pre5.0"`
	EnableTLS               bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;NATS_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services." introductionVersion:"pre5.0"`
}

// Tracing is the tracing config
type Tracing struct {
	Enabled   bool   `yaml:"enabled" env:"OCIS_TRACING_ENABLED;NATS_TRACING_ENABLED" desc:"Activates tracing." introductionVersion:"pre5.0"`
	Type      string `yaml:"type" env:"OCIS_TRACING_TYPE;NATS_TRACING_TYPE" desc:"The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger' and '' as of now." introductionVersion:"pre5.0"`
	Endpoint  string `yaml:"endpoint" env:"OCIS_TRACING_ENDPOINT;NATS_TRACING_ENDPOINT" desc:"The endpoint of the tracing agent." introductionVersion:"pre5.0"`
	Collector string `yaml:"collector" env:"OCIS_TRACING_COLLECTOR;NATS_TRACING_COLLECTOR" desc:"The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset." introductionVersion:"pre5.0"`
}
