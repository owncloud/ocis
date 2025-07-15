package defaults

import (
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/ocis-pkg/structs"
	"github.com/owncloud/ocis/v2/services/gateway/pkg/config"
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
			Addr:   "127.0.0.1:9143",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:      "127.0.0.1:9142",
			Namespace: "com.owncloud.api",
			Protocol:  "tcp",
		},
		Service: config.Service{
			Name: "gateway",
		},
		Reva:                       shared.DefaultRevaConfig(),
		CommitShareToStorageGrant:  true,
		ShareFolder:                "Shares",
		DisableHomeCreationOnLogin: true,
		TransferExpires:            24 * 60 * 60,
		Cache: config.Cache{
			ProviderCacheStore:      "noop",
			ProviderCacheNodes:      []string{"127.0.0.1:9233"},
			ProviderCacheDatabase:   "cache-providers",
			ProviderCacheTTL:        300 * time.Second,
			CreateHomeCacheStore:    "memory",
			CreateHomeCacheNodes:    []string{"127.0.0.1:9233"},
			CreateHomeCacheDatabase: "cache-createhome",
			CreateHomeCacheTTL:      300 * time.Second,
		},

		FrontendPublicURL: "https://localhost:9200",

		AppRegistryEndpoint:       "com.owncloud.api.app-registry",
		AuthAppEndpoint:           "com.owncloud.api.auth-app",
		AuthBasicEndpoint:         "com.owncloud.api.auth-basic",
		AuthMachineEndpoint:       "com.owncloud.api.auth-machine",
		AuthServiceEndpoint:       "com.owncloud.api.auth-service",
		GroupsEndpoint:            "com.owncloud.api.groups",
		PermissionsEndpoint:       "com.owncloud.api.settings",
		SharingEndpoint:           "com.owncloud.api.sharing",
		StoragePublicLinkEndpoint: "com.owncloud.api.storage-publiclink",
		StorageSharesEndpoint:     "com.owncloud.api.storage-shares",
		StorageUsersEndpoint:      "com.owncloud.api.storage-users",
		UsersEndpoint:             "com.owncloud.api.users",
		OCMEndpoint:               "com.owncloud.api.ocm",

		StorageRegistry: config.StorageRegistry{
			Driver: "spaces",
			JSON:   "",
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

	if cfg.Reva == nil && cfg.Commons != nil {
		cfg.Reva = structs.CopyOrZeroValue(cfg.Commons.Reva)
	}

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}

	if cfg.TransferSecret == "" && cfg.Commons != nil && cfg.Commons.TransferSecret != "" {
		cfg.TransferSecret = cfg.Commons.TransferSecret
	}

	if cfg.GRPC.TLS == nil && cfg.Commons != nil {
		cfg.GRPC.TLS = structs.CopyOrZeroValue(cfg.Commons.GRPCServiceTLS)
	}
}

// Sanitize sanitized the configuration
func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
