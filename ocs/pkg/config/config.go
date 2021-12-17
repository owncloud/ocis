package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"OCS_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"OCS_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"OCS_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"OCS_DEBUG_ZPAGES"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `ocisConfig:"allowed_origins"`
	AllowedMethods   []string `ocisConfig:"allowed_methods"`
	AllowedHeaders   []string `ocisConfig:"allowed_headers"`
	AllowCredentials bool     `ocisConfig:"allow_credentials"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"OCS_HTTP_ADDR"`
	Root      string `ocisConfig:"root" env:"OCS_HTTP_ROOT"`
	Namespace string
	CORS      CORS `ocisConfig:"cors"`
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;OCS_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;OCS_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;OCS_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;OCS_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"OCS_TRACING_SERVICE"`
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;OCS_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;OCS_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;OCS_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;OCS_LOG_FILE"`
}

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `ocisConfig:"address" env:"REVA_GATEWAY"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET;OCS_JWT_SECRET"`
}

// IdentityManagement keeps track of the OIDC address. This is because Reva requisite of uniqueness for users
// is based in the combination of IDP hostname + UserID. For more information see:
// https://github.com/cs3org/reva/blob/4fd0229f13fae5bc9684556a82dbbd0eced65ef9/pkg/storage/utils/decomposedfs/node/node.go#L856-L865
type IdentityManagement struct {
	Address string `ocisConfig:"address" env:"OCIS_URL;OCS_IDM_ADDRESS"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service `ocisConfig:"service"`

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	HTTP HTTP `ocisConfig:"http"`

	TokenManager TokenManager `ocisConfig:"token_manager"`
	Reva         Reva         `ocisConfig:"reva"`

	IdentityManagement IdentityManagement `ocisConfig:"identity_management"`

	AccountBackend     string `ocisConfig:"account_backend" env:"OCS_ACCOUNT_BACKEND_TYPE"`
	StorageUsersDriver string `ocisConfig:"storage_users_driver" env:"STORAGE_USERS_DRIVER;OCS_STORAGE_USERS_DRIVER"`
	MachineAuthAPIKey  string `ocisConfig:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;OCS_MACHINE_AUTH_API_KEY"`

	Context    context.Context
	Supervised bool
}

// DefaultConfig provides default values for a config struct.
func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9114",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9110",
			Root:      "/ocs",
			Namespace: "com.owncloud.web",
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		Service: Service{
			Name: "ocs",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "ocs",
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		AccountBackend:     "accounts",
		Reva:               Reva{Address: "127.0.0.1:9142"},
		StorageUsersDriver: "ocis",
		MachineAuthAPIKey:  "change-me-please",
		IdentityManagement: IdentityManagement{
			Address: "https://localhost:9200",
		},
	}
}
