package defaults

import (
	"github.com/owncloud/ocis/extensions/gateway/pkg/config"
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
		Reva: &config.Reva{
			Address: "127.0.0.1:9142",
		},

		CommitShareToStorageGrant:  true,
		CommitShareToStorageRef:    true,
		ShareFolder:                "Shares",
		DisableHomeCreationOnLogin: true,
		TransferExpires:            24 * 60 * 60,
		HomeMapping:                "",
		EtagCacheTTL:               0,

		FrontendPublicURL: "https://localhost:9200",

		AppRegistryEndpoint:       "localhost:9242",
		AuthBasicEndpoint:         "localhost:9146",
		AuthBearerEndpoint:        "localhost:9148",
		AuthMachineEndpoint:       "localhost:9166",
		GroupsEndpoint:            "localhost:9160",
		PermissionsEndpoint:       "localhost:9191",
		SharingEndpoint:           "localhost:9150",
		StoragePublicLinkEndpoint: "localhost:9178",
		StorageSharesEndpoint:     "localhost:9154",
		StorageUsersEndpoint:      "localhost:9157",
		UsersEndpoint:             "localhost:9144",

		StorageRegistry: config.StorageRegistry{
			Driver: "spaces",
			JSON:   "",
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

	if cfg.TransferSecret == "" && cfg.Commons != nil && cfg.Commons.TransferSecret != "" {
		cfg.TransferSecret = cfg.Commons.TransferSecret
	}
}

func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
