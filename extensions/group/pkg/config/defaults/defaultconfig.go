package defaults

import (
	"path/filepath"

	"github.com/owncloud/ocis/extensions/group/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9161",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:     "127.0.0.1:9160",
			Protocol: "tcp",
		},
		Service: config.Service{
			Name: "group",
		},
		GroupMembersCacheExpiration: 5,
		GatewayEndpoint:             "127.0.0.1:9142",
		JWTSecret:                   "Pive-Fumkiu4",
		Driver:                      "ldap",
		Drivers: config.Drivers{
			LDAP: config.LDAPDriver{
				URI:              "ldaps://localhost:9235",
				CACert:           filepath.Join(defaults.BaseDataPath(), "idm", "ldap.crt"),
				Insecure:         false,
				UserBaseDN:       "ou=users,o=libregraph-idm",
				GroupBaseDN:      "ou=groups,o=libregraph-idm",
				UserScope:        "sub",
				GroupScope:       "sub",
				LoginAttributes:  []string{"uid", "mail"},
				UserFilter:       "",
				GroupFilter:      "",
				UserObjectClass:  "inetOrgPerson",
				GroupObjectClass: "groupOfNames",
				BindDN:           "uid=reva,ou=sysusers,o=libregraph-idm",
				BindPassword:     "reva",
				IDP:              "https://localhost:9200",
				UserSchema: config.LDAPUserSchema{
					ID:          "ownclouduuid",
					Mail:        "mail",
					DisplayName: "displayname",
					Username:    "uid",
				},
				GroupSchema: config.LDAPGroupSchema{
					ID:          "ownclouduuid",
					Mail:        "mail",
					DisplayName: "cn",
					Groupname:   "cn",
					Member:      "member",
				},
			},
			JSON: config.JSONDriver{},
			OwnCloudSQL: config.OwnCloudSQLDriver{
				DBUsername:         "owncloud",
				DBPassword:         "secret",
				DBHost:             "mysql",
				DBPort:             3306,
				DBName:             "owncloud",
				IDP:                "https://localhost:9200",
				Nobody:             90,
				JoinUsername:       false,
				JoinOwnCloudUUID:   false,
				EnableMedialSearch: false,
			},
			REST: config.RESTProvider{
				RedisAddr: "localhost:6379",
			},
		},
	}
}

func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for BindEnv.
	if cfg.Logging == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Logging = &config.Logging{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Logging == nil {
		cfg.Logging = &config.Logging{}
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
}

func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
