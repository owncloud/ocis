package defaults

import (
	"path"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/idm/pkg/config"
)

// FullDefaultConfig returns a fully initialized default configuration
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig returns a basic default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9239",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		Service: config.Service{
			Name: "idm",
		},
		CreateDemoUsers:    false,
		DemoUsersIssuerUrl: "https://localhost:9200",
		IDM: config.Settings{
			LDAPSAddr:    "127.0.0.1:9235",
			Cert:         path.Join(defaults.BaseDataPath(), "idm", "ldap.crt"),
			Key:          path.Join(defaults.BaseDataPath(), "idm", "ldap.key"),
			DatabasePath: path.Join(defaults.BaseDataPath(), "idm", "ocis.boltdb"),
		},
	}
}

// EnsureDefaults adds default values to the configuration if they are not set yet
func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for "envdecode".
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
	// provide with defaults for shared tracing, since we need a valid destination address for "envdecode".
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

	if cfg.AdminUserID == "" && cfg.Commons != nil {
		cfg.AdminUserID = cfg.Commons.AdminUserID
	}
}

// Sanitize sanitizes the configuration
func Sanitize(cfg *config.Config) {
	// nothing to sanitize here
}
