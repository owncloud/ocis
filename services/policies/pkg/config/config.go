package config

import (
	"context"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"time"
)

// Config combines all available configuration parts.
type Config struct {
	Commons           *shared.Commons       `yaml:"-"` // don't use this directly as configuration for a service
	GRPC              GRPC                  `yaml:"grpc"`
	Service           Service               `yaml:"-"`
	TokenManager      *TokenManager         `yaml:"token_manager"`
	Events            Events                `yaml:"events"`
	Reva              *shared.Reva          `yaml:"reva"`
	GRPCClientTLS     *shared.GRPCClientTLS `yaml:"grpc_client_tls"`
	MachineAuthAPIKey string                `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;POLICIES_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services."`
	Context           context.Context       `yaml:"-"`
	Log               *Log                  `yaml:"log"`
	Engine            Engine                `yaml:"engines"`
	Postprocessing    Postprocessing        `yaml:"postprocessing"`
}

// Service defines the available service configuration.
type Service struct {
	Name string `yaml:"-"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string                 `ocisConfig:"addr" env:"POLICIES_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string                 `ocisConfig:"-" yaml:"-"`
	TLS       *shared.GRPCServiceTLS `yaml:"tls"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;POLICIES_JWT_SECRET" desc:"The secret to mint and validate jwt tokens."`
}

// Engine configures the policy engine.
type Engine struct {
	Timeout  time.Duration `yaml:"timeout" env:"POLICIES_ENGINE_TIMEOUT" desc:"Sets the timeout."`
	Policies []string      `yaml:"policies"`
}

// Postprocessing defines the config options for the postprocessing policy handling.
type Postprocessing struct {
	Query string `yaml:"query" env:"POLICIES_POSTPROCESSING_QUERY" desc:"Sets the postprocessing query."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"POLICIES_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster              string `yaml:"cluster" env:"POLICIES_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;POLICIES_EVENTS_TLS_INSECURE" desc:"Whether the server should skip the client certificate verification during the TLS handshake."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"POLICIES_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided POLICIES_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;POLICIES_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services."`
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;POLICIES_LOG_LEVEL" desc:"The log level. Valid values are: \"panic\", \"fatal\", \"error\", \"warn\", \"info\", \"debug\", \"trace\"."`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;POLICIES_LOG_PRETTY" desc:"Activates pretty log output."`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;POLICIES_LOG_COLOR" desc:"Activates colorized log output."`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;POLICIES_LOG_FILE" desc:"The path to the log file. Activates logging to this file if set."`
}
