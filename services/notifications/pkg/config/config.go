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

	Notifications Notifications `yaml:"notifications"`

	Context context.Context `yaml:"-"`
}

// Notifications definces the config options for the notifications service.
type Notifications struct {
	SMTP              SMTP   `yaml:"SMTP"`
	Events            Events `yaml:"events"`
	RevaGateway       string `yaml:"reva_gateway" env:"REVA_GATEWAY;NOTIFICATIONS_REVA_GATEWAY" desc:"CS3 gateway used to look up user metadata"`
	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;NOTIFICATIONS_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary to access resources from other services."`
}

// SMTP combines the smtp configuration options.
type SMTP struct {
	Host     string `yaml:"smtp_host" env:"NOTIFICATIONS_SMTP_HOST" desc:"SMTP host to connect to."`
	Port     string `yaml:"smtp_port" env:"NOTIFICATIONS_SMTP_PORT" desc:"Port of the SMTP host to connect to."`
	Sender   string `yaml:"smtp_sender" env:"NOTIFICATIONS_SMTP_SENDER" desc:"Sender of emails that will be sent."`
	Password string `yaml:"smtp_password" env:"NOTIFICATIONS_SMTP_PASSWORD" desc:"Password of the SMTP host to connect to."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint      string `yaml:"endpoint" env:"NOTIFICATIONS_EVENTS_ENDPOINT" desc:"Endpoint of the event system."`
	Cluster       string `yaml:"cluster" env:"NOTIFICATIONS_EVENTS_CLUSTER" desc:"Cluster ID of the event system."`
	ConsumerGroup string `yaml:"group" env:"NOTIFICATIONS_EVENTS_GROUP" desc:"Name of the event group / queue on the event system."`
}
