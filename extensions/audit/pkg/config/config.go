package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Log   *Log  `yaml:"log,omitempty"`
	Debug Debug `yaml:"debug,omitempty"`

	Events   Events   `yaml:"events,omitempty"`
	Auditlog Auditlog `yaml:"auditlog,omitempty"`

	Context context.Context `yaml:"-"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint      string `yaml:"events_endpoint" env:"AUDIT_EVENTS_ENDPOINT" desc:"the address of the streaming service"`
	Cluster       string `yaml:"events_cluster" env:"AUDIT_EVENTS_CLUSTER" desc:"the clusterID of the streaming service. Mandatory when using nats"`
	ConsumerGroup string `yaml:"events_group" env:"AUDIT_EVENTS_GROUP" desc:"the customergroup of the service. One group will only get one vopy of an event"`
}

// Auditlog holds audit log information
type Auditlog struct {
	LogToConsole bool   `yaml:"log_to_console" env:"AUDIT_LOG_TO_CONSOLE" desc:"logs to Stdout if true"`
	LogToFile    bool   `yaml:"log_to_file" env:"AUDIT_LOG_TO_FILE" desc:"logs to file if true"`
	FilePath     string `yaml:"filepath" env:"AUDIT_FILEPATH" desc:"filepath to the logfile. Mandatory if LogToFile is true"`
	Format       string `yaml:"format" env:"AUDIT_FORMAT" desc:"log format. using json is advised"`
}
