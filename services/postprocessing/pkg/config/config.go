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

	Log *Log `yaml:"log"`

	Postprocessing Postprocessing `yaml:"postprocessing"`

	Context context.Context `yaml:"-"`
}

// Postprocessing definces the config options for the postprocessing service.
type Postprocessing struct {
	Events          Events        `yaml:"events"`
	Virusscan       bool          `yaml:"virusscan" env:"POSTPROCESSING_VIRUSSCAN" desc:"After uploading a file but before making it available for download, virus scanning the file can be enabled. Needs as prerequisite the antivirus service to be enabled and configured."`
	Delayprocessing time.Duration `yaml:"delayprocessing" env:"POSTPROCESSING_DELAY" desc:"After uploading a file but before making it available for download, a delay step can be added. Intended for developing purposes only. The duration can be set as number followed by a unit identifier like s, m or h."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint string `yaml:"endpoint" env:"POSTPROCESSING_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster  string `yaml:"cluster" env:"POSTPROCESSING_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`

	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;SEARCH_EVENTS_TLS_INSECURE" desc:"Whether the ocis server should skip the client certificate verification during the TLS handshake."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"SEARCH_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided SEARCH_EVENTS_TLS_INSECURE will be seen as false."`
}
