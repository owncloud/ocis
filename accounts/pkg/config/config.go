// Package config should be moved to internal
package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `ocisConfig:"allowed_origins"`
	AllowedMethods   []string `ocisConfig:"allowed_methods"`
	AllowedHeaders   []string `ocisConfig:"allowed_headers"`
	AllowCredentials bool     `ocisConfig:"allowed_credentials"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"ACCOUNTS_HTTP_ADDR"`
	Namespace string
	Root      string `ocisConfig:"root" env:"ACCOUNTS_HTTP_ROOT"`
	CacheTTL  int    `ocisConfig:"cache_ttl" env:"ACCOUNTS_CACHE_TTL"`
	CORS      CORS   `ocisConfig:"cors"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `ocisConfig:"addr" env:"ACCOUNTS_GRPC_ADDR"`
	Namespace string
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path" env:"ACCOUNTS_ASSET_PATH"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET;ACCOUNTS_JWT_SECRET"`
}

// Repo defines which storage implementation is to be used.
type Repo struct {
	Backend string `ocisConfig:"backend"  env:"ACCOUNTS_STORAGE_BACKEND"`
	Disk    Disk   `ocisConfig:"disk"`
	CS3     CS3    `ocisConfig:"cs3"`
}

// Disk is the local disk implementation of the storage.
type Disk struct {
	Path string `ocisConfig:"path" env:"ACCOUNTS_STORAGE_DISK_PATH"`
}

// CS3 is the cs3 implementation of the storage.
type CS3 struct {
	ProviderAddr string `ocisConfig:"provider_addr" env:"ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR"`
	JWTSecret    string `ocisConfig:"jwt_secret" env:"ACCOUNTS_STORAGE_CS3_JWT_SECRET"`
}

// ServiceUser defines the user required for EOS.
type ServiceUser struct {
	UUID     string `ocisConfig:"uuid" env:"ACCOUNTS_SERVICE_USER_UUID"`
	Username string `ocisConfig:"username" env:"ACCOUNTS_SERVICE_USER_USERNAME"`
	UID      int64  `ocisConfig:"uid" env:"ACCOUNTS_SERVICE_USER_UID"`
	GID      int64  `ocisConfig:"gid" env:"ACCOUNTS_SERVICE_USER_GID"`
}

// Index defines config for indexes.
type Index struct {
	UID UIDBound `ocisConfig:"uid"`
	GID GIDBound `ocisConfig:"gid"`
}

// GIDBound defines a lower and upper bound.
type GIDBound struct {
	Lower int64 `ocisConfig:"lower" env:"ACCOUNTS_GID_INDEX_LOWER_BOUND"`
	Upper int64 `ocisConfig:"upper" env:"ACCOUNTS_GID_INDEX_UPPER_BOUND"`
}

// UIDBound defines a lower and upper bound.
type UIDBound struct {
	Lower int64 `ocisConfig:"lower" env:"ACCOUNTS_UID_INDEX_LOWER_BOUND"`
	Upper int64 `ocisConfig:"upper" env:"ACCOUNTS_UID_INDEX_UPPER_BOUND"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;ACCOUNTS_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;ACCOUNTS_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;ACCOUNTS_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;ACCOUNTS_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"ACCOUNTS_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;ACCOUNTS_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;ACCOUNTS_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;ACCOUNTS_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;ACCOUNTS_LOG_FILE"`
}

// Config merges all Account config parameters.
type Config struct {
	//*shared.Commons

	HTTP               HTTP         `ocisConfig:"http"`
	GRPC               GRPC         `ocisConfig:"grpc"`
	Service            Service      `ocisConfig:"service"`
	Asset              Asset        `ocisConfig:"asset"`
	Log                Log          `ocisConfig:"log"`
	TokenManager       TokenManager `ocisConfig:"token_manager"`
	Repo               Repo         `ocisConfig:"repo"`
	Index              Index        `ocisConfig:"index"`
	ServiceUser        ServiceUser  `ocisConfig:"service_user"`
	HashDifficulty     int          `ocisConfig:"hash_difficulty" env:"ACCOUNTS_HASH_DIFFICULTY"`
	DemoUsersAndGroups bool         `ocisConfig:"demo_users_and_groups" env:"ACCOUNTS_DEMO_USERS_AND_GROUPS"`
	Tracing            Tracing      `ocisConfig:"tracing"`

	Context    context.Context
	Supervised bool
}

// New returns a new config.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{

		HTTP: HTTP{
			Addr:      "127.0.0.1:9181",
			Namespace: "com.owncloud.web",
			Root:      "/",
			CacheTTL:  604800, // 7 days
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		GRPC: GRPC{
			Addr:      "127.0.0.1:9180",
			Namespace: "com.owncloud.api",
		},
		Service: Service{
			Name: "accounts",
		},
		Asset: Asset{},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		HashDifficulty:     11,
		DemoUsersAndGroups: true,
		Repo: Repo{
			Backend: "CS3",
			Disk: Disk{
				Path: path.Join(defaults.BaseDataPath(), "accounts"),
			},
			CS3: CS3{
				ProviderAddr: "localhost:9215",
				JWTSecret:    "Pive-Fumkiu4",
			},
		},
		Index: Index{
			UID: UIDBound{
				Lower: 0,
				Upper: 1000,
			},
			GID: GIDBound{
				Lower: 0,
				Upper: 1000,
			},
		},
		ServiceUser: ServiceUser{
			UUID:     "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
			Username: "",
			UID:      0,
			GID:      0,
		},
		Tracing: Tracing{
			Type:    "jaeger",
			Service: "accounts",
		},
	}
}
