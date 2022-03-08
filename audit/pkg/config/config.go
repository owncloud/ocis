package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Log   *Log  `ocisConfig:"log"`
	Debug Debug `ocisConfig:"debug"`

	Events   Events   `ocisConfig:"events"`
	Auditlog Auditlog `ocisConfig:"auditlog"`

	Context context.Context
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint      string `ocisConfig:"events_endpoint" env:"AUDIT_EVENTS_ENDPOINT"`
	Cluster       string `ocisConfig:"events_cluster" env:"AUDIT_EVENTS_CLUSTER"`
	ConsumerGroup string `ocisConfig:"events_group" env:"AUDIT_EVENTS_GROUP"`
}

// Auditlog holds audit log information
type Auditlog struct {
	LogToConsole bool   `ocisConfig:"log_to_console" env:"AUDIT_LOG_TO_CONSOLE"`
	LogToFile    bool   `ocisConfig:"log_to_file" env:"AUDIT_LOG_TO_FILE"`
	FilePath     string `ocisConfig:"filepath" env:"AUDIT_FILEPATH"`
	Format       string `ocisConfig:"format" env:"AUDIT_FORMAT"`
}
