package defaults

import (
	"path/filepath"

	"github.com/owncloud/ocis/extensions/auth-basic/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Service: config.ServiceConfig{
			Network:         "tcp",
			Address:         "127.0.0.1:9146",
			GatewayEndpoint: "127.0.0.1:9142",
			JWTSecret:       "Pive-Fumkiu4",
			AuthManager:     "ldap",
			AuthManagers: config.AuthManagers{
				LDAP: config.LDAPManager{
					URI:              "ldaps:localhost:9126",
					CACert:           filepath.Join(defaults.BaseDataPath(), "ldap", "ldap.crt"),
					Insecure:         false,
					UserBaseDN:       "dc=ocis,dc=test",
					GroupBaseDN:      "dc=ocis,dc=test",
					LoginAttributes:  []string{"cn", "mail"},
					UserFilter:       "",
					GroupFilter:      "",
					UserObjectClass:  "posixAccount",
					GroupObjectClass: "posixGroup",
					BindDN:           "cn=reva,ou=sysusers,dc=ocis,dc=test",
					BindPassword:     "reva",
					IDP:              "https://localhost:9200",
					UserSchema: config.LDAPUserSchema{
						ID:          "ownclouduuid",
						Mail:        "mail",
						DisplayName: "displayname",
						Username:    "cn",
					},
					GroupSchema: config.LDAPGroupSchema{
						ID:          "cn",
						Mail:        "mail",
						DisplayName: "cn",
						Groupname:   "cn",
						Member:      "cn",
					},
				},
			},
		},
	}
}

func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for BindEnv.
	if cfg.Logging == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Logging = &config.LoggingConfig{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Logging == nil {
		cfg.Logging = &config.LoggingConfig{}
	}
	// provide with defaults for shared tracing, since we need a valid destination address for BindEnv.
	if cfg.Tracing == nil && cfg.Commons != nil && cfg.Commons.Tracing != nil {
		cfg.Tracing = &config.TracingConfig{
			Enabled:   cfg.Commons.Tracing.Enabled,
			Type:      cfg.Commons.Tracing.Type,
			Endpoint:  cfg.Commons.Tracing.Endpoint,
			Collector: cfg.Commons.Tracing.Collector,
		}
	} else if cfg.Tracing == nil {
		cfg.Tracing = &config.TracingConfig{}
	}
}
