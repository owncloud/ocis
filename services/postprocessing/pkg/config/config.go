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

// Postprocessing defines the config options for the postprocessing service.
type Postprocessing struct {
	Events          Events        `yaml:"events"`
	Steps           []string      `yaml:"steps" env:"POSTPROCESSING_STEPS" desc:"A comma separated list of postprocessing steps, processed in order of their appearance. Currently supported values by the system are: 'virusscan', 'policies' and 'delay'. Custom steps are allowed. See the documentation for instructions."`
	Virusscan       bool          `yaml:"virusscan" env:"POSTPROCESSING_VIRUSSCAN" desc:"After uploading a file but before making it available for download, virus scanning the file can be enabled. Needs as prerequisite the antivirus service to be enabled and configured." deprecationVersion:"master" removalVersion:"master" deprecationVersion:"3.0" removalVersion:"4.0.0" deprecationInfo:"POSTPROCESSING_VIRUSSCAN is not longer necessary and is replaced by POSTPROCESSING_STEPS which also holds information about the order of steps" deprecationReplacement:"POSTPROCESSING_STEPS"`
	Delayprocessing time.Duration `yaml:"delayprocessing" env:"POSTPROCESSING_DELAY" desc:"After uploading a file but before making it available for download, a delay step can be added. Intended for developing purposes only. The duration can be set as number followed by a unit identifier like s, m or h. If a duration is set but the keyword 'delay' is not explicitely added to 'POSTPROCESSING_STEPS', the delay step will be processed as last step. In such a case, a log entry will be written on service startup to remind the admin about that situation."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;POSTPROCESSING_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster  string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;POSTPROCESSING_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`

	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;POSTPROCESSING_EVENTS_TLS_INSECURE" desc:"Whether the ocis server should skip the client certificate verification during the TLS handshake."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"POSTPROCESSING_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided POSTPROCESSING_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;POSTPROCESSING_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services."`
}
