package defaults

import (
	"log"

	"github.com/owncloud/ocis/extensions/auth-machine/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9167",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:     "127.0.0.1:9166",
			Protocol: "tcp",
		},
		Service: config.Service{
			Name: "auth-machine",
		},
		Reva: &config.Reva{
			Address: "127.0.0.1:9142",
		},
		AuthProvider: "ldap",
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

	if cfg.Reva == nil && cfg.Commons != nil && cfg.Commons.Reva != nil {
		cfg.Reva = &config.Reva{
			Address: cfg.Commons.Reva.Address,
		}
	} else if cfg.Reva == nil {
		cfg.Reva = &config.Reva{}
	}

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}

	if cfg.AuthProviders.Machine.APIKey == "" && cfg.Commons != nil && cfg.Commons.MachineAuthAPIKey != "" {
		cfg.AuthProviders.Machine.APIKey = cfg.Commons.MachineAuthAPIKey
	} else if cfg.AuthProviders.Machine.APIKey == "" {
		log.Fatalf("machine auth api key is not set up properly, bailing out (%s)", cfg.Service.Name)
	}
}

func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
