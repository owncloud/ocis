// Package config should be moved to internal
package config

import (
	"context"
	"fmt"
	"path"
	"reflect"

	gofig "github.com/gookit/config/v2"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// LDAP defines the available ldap configuration.
type LDAP struct {
	Hostname     string     `mapstructure:"hostname"`
	Port         int        `mapstructure:"port"`
	BaseDN       string     `mapstructure:"base_dn"`
	UserFilter   string     `mapstructure:"user_filter"`
	GroupFilter  string     `mapstructure:"group_filter"`
	BindDN       string     `mapstructure:"bind_dn"`
	BindPassword string     `mapstructure:"bind_password"`
	IDP          string     `mapstructure:"idp"`
	Schema       LDAPSchema `mapstructure:"schema"`
}

// LDAPSchema defines the available ldap schema configuration.
type LDAPSchema struct {
	AccountID   string `mapstructure:"account_id"`
	Identities  string `mapstructure:"identities"`
	Username    string `mapstructure:"username"`
	DisplayName string `mapstructure:"display_name"`
	Mail        string `mapstructure:"mail"`
	Groups      string `mapstructure:"groups"`
}

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	AllowCredentials bool     `mapstructure:"allowed_credentials"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `mapstructure:"addr"`
	Namespace string `mapstructure:"namespace"`
	Root      string `mapstructure:"root"`
	CacheTTL  int    `mapstructure:"cache_ttl"`
	CORS      CORS   `mapstructure:"cors"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `mapstructure:"addr"`
	Namespace string `mapstructure:"namespace"`
}

// Server configures a server.
type Server struct {
	Version            string `mapstructure:"version"`
	Name               string `mapstructure:"name"`
	HashDifficulty     int    `mapstructure:"hash_difficulty"`
	DemoUsersAndGroups bool   `mapstructure:"demo_users_and_groups"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `mapstructure:"path"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `mapstructure:"jwt_secret"`
}

// Log defines the available logging configuration.
type Log struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Color  bool   `mapstructure:"color"`
	File   string `mapstructure:"file"`
}

// Repo defines which storage implementation is to be used.
type Repo struct {
	Backend string `mapstructure:"backend"`
	Disk    Disk   `mapstructure:"disk"`
	CS3     CS3    `mapstructure:"cs3"`
}

// Disk is the local disk implementation of the storage.
type Disk struct {
	Path string `mapstructure:"path"`
}

// CS3 is the cs3 implementation of the storage.
type CS3 struct {
	ProviderAddr string `mapstructure:"provider_addr"`
	DataURL      string `mapstructure:"data_url"`
	DataPrefix   string `mapstructure:"data_prefix"`
	JWTSecret    string `mapstructure:"jwt_secret"`
}

// ServiceUser defines the user required for EOS.
type ServiceUser struct {
	UUID     string `mapstructure:"uuid"`
	Username string `mapstructure:"username"`
	UID      int64  `mapstructure:"uid"`
	GID      int64  `mapstructure:"gid"`
}

// Index defines config for indexes.
type Index struct {
	UID Bound `mapstructure:"uid"`
	GID Bound `mapstructure:"gid"`
}

// Bound defines a lower and upper bound.
type Bound struct {
	Lower int64 `mapstructure:"lower"`
	Upper int64 `mapstructure:"upper"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"enabled"`
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	Collector string `mapstructure:"collector"`
	Service   string `mapstructure:"service"`
}

// Config merges all Account config parameters.
type Config struct {
	LDAP         LDAP         `mapstructure:"ldap"`
	HTTP         HTTP         `mapstructure:"http"`
	GRPC         GRPC         `mapstructure:"grpc"`
	Server       Server       `mapstructure:"server"`
	Asset        Asset        `mapstructure:"asset"`
	Log          Log          `mapstructure:"log"`
	TokenManager TokenManager `mapstructure:"token_manager"`
	Repo         Repo         `mapstructure:"repo"`
	Index        Index        `mapstructure:"index"`
	ServiceUser  ServiceUser  `mapstructure:"service_user"`
	Tracing      Tracing      `mapstructure:"tracing"`

	Context    context.Context
	Supervised bool
}

// New returns a new config.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		LDAP: LDAP{},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9181",
			Namespace: "com.owncloud.web",
			Root:      "/",
			CacheTTL:  604800,
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
		Log:   Log{},
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
				DataURL:      "http://localhost:9216",
				DataPrefix:   "data",
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
func GetEnv() []string {
	var r = make([]string, len(structMappings(&Config{})))
	for i := range structMappings(&Config{}) {
		r = append(r, structMappings(&Config{})[i].EnvVars...)
	}

	return r
}

