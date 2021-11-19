package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
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
	Addr string `ocisConfig:"addr"`
	Root string `ocisConfig:"root"`
	CORS CORS   `ocisConfig:"cors"`
}

// Service defines the available service configuration.
type Service struct {
	Name      string `ocisConfig:"name"`
	Namespace string `ocisConfig:"namespace"`
	Version   string `ocisConfig:"version"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `ocisConfig:"address"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret"`
}

// IdentityManagement keeps track of the OIDC address. This is because Reva requisite of uniqueness for users
// is based in the combination of IDP hostname + UserID. For more information see:
// https://github.com/cs3org/reva/blob/4fd0229f13fae5bc9684556a82dbbd0eced65ef9/pkg/storage/utils/decomposedfs/node/node.go#L856-L865
type IdentityManagement struct {
	Address string `ocisConfig:"address"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File               string             `ocisConfig:"file"`
	Log                *shared.Log        `ocisConfig:"log"`
	Debug              Debug              `ocisConfig:"debug"`
	HTTP               HTTP               `ocisConfig:"http"`
	Tracing            Tracing            `ocisConfig:"tracing"`
	TokenManager       TokenManager       `ocisConfig:"token_manager"`
	Service            Service            `ocisConfig:"service"`
	AccountBackend     string             `ocisConfig:"account_backend"`
	Reva               Reva               `ocisConfig:"reva"`
	StorageUsersDriver string             `ocisConfig:"storage_users_driver"`
	MachineAuthAPIKey  string             `ocisConfig:"machine_auth_api_key"`
	IdentityManagement IdentityManagement `ocisConfig:"identity_management"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
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
			Addr: "127.0.0.1:9110",
			Root: "/ocs",
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
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
		Service: Service{
			Name:      "ocs",
			Namespace: "com.owncloud.web",
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
