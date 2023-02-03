package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Log   *Log  `yaml:"log"`
	Debug Debug `yaml:"debug"`

	Nats Nats `ociConfig:"nats"`

	Context context.Context `yaml:"-"`
}

// Nats is the nats config
type Nats struct {
	Host                    string `yaml:"host" env:"NATS_NATS_HOST" desc:"Bind address."`
	Port                    int    `yaml:"port" env:"NATS_NATS_PORT" desc:"Bind port."`
	ClusterID               string `yaml:"clusterid" env:"NATS_NATS_CLUSTER_ID" desc:"ID of the NATS cluster."`
	StoreDir                string `yaml:"store_dir" env:"NATS_NATS_STORE_DIR" desc:"The directory where the filesystem storage will store NATS JetStream data. If not definied, the root directory derives from $OCIS_BASE_DATA_PATH:/nats."`
	TLSCert                 string `yaml:"tls_cert" env:"NATS_TLS_CERT" desc:"Path/File name of the TLS server certificate (in PEM format) for the NATS listener. If not definied, the root directory derives from $OCIS_BASE_DATA_PATH:/nats."`
	TLSKey                  string `yaml:"tls_key" env:"NATS_TLS_KEY" desc:"Path/File name for the TLS certificate key (in PEM format) for the NATS listener. If not definied, the root directory derives from $OCIS_BASE_DATA_PATH:/nats."`
	TLSSkipVerifyClientCert bool   `yaml:"tls_skip_verify_client_cert" env:"OCIS_INSECURE;NATS_TLS_SKIP_VERIFY_CLIENT_CERT" desc:"Whether the NATS server should skip the client certificate verification during the TLS handshake."`
	EnableTLS               bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;NATS_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.."`
}
