package defaults

import (
	"github.com/owncloud/ocis/extensions/gateway/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

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
			Addr:     "127.0.0.1:9142",
			Protocol: "tcp",
		},
		Service: config.Service{
			Name: "gateway",
		},
		GatewayEndpoint: "127.0.0.1:9142",
		JWTSecret:       "Pive-Fumkiu4",

		CommitShareToStorageGrant:  true,
		CommitShareToStorageRef:    true,
		ShareFolder:                "Shares",
		DisableHomeCreationOnLogin: true,
		TransferSecret:             "replace-me-with-a-transfer-secret",
		TransferExpires:            24 * 60 * 60,
		HomeMapping:                "",
		EtagCacheTTL:               0,

		UsersEndpoint:             "localhost:9144",
		GroupsEndpoint:            "localhost:9160",
		PermissionsEndpoint:       "localhost:9191",
		SharingEndpoint:           "localhost:9150",
		DataGatewayPublicURL:      "",
		FrontendPublicURL:         "https://localhost:9200",
		AuthBasicEndpoint:         "localhost:9146",
		AuthBearerEndpoint:        "localhost:9148",
		AuthMachineEndpoint:       "localhost:9166",
		StoragePublicLinkEndpoint: "localhost:9178",
		StorageUsersEndpoint:      "localhost:9157",
		StorageSharesEndpoint:     "localhost:9154",

		StorageRegistry: config.StorageRegistry{
			Driver: "spaces",
			JSON:   "",
		},
		AppRegistry: config.AppRegistry{
			MimetypesJSON: "",
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
