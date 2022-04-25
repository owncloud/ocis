package defaults

import (
	"github.com/owncloud/ocis/extensions/ocdav/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9163",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTPConfig{
			Addr:     "127.0.0.1:0", // :0 to pick any free local port
			Protocol: "tcp",
			Prefix:   "",
		},
		Service: config.Service{
			Name: "ocdav",
		},
		GatewayEndpoint: "127.0.0.1:9142",
		JWTSecret:       "Pive-Fumkiu4",
		WebdavNamespace: "/users/{{.Id.OpaqueId}}",
		FilesNamespace:  "/users/{{.Id.OpaqueId}}",
		SharesNamespace: "/Shares",
		PublicURL:       "https://localhost:9200",
		Insecure:        false,
		Timeout:         84300,
		Middleware: config.Middleware{
			Auth: config.Auth{
				CredentialsByUserAgent: map[string]string{},
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
