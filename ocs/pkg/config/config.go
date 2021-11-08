package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `mapstructure:"addr"`
	Token  string `mapstructure:"token"`
	Pprof  bool   `mapstructure:"pprof"`
	Zpages bool   `mapstructure:"zpages"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr string `mapstructure:"addr"`
	Root string `mapstructure:"root"`
	CORS CORS   `mapstructure:"cors"`
}

// Service defines the available service configuration.
type Service struct {
	Name      string `mapstructure:"name"`
	Namespace string `mapstructure:"namespace"`
	Version   string `mapstructure:"version"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"enabled"`
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	Collector string `mapstructure:"collector"`
	Service   string `mapstructure:"service"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

// IdentityManagement keeps track of the OIDC address. This is because Reva requisite of uniqueness for users
// is based in the combination of IDP hostname + UserID. For more information see:
// https://github.com/cs3org/reva/blob/4fd0229f13fae5bc9684556a82dbbd0eced65ef9/pkg/storage/utils/decomposedfs/node/node.go#L856-L865
type IdentityManagement struct {
	Address string `mapstructure:"address"`
}

// Config combines all available configuration parts.
type Config struct {
	File               string             `mapstructure:"file"`
	Log                shared.Log         `mapstructure:"log"`
	Debug              Debug              `mapstructure:"debug"`
	HTTP               HTTP               `mapstructure:"http"`
	Tracing            Tracing            `mapstructure:"tracing"`
	TokenManager       TokenManager       `mapstructure:"token_manager"`
	Service            Service            `mapstructure:"service"`
	AccountBackend     string             `mapstructure:"account_backend"`
	RevaAddress        string             `mapstructure:"reva_address"`
	StorageUsersDriver string             `mapstructure:"storage_users_driver"`
	MachineAuthAPIKey  string             `mapstructure:"machine_auth_api_key"`
	IdentityManagement IdentityManagement `mapstructure:"identity_management"`

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
		Log: shared.Log{},
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
		RevaAddress:        "127.0.0.1:9142",
		StorageUsersDriver: "ocis",
		MachineAuthAPIKey:  "change-me-please",
		IdentityManagement: IdentityManagement{
			Address: "https://localhost:9200",
		},
	}
}

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv() []string {
	var r = make([]string, len(structMappings(&Config{})))
	for i := range structMappings(&Config{}) {
		r = append(r, structMappings(&Config{})[i].EnvVars...)
	}

	return r
}
