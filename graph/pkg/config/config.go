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

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `ocisConfig:"addr"`
	Namespace string `ocisConfig:"namespace"`
	Root      string `ocisConfig:"root"`
}

// Server configures a server.
type Server struct {
	Version string `ocisConfig:"version"`
	Name    string `ocisConfig:"name"`
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

type Spaces struct {
	WebDavBase   string `ocisConfig:"webdav_base"`
	WebDavPath   string `ocisConfig:"webdav_path"`
	DefaultQuota string `ocisConfig:"default_quota"`
}

type LDAP struct {
	URI          string `ocisConfig:"uri"`
	BindDN       string `ocisConfig:"bind_dn"`
	BindPassword string `ocisConfig:"bind_password"`

	UserBaseDN               string `ocisConfig:"user_base_dn"`
	UserSearchScope          string `ocisConfig:"user_search_scope"`
	UserFilter               string `ocisConfig:"user_filter"`
	UserEmailAttribute       string `ocisConfig:"user_mail_attribute"`
	UserDisplayNameAttribute string `ocisConfig:"user_displayname_attribute"`
	UserNameAttribute        string `ocisConfig:"user_name_attribute"`
	UserIDAttribute          string `ocisConfig:"user_id_attribute"`

	GroupBaseDN        string `ocisConfig:"group_base_dn"`
	GroupSearchScope   string `ocisConfig:"group_search_scope"`
	GroupFilter        string `ocisConfig:"group_filter"`
	GroupNameAttribute string `ocisConfig:"group_name_attribute"`
	GroupIDAttribute   string `ocisConfig:"group_id_attribute"`
}

type Identity struct {
	Backend string `ocisConfig:"backend"`
	LDAP    LDAP   `ocisConfig:"ldap"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File         string       `ocisConfig:"file"`
	Log          *shared.Log  `ocisConfig:"log"`
	Debug        Debug        `ocisConfig:"debug"`
	HTTP         HTTP         `ocisConfig:"http"`
	Server       Server       `ocisConfig:"server"`
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
			Namespace: "com.owncloud.web",
			Root:      "/graph",
		},
		Server: Server{},
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