// UnmapEnv loads values from the gooconf.Config argument and sets them in the expected destination.
func (c *Config) UnmapEnv(gooconf *gofig.Config) error {
	vals := structMappings(c)
	for i := range vals {
		for j := range vals[i].EnvVars {
			// we need to guard against v != "" because this is the condition that checks that the value is set from the environment.
			// the `ok` guard is not enough, apparently.
			if v, ok := gooconf.GetValue(vals[i].EnvVars[j]); ok && v != "" {

				// get the destination type from destination
				switch reflect.ValueOf(vals[i].Destination).Type().String() {
				case "*bool":
					r := gooconf.Bool(vals[i].EnvVars[j])
					*vals[i].Destination.(*bool) = r
				case "*string":
					r := gooconf.String(vals[i].EnvVars[j])
					*vals[i].Destination.(*string) = r
				case "*int":
					r := gooconf.Int(vals[i].EnvVars[j])
					*vals[i].Destination.(*int) = r
				case "*float64":
					// defaults to float64
					r := gooconf.Float(vals[i].EnvVars[j])
					*vals[i].Destination.(*float64) = r
				default:
					// it is unlikely we will ever get here. Let this serve more as a runtime check for when debugging.
					return fmt.Errorf("invalid type for env var: `%v`", vals[i].EnvVars[j])
				}
			}
		}
	}

	return nil
}

type mapping struct {
	EnvVars     []string    // name of the EnvVars var.
	Destination interface{} // memory address of the original config value to modify.
}

// structMappings binds a set of environment variables to a destination on cfg.
func structMappings(cfg *Config) []mapping {
	return []mapping{
		{
			EnvVars:     []string{"ACCOUNTS_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		{
			EnvVars:     []string{"ACCOUNTS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		{
			EnvVars:     []string{"ACCOUNTS_HTTP_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		{
			EnvVars:     []string{"ACCOUNTS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		{
			EnvVars:     []string{"ACCOUNTS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		{
			EnvVars:     []string{"ACCOUNTS_CACHE_TTL"},
			Destination: &cfg.HTTP.CacheTTL,
		},
		{
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		{
			EnvVars:     []string{"ACCOUNTS_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		{
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		{
			EnvVars:     []string{"ACCOUNTS_HASH_DIFFICULTY"},
			Destination: &cfg.Server.HashDifficulty,
		},
		{
			EnvVars:     []string{"ACCOUNTS_DEMO_USERS_AND_GROUPS"},
			Destination: &cfg.Server.DemoUsersAndGroups,
		},
		{
			EnvVars:     []string{"ACCOUNTS_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		{
			EnvVars:     []string{"ACCOUNTS_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_BACKEND"},
			Destination: &cfg.Repo.Backend,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_DISK_PATH"},
			Destination: &cfg.Repo.Disk.Path,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_PROVIDER_ADDR"},
			Destination: &cfg.Repo.CS3.ProviderAddr,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_DATA_URL"},
			Destination: &cfg.Repo.CS3.DataURL,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_DATA_PREFIX"},
			Destination: &cfg.Repo.CS3.DataPrefix,
		},
		{
			EnvVars:     []string{"ACCOUNTS_STORAGE_CS3_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.Repo.CS3.JWTSecret,
		},
		{
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_UUID"},
			Destination: &cfg.ServiceUser.UUID,
		},
		{
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_USERNAME"},
			Destination: &cfg.ServiceUser.Username,
		},
		{
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_UID"},
			Destination: &cfg.ServiceUser.UID,
		},
		{
			EnvVars:     []string{"ACCOUNTS_SERVICE_USER_GID"},
			Destination: &cfg.ServiceUser.GID,
		},
		{
			EnvVars:     []string{"ACCOUNTS_UID_INDEX_LOWER_BOUND"},
			Destination: &cfg.Index.UID.Lower,
		},
		{
			EnvVars:     []string{"ACCOUNTS_GID_INDEX_LOWER_BOUND"},
			Destination: &cfg.Index.GID.Lower,
		},
		{
			EnvVars:     []string{"ACCOUNTS_UID_INDEX_UPPER_BOUND"},
			Destination: &cfg.Index.UID.Upper,
		},
		{
			EnvVars:     []string{"ACCOUNTS_GID_INDEX_UPPER_BOUND"},
			Destination: &cfg.Index.GID.Upper,
		},
	}
}

// TODO(refs) What is with the variables with no destination defined?
//&cli.StringSliceFlag{
//Name:    "cors-allowed-origins",
//Value:   cli.NewStringSlice("*"),
//Usage:   "Set the allowed CORS origins",
//EnvVars: []string{"ACCOUNTS_CORS_ALLOW_ORIGINS", "OCIS_CORS_ALLOW_ORIGINS"},
//},
//&cli.StringSliceFlag{
//Name:    "cors-allowed-methods",
//Value:   cli.NewStringSlice("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"),
//Usage:   "Set the allowed CORS origins",
//EnvVars: []string{"ACCOUNTS_CORS_ALLOW_METHODS", "OCIS_CORS_ALLOW_METHODS"},
//},
//&cli.StringSliceFlag{
//Name:    "cors-allowed-headers",
//Value:   cli.NewStringSlice("Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"),
//Usage:   "Set the allowed CORS origins",
//EnvVars: []string{"ACCOUNTS_CORS_ALLOW_HEADERS", "OCIS_CORS_ALLOW_HEADERS"},
//},
//&cli.BoolFlag{
//Name:    "cors-allow-credentials",
//Value:   flags.OverrideDefaultBool(cfg.HTTP.CORS.AllowCredentials, true),
//Usage:   "Allow credentials for CORS",
//EnvVars: []string{"ACCOUNTS_CORS_ALLOW_CREDENTIALS", "OCIS_CORS_ALLOW_CREDENTIALS"},
//},
