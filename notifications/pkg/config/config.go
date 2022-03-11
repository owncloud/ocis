package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `ocisConfig:"-" yaml:"-"`

	Service Service `ocisConfig:"-" yaml:"-"`

	Log   *Log  `ocisConfig:"log"`
	Debug Debug `ocisConfig:"debug"`

	Notifications Notifications `ocisConfig:"notifications"`

	Context context.Context `ocisConfig:"-" yaml:"-"`
}

// Notifications definces the config options for the notifications service.
type Notifications struct {
	SMTP              SMTP   `ocisConfig:"SMTP"`
	Events            Events `ocisConfig:"events"`
	RevaGateway       string `ocisConfig:"reva_gateway" env:"REVA_GATEWAY;NOTIFICATIONS_REVA_GATEWAY"`
	MachineAuthSecret string `ocisConfig:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;NOTIFICATIONS_MACHINE_AUTH_API_KEY"`
}

// SMTP combines the smtp configuration options.
type SMTP struct {
	Host     string `ocisConfig:"smtp_host" env:"NOTIFICATIONS_SMTP_HOST"`
	Port     string `ocisConfig:"smtp_port" env:"NOTIFICATIONS_SMTP_PORT"`
	Sender   string `ocisConfig:"smtp_sender" env:"NOTIFICATIONS_SMTP_SENDER"`
	Password string `ocisConfig:"smtp_password" env:"NOTIFICATIONS_SMTP_PASSWORD"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint      string `ocisConfig:"events_endpoint" env:"NOTIFICATIONS_EVENTS_ENDPOINT"`
	Cluster       string `ocisConfig:"events_cluster" env:"NOTIFICATIONS_EVENTS_CLUSTER"`
	ConsumerGroup string `ocisConfig:"events_group" env:"NOTIFICATIONS_EVENTS_GROUP"`
}
