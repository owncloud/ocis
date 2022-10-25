package defaults

import (
	"path"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig returns the default config
func DefaultConfig() *config.Config {
	return &config.Config{
		Service: config.Service{
			Name: "settings",
		},
		Debug: config.Debug{
			Addr:   "127.0.0.1:9194",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9190",
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
			Addr:      "127.0.0.1:9191",
			Namespace: "com.owncloud.api",
		},
		StoreType: "metadata", // use metadata or filesystem
		DataPath:  path.Join(defaults.BaseDataPath(), "settings"),
		Asset: config.Asset{
			Path: "",
		},
		SetupDefaultAssignments: false,
		Metadata: config.Metadata{
			GatewayAddress: "127.0.0.1:9215", // system storage
			StorageAddress: "127.0.0.1:9215",
			SystemUserIDP:  "internal",
		},
	}
}

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

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}

	if cfg.Metadata.SystemUserAPIKey == "" && cfg.Commons != nil && cfg.Commons.SystemUserAPIKey != "" {
		cfg.Metadata.SystemUserAPIKey = cfg.Commons.SystemUserAPIKey
	}

	if cfg.Metadata.SystemUserID == "" && cfg.Commons != nil && cfg.Commons.SystemUserID != "" {
		cfg.Metadata.SystemUserID = cfg.Commons.SystemUserID
	}

	if cfg.AdminUserID == "" && cfg.Commons != nil {
		cfg.AdminUserID = cfg.Commons.AdminUserID
	}

	if cfg.MicroGRPCClient == nil {
		cfg.MicroGRPCClient = &shared.MicroGRPCClient{}
		if cfg.Commons != nil && cfg.Commons.MicroGRPCClient != nil {
			cfg.MicroGRPCClient.TLSMode = cfg.Commons.MicroGRPCClient.TLSMode
			cfg.MicroGRPCClient.TLSCACert = cfg.Commons.MicroGRPCClient.TLSCACert
		}
	}
	if cfg.MicroGRPCService == nil {
		cfg.MicroGRPCService = &shared.MicroGRPCService{}
		if cfg.Commons != nil && cfg.Commons.MicroGRPCService != nil {
			cfg.MicroGRPCService.TLSEnabled = cfg.Commons.MicroGRPCService.TLSEnabled
			cfg.MicroGRPCService.TLSCert = cfg.Commons.MicroGRPCService.TLSCert
			cfg.MicroGRPCService.TLSKey = cfg.Commons.MicroGRPCService.TLSKey
		}
	}
}

func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}
