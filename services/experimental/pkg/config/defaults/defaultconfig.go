package defaults

import (
	"github.com/owncloud/ocis/v2/services/experimental/pkg/config"
	"strings"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

	return cfg
}

// DefaultConfig sets default service configuration.
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:  "127.0.0.1:9181",
			Token: "",
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9180",
			Namespace: "com.owncloud.experimental",
			Root:      "/experimental",
		},
		Service: config.Service{
			Name: "experimental",
		},
		Events: config.Events{
			Endpoint: "127.0.0.1:9233",
			Cluster:  "ocis-cluster",
		},
		Activities: config.Activities{
			Storage: config.ActivitiesStorage{
				Type: "mem_storage",
				MemStore: config.ActivitiesMemStorage{
					Capacity: 2500,
				},
			},
		},
	}
}

// EnsureDefaults ensures that all default values are applied.
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

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}
}

// Sanitize config
func Sanitize(cfg *config.Config) {
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}
