package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"GRAPH_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"GRAPH_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"GRAPH_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"GRAPH_DEBUG_ZPAGES"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr" env:"GRAPH_HTTP_ADDR"`
	Namespace string
	Root      string `ocisConfig:"root" env:"GRAPH_HTTP_ROOT"`
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;GRAPH_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;GRAPH_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;GRAPH_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;GRAPH_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"GRAPH_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;GRAPH_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;GRAPH_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;GRAPH_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;GRAPH_LOG_FILE"`
}

// Reva defines all available REVA configuration.
type Reva struct {
	Address string `ocisConfig:"address" env:"REVA_GATEWAY"`
}

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string `ocisConfig:"jwt_secret" env:"OCIS_JWT_SECRET;GRAPH_JWT_SECRET"`
}

type Spaces struct {
	WebDavBase   string `ocisConfig:"webdav_base" env:"GRAPH_SPACES_WEBDAV_BASE"`
	WebDavPath   string `ocisConfig:"webdav_path" env:"GRAPH_SPACES_WEBDAV_PATH"`
	DefaultQuota string `ocisConfig:"default_quota" env:"GRAPH_SPACES_DEFAULT_QUOTA"`
}

// TODO: do we really need a ldap backend if CS3 also does LDAP!?
type LDAP struct {
	URI          string `ocisConfig:"uri" env:"GRAPH_LDAP_URI"`
	BindDN       string `ocisConfig:"bind_dn" env:"GRAPH_LDAP_BIND_DN"`
	BindPassword string `ocisConfig:"bind_password" env:"GRAPH_LDAP_BIND_PASSWORD"`

	UserBaseDN               string `ocisConfig:"user_base_dn" env:"GRAPH_LDAP_USER_BASE_DN"`
	UserSearchScope          string `ocisConfig:"user_search_scope" env:"GRAPH_LDAP_USER_SCOPE"`
	UserFilter               string `ocisConfig:"user_filter" env:"GRAPH_LDAP_USER_FILTER"`
	UserEmailAttribute       string `ocisConfig:"user_mail_attribute" env:"GRAPH_LDAP_USER_EMAIL_ATTRIBUTE"`
	UserDisplayNameAttribute string `ocisConfig:"user_displayname_attribute" env:"GRAPH_LDAP_USER_DISPLAYNAME_ATTRIBUTE"`
	UserNameAttribute        string `ocisConfig:"user_name_attribute" env:"GRAPH_LDAP_USER_NAME_ATTRIBUTE"`
	UserIDAttribute          string `ocisConfig:"user_id_attribute" env:"GRAPH_LDAP_USER_UID_ATTRIBUTE"`

	GroupBaseDN        string `ocisConfig:"group_base_dn" env:"GRAPH_LDAP_GROUP_BASE_DN"`
	GroupSearchScope   string `ocisConfig:"group_search_scope" env:"GRAPH_LDAP_GROUP_SEARCH_SCOPE"`
	GroupFilter        string `ocisConfig:"group_filter" env:"GRAPH_LDAP_GROUP_FILTER"`
	GroupNameAttribute string `ocisConfig:"group_name_attribute" env:"GRAPH_LDAP_GROUP_NAME_ATTRIBUTE"`
	GroupIDAttribute   string `ocisConfig:"group_id_attribute" env:"GRAPH_LDAP_GROUP_ID_ATTRIBUTE"`
}

type Identity struct {
	Backend string `ocisConfig:"backend" env:"GRAPH_IDENTITY_BACKEND"`
	LDAP    LDAP   `ocisConfig:"ldap"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Log          Log          `ocisConfig:"log"`
	Debug        Debug        `ocisConfig:"debug"`
	HTTP         HTTP         `ocisConfig:"http"`
	Service      Service      `ocisConfig:"service"`
	Tracing      Tracing      `ocisConfig:"tracing"`
	Reva         Reva         `ocisConfig:"reva"`
	TokenManager TokenManager `ocisConfig:"token_manager"`
	Spaces       Spaces       `ocisConfig:"spaces"`
	Identity     Identity     `ocisConfig:"identity"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:  "127.0.0.1:9124",
			Token: "",
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9120",
			Namespace: "com.owncloud.graph",
			Root:      "/graph",
		},
		Service: Service{
			Name: "graph",
		},
		Tracing: Tracing{
			Enabled: false,
			Type:    "jaeger",
			Service: "graph",
		},
		Reva: Reva{
			Address: "127.0.0.1:9142",
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		Spaces: Spaces{
			WebDavBase:   "https://localhost:9200",
			WebDavPath:   "/dav/spaces/",
			DefaultQuota: "1000000000",
		},
		Identity: Identity{
			Backend: "cs3",
			LDAP: LDAP{
				URI:                      "ldap://localhost:9125",
				BindDN:                   "",
				BindPassword:             "",
				UserBaseDN:               "ou=users,dc=ocis,dc=test",
				UserSearchScope:          "sub",
				UserFilter:               "(objectClass=posixaccount)",
				UserEmailAttribute:       "mail",
				UserDisplayNameAttribute: "displayName",
				UserNameAttribute:        "uid",
				// FIXME: switch this to some more widely available attribute by default
				//        ideally this needs to	be constant for the lifetime of a users
				UserIDAttribute:    "ownclouduuid",
				GroupBaseDN:        "ou=groups,dc=ocis,dc=test",
				GroupSearchScope:   "sub",
				GroupFilter:        "(objectclass=groupOfNames)",
				GroupNameAttribute: "cn",
				GroupIDAttribute:   "cn",
			},
		},
	}
}
