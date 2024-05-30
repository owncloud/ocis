package defaults

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/generators"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
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
	secret, _ := generators.GenerateRandomString(generators.AlphaNumChars, 32) // anything to do with the error?
	return &config.Config{
		Service: config.Service{
			Name: "collaboration",
		},
		App: config.App{
			Name:        "Collabora Online",
			Description: "Open office documents with Collabora Online",
			Icon:        "image-edit",
			LockName:    "com.github.owncloud.collaboration",
		},
		JWTSecret: secret,
		GRPC: config.GRPC{
			Addr:      "0.0.0.0:9301",
			Namespace: "com.owncloud.collaboration",
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9300",
			BindAddr:  "0.0.0.0:9300",
			Namespace: "com.owncloud.collaboration",
			Scheme:    "https",
		},
		Debug: config.Debug{
			Addr:   "127.0.0.1:9304",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		WopiApp: config.WopiApp{
			Addr:     "https://127.0.0.1:8080",
			Insecure: false,
		},
		CS3Api: config.CS3Api{
			Gateway: config.Gateway{
				Name: "com.owncloud.api.gateway",
			},
			DataGateway: config.DataGateway{
				Insecure: false,
			},
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
}

// Sanitize sanitized the configuration
func Sanitize(cfg *config.Config) {
	// sanitize config
}
