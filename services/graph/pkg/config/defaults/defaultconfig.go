package defaults

import (
	"path"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:  "127.0.0.1:9124",
			Token: "",
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9120",
			Namespace: "com.owncloud.graph",
			Root:      "/graph",
		},
		Service: config.Service{
			Name: "graph",
		},
		Reva: &config.Reva{
			Address: "127.0.0.1:9142",
		},
		Spaces: config.Spaces{
			WebDavBase:   "https://localhost:9200",
			WebDavPath:   "/dav/spaces/",
			DefaultQuota: "1000000000",
		},
		Identity: config.Identity{
			Backend: "ldap",
			LDAP: config.LDAP{
				URI:                      "ldaps://localhost:9235",
				Insecure:                 false,
				CACert:                   path.Join(defaults.BaseDataPath(), "idm", "ldap.crt"),
				BindDN:                   "uid=libregraph,ou=sysusers,o=libregraph-idm",
				UseServerUUID:            false,
				UsePasswordModExOp:       true,
				WriteEnabled:             true,
				UserBaseDN:               "ou=users,o=libregraph-idm",
				UserSearchScope:          "sub",
				UserFilter:               "",
				UserObjectClass:          "inetOrgPerson",
				UserEmailAttribute:       "mail",
				UserDisplayNameAttribute: "displayName",
				UserNameAttribute:        "uid",
				// FIXME: switch this to some more widely available attribute by default
				//        ideally this needs to	be constant for the lifetime of a users
				UserIDAttribute:    "owncloudUUID",
				GroupBaseDN:        "ou=groups,o=libregraph-idm",
				GroupSearchScope:   "sub",
				GroupFilter:        "",
				GroupObjectClass:   "groupOfNames",
				GroupNameAttribute: "cn",
				GroupIDAttribute:   "owncloudUUID",
			},
		},
		Events: config.Events{
			Endpoint: "127.0.0.1:9233",
			Cluster:  "ocis-cluster",
		},
	}
}

func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for BindEnv.
	if cfg.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Log = &config.Log{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Log == nil {
		cfg.Log = &config.Log{}
	}
	// provide with defaults for shared tracing, since we need a valid destination address for BindEnv.
	if cfg.Tracing == nil && cfg.Commons != nil && cfg.Commons.Tracing != nil {
		cfg.Tracing = &config.Tracing{
			Enabled:   cfg.Commons.Tracing.Enabled,
			Type:      cfg.Commons.Tracing.Type,
			Endpoint:  cfg.Commons.Tracing.Endpoint,
			Collector: cfg.Commons.Tracing.Collector,
		}
	} else if cfg.Tracing == nil {
		cfg.Tracing = &config.Tracing{}
	}

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}
}

func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}
