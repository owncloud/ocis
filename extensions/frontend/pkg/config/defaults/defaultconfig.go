package defaults

import (
	"github.com/owncloud/ocis/extensions/frontend/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9141",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTPConfig{
			Addr:     "127.0.0.1:9140",
			Protocol: "tcp",
			Prefix:   "",
		},
		Service: config.Service{
			Name: "frontend",
		},
		GatewayEndpoint:          "127.0.0.1:9142",
		JWTSecret:                "Pive-Fumkiu4",
		PublicURL:                "https://localhost:9200",
		EnableFavorites:          false,
		EnableProjectSpaces:      true,
		UploadMaxChunkSize:       1e+8,
		UploadHTTPMethodOverride: "",
		DefaultUploadProtocol:    "tus",
		TransferSecret:           "replace-me-with-a-transfer-secret",
		Checksums: config.Checksums{
			SupportedTypes:      []string{"sha1", "md5", "adler32"},
			PreferredUploadType: "",
		},
		AppProvider: config.AppProvider{
			Prefix:   "",
			Insecure: false,
		},
		Archiver: config.Archiver{
			Insecure: false,
			Prefix:   "archiver",
		},
		DataGateway: config.DataGateway{
			Prefix: "data",
		},
		OCS: config.OCS{
			Prefix:                  "ocs",
			SharePrefix:             "/Shares",
			HomeNamespace:           "/users/{{.Id.OpaqueId}}",
			CacheWarmupDriver:       "",
			AdditionalInfoAttribute: "{{.Mail}}",
			ResourceInfoCacheTTL:    0,
		},
		AuthMachine: config.AuthMachine{
			APIKey: "change-me-please",
		},
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
