// Package config should be moved to internal
package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// LDAP defines the available ldap configuration.
type LDAP struct {
	Hostname     string     `ocisConfig:"hostname"`
	Port         int        `ocisConfig:"port"`
	BaseDN       string     `ocisConfig:"base_dn"`
	UserFilter   string     `ocisConfig:"user_filter"`
	GroupFilter  string     `ocisConfig:"group_filter"`
	BindDN       string     `ocisConfig:"bind_dn"`
	BindPassword string     `ocisConfig:"bind_password"`
	IDP          string     `ocisConfig:"idp"`
	Schema       LDAPSchema `ocisConfig:"schema"`
}

// LDAPSchema defines the available ldap schema configuration.
type LDAPSchema struct {
	AccountID   string `ocisConfig:"account_id"`
	Identities  string `ocisConfig:"identities"`
	Username    string `ocisConfig:"username"`
	DisplayName string `ocisConfig:"display_name"`
	Mail        string `ocisConfig:"mail"`
	Groups      string `ocisConfig:"groups"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `ocisConfig:"allowed_origins"`
	AllowedMethods   []string `ocisConfig:"allowed_methods"`
	AllowedHeaders   []string `ocisConfig:"allowed_headers"`
	AllowCredentials bool     `ocisConfig:"allowed_credentials"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr"`
	Namespace string `ocisConfig:"namespace"`
	Root      string `ocisConfig:"root"`
	CacheTTL  int    `ocisConfig:"cache_ttl"`
	CORS      CORS   `ocisConfig:"cors"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `ocisConfig:"addr"`
	Namespace string `ocisConfig:"namespace"`
}

// Server configures a server.
type Server struct {
	Version            string `ocisConfig:"version"`
	Name               string `ocisConfig:"name"`
	HashDifficulty     int    `ocisConfig:"hash_difficulty"`
	DemoUsersAndGroups bool   `ocisConfig:"demo_users_and_groups"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `ocisConfig:"path"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret"`
}

// Repo defines which storage implementation is to be used.
type Repo struct {
	Backend string `ocisConfig:"backend"`
	Disk    Disk   `ocisConfig:"disk"`
	CS3     CS3    `ocisConfig:"cs3"`
}

// Disk is the local disk implementation of the storage.
type Disk struct {
	Path string `ocisConfig:"path"`
}

// CS3 is the cs3 implementation of the storage.
type CS3 struct {
	ProviderAddr string `ocisConfig:"provider_addr"`
	JWTSecret    string `ocisConfig:"jwt_secret"`
}

// ServiceUser defines the user required for EOS.
type ServiceUser struct {
	UUID     string `ocisConfig:"uuid"`
	Username string `ocisConfig:"username"`
	UID      int64  `ocisConfig:"uid"`
	GID      int64  `ocisConfig:"gid"`
}

// Index defines config for indexes.
type Index struct {
	UID Bound `ocisConfig:"uid"`
	GID Bound `ocisConfig:"gid"`
}

// Bound defines a lower and upper bound.
type Bound struct {
	Lower int64 `ocisConfig:"lower"`
	Upper int64 `ocisConfig:"upper"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// Config merges all Account config parameters.
type Config struct {
	*shared.Commons

	LDAP         LDAP         `ocisConfig:"ldap"`
	HTTP         HTTP         `ocisConfig:"http"`
	GRPC         GRPC         `ocisConfig:"grpc"`
	Server       Server       `ocisConfig:"server"`
	Asset        Asset        `ocisConfig:"asset"`
	Log          *shared.Log  `ocisConfig:"log"`
	TokenManager TokenManager `ocisConfig:"token_manager"`
	Repo         Repo         `ocisConfig:"repo"`
	Index        Index        `ocisConfig:"index"`
	ServiceUser  ServiceUser  `ocisConfig:"service_user"`
	Tracing      Tracing      `ocisConfig:"tracing"`

	Context    context.Context
	Supervised bool
}

// New returns a new config.
func New() *Config {
	return &Config{
		Log: &shared.Log{},
	}
}

func DefaultConfig() *Config {
	return &Config{
		LDAP: LDAP{},
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
		Server: Server{
			Name:               "accounts",
			HashDifficulty:     11,
			DemoUsersAndGroups: true,
		},
		Asset: Asset{},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
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
			UID: Bound{
				Lower: 0,
				Upper: 1000,
			},
			GID: Bound{
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

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv(cfg *Config) []string {
	var r = make([]string, len(structMappings(cfg)))
	for i := range structMappings(cfg) {
		r = append(r, structMappings(cfg)[i].EnvVars...)
	}

	return r
}
