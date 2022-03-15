package defaults

import (
	"path"
	"strings"

	"github.com/owncloud/ocis/accounts/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
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
			Addr:   "127.0.0.1:9182",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9181",
			Namespace: "com.owncloud.web",
			Root:      "/",
			CacheTTL:  604800, // 7 days
			CORS: config.CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		GRPC: config.GRPC{
			Addr:      "127.0.0.1:9180",
			Namespace: "com.owncloud.api",
		},
		Service: config.Service{
			Name: "accounts",
		},
		Asset: config.Asset{},
		TokenManager: config.TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		HashDifficulty:     11,
		DemoUsersAndGroups: true,
		Repo: config.Repo{
			Backend: "CS3",
			Disk: config.Disk{
				Path: path.Join(defaults.BaseDataPath(), "accounts"),
			},
			CS3: config.CS3{
				ProviderAddr: "localhost:9215",
			},
		},
		Index: config.Index{
			UID: config.UIDBound{
				Lower: 0,
				Upper: 1000,
			},
			GID: config.GIDBound{
				Lower: 0,
				Upper: 1000,
			},
		},
		ServiceUser: config.ServiceUser{
			UUID:     "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
			Username: "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
			UID:      0,
			GID:      0,
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
}

func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
	cfg.Repo.Backend = strings.ToLower(cfg.Repo.Backend)
}
