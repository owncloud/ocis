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
	Virusscan       bool          `yaml:"virusscan" env:"POSTPROCESSING_VIRUSSCAN" desc:"should the system do a virusscan? Needs antivirus service"`
	Delayprocessing time.Duration `yaml:"delayprocessing" env:"POSTPROCESSING_DELAY" desc:"the sytem sleeps for this time while postprocessing"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint string `yaml:"endpoint" env:"POSTPROCESSING_EVENTS_ENDPOINT" desc:"Endpoint of the event system."`
	Cluster  string `yaml:"cluster" env:"POSTPROCESSING_EVENTS_CLUSTER" desc:"Cluster ID of the event system."`

	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;POSTPROCESSING_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"POSTPROCESSING_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided POSTPROCESSING_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;POSTPROCESSING_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.."`
}
