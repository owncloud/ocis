package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `ocisConfig:"-" yaml:"-"`

	Service Service `ocisConfig:"-" yaml:"-"`

	Tracing *Tracing `ocisConfig:"tracing"`
	Log     *Log     `ocisConfig:"log"`
	Debug   Debug    `ocisConfig:"debug"`

	GRPC GRPC `ocisConfig:"grpc"`

	Datapath string `yaml:"data_path" env:"SEARCH_DATA_PATH"`
	Reva     Reva   `ocisConfig:"reva"`
	Events   Events `yaml:"events"`

	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;SEARCH_MACHINE_AUTH_API_KEY"`

	Context context.Context `ocisConfig:"-" yaml:"-"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint      string `yaml:"events_endpoint" env:"SEARCH_EVENTS_ENDPOINT" desc:"the address of the streaming service"`
	Cluster       string `yaml:"events_cluster" env:"SEARCH_EVENTS_CLUSTER" desc:"the clusterID of the streaming service. Mandatory when using nats"`
	ConsumerGroup string `yaml:"events_group" env:"SEARCH_EVENTS_GROUP" desc:"the customergroup of the service. One group will only get one copy of an event"`
}
