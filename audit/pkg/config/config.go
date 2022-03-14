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

	Events   Events   `ocisConfig:"events"`
	Auditlog Auditlog `ocisConfig:"auditlog"`

	Context context.Context `ocisConfig:"-" yaml:"-"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint      string `ocisConfig:"events_endpoint" env:"AUDIT_EVENTS_ENDPOINT" desc:"the address of the streaming service"`
	Cluster       string `ocisConfig:"events_cluster" env:"AUDIT_EVENTS_CLUSTER" desc:"the clusterID of the streaming service. Mandatory when using nats"`
	ConsumerGroup string `ocisConfig:"events_group" env:"AUDIT_EVENTS_GROUP" desc:"the customergroup of the service. One group will only get one vopy of an event"`
}

// Auditlog holds audit log information
type Auditlog struct {
	LogToConsole bool   `ocisConfig:"log_to_console" env:"AUDIT_LOG_TO_CONSOLE" desc:"logs to Stdout if true"`
	LogToFile    bool   `ocisConfig:"log_to_file" env:"AUDIT_LOG_TO_FILE" desc:"logs to file if true"`
	FilePath     string `ocisConfig:"filepath" env:"AUDIT_FILEPATH" desc:"filepath to the logfile. Mandatory if LogToFile is true"`
	Format       string `ocisConfig:"format" env:"AUDIT_FORMAT" desc:"log format. using json is advised"`
}
