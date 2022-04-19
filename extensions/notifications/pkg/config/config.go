package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Log   *Log  `yaml:"log"`
	Debug Debug `yaml:"debug"`

	Notifications Notifications `yaml:"notifications"`

	ConfigFile string `yaml:"-" env:"NOTIFICATIONS_CONFIG_FILE" desc:"config file to be used by the notifications extension"`

	Context context.Context `yaml:"-"`
}

// Notifications definces the config options for the notifications service.
type Notifications struct {
	SMTP              SMTP   `yaml:"SMTP"`
	Events            Events `yaml:"events"`
	RevaGateway       string `yaml:"reva_gateway" env:"REVA_GATEWAY;NOTIFICATIONS_REVA_GATEWAY"`
	MachineAuthSecret string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;NOTIFICATIONS_MACHINE_AUTH_API_KEY"`
}

// SMTP combines the smtp configuration options.
type SMTP struct {
	Host     string `yaml:"smtp_host" env:"NOTIFICATIONS_SMTP_HOST"`
	Port     string `yaml:"smtp_port" env:"NOTIFICATIONS_SMTP_PORT"`
	Sender   string `yaml:"smtp_sender" env:"NOTIFICATIONS_SMTP_SENDER"`
	Password string `yaml:"smtp_password" env:"NOTIFICATIONS_SMTP_PASSWORD"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint      string `yaml:"events_endpoint" env:"NOTIFICATIONS_EVENTS_ENDPOINT"`
	Cluster       string `yaml:"events_cluster" env:"NOTIFICATIONS_EVENTS_CLUSTER"`
	ConsumerGroup string `yaml:"events_group" env:"NOTIFICATIONS_EVENTS_GROUP"`
}
