package config

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons           *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service
	HTTP              HTTP            `yaml:"http"`
	GRPC              GRPC            `yaml:"grpc"`
	Service           Service         `yaml:"-"`
	TokenManager      *TokenManager   `yaml:"token_manager"`
	Events            Events          `yaml:"events"`
	Reva              *shared.Reva    `yaml:"reva"`
	MachineAuthAPIKey string          `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;AUTHZ_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services."`
	Context           context.Context `yaml:"-"`
	OPA               OPA             `yaml:"opa"`
}

// Service defines the available service configuration.
type Service struct {
	Name string `yaml:"-"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `yaml:"addr" env:"AUTHZ_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string `yaml:"-"`
	Root      string `yaml:"root" env:"AUTHZ_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `yaml:"addr" env:"AUTHZ_GRPC_ADDR" desc:"The bind address of the GRPC service."`
	Namespace string `yaml:"-"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;AUTHZ_JWT_SECRET" desc:"The secret to mint and validate jwt tokens."`
}

// OPA configures the opa policy middleware.
type OPA struct {
	Enabled  bool          `yaml:"enabled"`
	Policies []string      `yaml:"policies"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"AUTHZ_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster              string `yaml:"cluster" env:"AUTHZ_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;AUTHZ_EVENTS_TLS_INSECURE" desc:"Whether the server should skip the client certificate verification during the TLS handshake."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"AUTHZ_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided AUTHZ_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;AUTHZ_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services."`
}
