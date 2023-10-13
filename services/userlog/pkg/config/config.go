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

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP          HTTP                  `yaml:"http"`
	GRPCClientTLS *shared.GRPCClientTLS `yaml:"grpc_client_tls"`

	TokenManager *TokenManager `yaml:"token_manager"`

	RevaGateway     string      `yaml:"reva_gateway" env:"OCIS_REVA_GATEWAY" desc:"CS3 gateway used to look up user metadata"`
	TranslationPath string      `yaml:"translation_path" env:"OCIS_TRANSLATION_PATH;USERLOG_TRANSLATION_PATH" desc:"(optional) Set this to a path with custom translations to overwrite the builtin translations. Note that file and folder naming rules apply, see the documentation for more details."`
	DefaultLanguage string      `yaml:"default_language" env:"OCIS_DEFAULT_LANGUAGE" desc:"(optional) The default language. If not defined, English will be used as default. See the 'settings' service documentation for more details."`
	Events          Events      `yaml:"events"`
	Persistence     Persistence `yaml:"persistence"`

	DisableSSE bool `yaml:"disable_sse" env:"OCIS_DISABLE_SSE,USERLOG_DISABLE_SSE" desc:"Disables server-sent events (sse). When disabled, clients will no longer receive sse notifications."`

	GlobalNotificationsSecret string `yaml:"global_notifications_secret" env:"USERLOG_GLOBAL_NOTIFICATIONS_SECRET" desc:"The secret to secure the global notifications endpoint. Only system admins and users knowing that secret can call the global notifications POST/DELETE endpoints."`

	ServiceAccount ServiceAccount `yaml:"service_account"`

	Context context.Context `yaml:"-"`
}

// Persistence configures the store to use
type Persistence struct {
	Store    string        `yaml:"store" env:"OCIS_PERSISTENT_STORE;USERLOG_STORE" desc:"The type of the store. Supported values are: 'memory', 'ocmem', 'etcd', 'redis', 'redis-sentinel', 'nats-js', 'noop'. See the text description for details."`
	Nodes    []string      `yaml:"nodes" env:"OCIS_PERSISTENT_STORE_NODES;USERLOG_STORE_NODES" desc:"A comma separated list of nodes to access the configured store. This has no effect when 'memory' or 'ocmem' stores are configured. Note that the behaviour how nodes are used is dependent on the library of the configured store."`
	Database string        `yaml:"database" env:"USERLOG_STORE_DATABASE" desc:"The database name the configured store should use."`
	Table    string        `yaml:"table" env:"USERLOG_STORE_TABLE" desc:"The database table the store should use."`
	TTL      time.Duration `yaml:"ttl" env:"OCIS_PERSISTENT_STORE_TTL;USERLOG_STORE_TTL" desc:"Time to live for events in the store. The duration can be set as number followed by a unit identifier like s, m or h. Defaults to '336h' (2 weeks)."`
	Size     int           `yaml:"size" env:"OCIS_PERSISTENT_STORE_SIZE;USERLOG_STORE_SIZE" desc:"The maximum quantity of items in the store. Only applies when store type 'ocmem' is configured. Defaults to 512."`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"OCIS_EVENTS_ENDPOINT;USERLOG_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster              string `yaml:"cluster" env:"OCIS_EVENTS_CLUSTER;USERLOG_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;USERLOG_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"OCIS_EVENTS_TLS_ROOT_CA_CERTIFICATE;USERLOG_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided NOTIFICATIONS_EVENTS_TLS_INSECURE will be seen as false."`
	EnableTLS            bool   `yaml:"enable_tls" env:"OCIS_EVENTS_ENABLE_TLS;USERLOG_EVENTS_ENABLE_TLS" desc:"Enable TLS for the connection to the events broker. The events broker is the ocis service which receives and delivers events between the services.."`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `yaml:"allow_origins" env:"OCIS_CORS_ALLOW_ORIGINS;USERLOG_CORS_ALLOW_ORIGINS" desc:"A comma-separated list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin"`
	AllowedMethods   []string `yaml:"allow_methods" env:"OCIS_CORS_ALLOW_METHODS;USERLOG_CORS_ALLOW_METHODS" desc:"A comma-separated list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method"`
	AllowedHeaders   []string `yaml:"allow_headers" env:"OCIS_CORS_ALLOW_HEADERS;USERLOG_CORS_ALLOW_HEADERS" desc:"A blank or comma-separated list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers."`
	AllowCredentials bool     `yaml:"allow_credentials" env:"OCIS_CORS_ALLOW_CREDENTIALS;USERLOG_CORS_ALLOW_CREDENTIALS" desc:"Allow credentials for CORS.See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials."`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string                `yaml:"addr" env:"USERLOG_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	Namespace string                `yaml:"-"`
	Root      string                `yaml:"root" env:"USERLOG_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	CORS      CORS                  `yaml:"cors"`
	TLS       shared.HTTPServiceTLS `yaml:"tls"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `yaml:"jwt_secret" env:"OCIS_JWT_SECRET;USERLOG_JWT_SECRET" desc:"The secret to mint and validate jwt tokens."`
}

// ServiceAccount is the configuration for the used service account
type ServiceAccount struct {
	ServiceAccountID     string `yaml:"service_account_id" env:"OCIS_SERVICE_ACCOUNT_ID;USERLOG_SERVICE_ACCOUNT_ID" desc:"The ID of the service account the service should use. See the 'auth-service' service description for more details."`
	ServiceAccountSecret string `yaml:"service_account_secret" env:"OCIS_SERVICE_ACCOUNT_SECRET;USERLOG_SERVICE_ACCOUNT_SECRET" desc:"The service account secret."`
}
