package defaults

import (
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/ocis-pkg/structs"
	"github.com/owncloud/ocis/v2/services/ocm/pkg/config"
)

// FullDefaultConfig returns the full default config
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig return the default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:   "127.0.0.1:9281",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: config.HTTPConfig{
			Addr:      "127.0.0.1:9280",
			Namespace: "com.owncloud.web",
			Protocol:  "tcp",
			Prefix:    "",
			CORS: config.CORS{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{
					"OPTIONS",
					"HEAD",
					"GET",
					"PUT",
					"POST",
					"DELETE",
					"MKCOL",
					"PROPFIND",
					"PROPPATCH",
					"MOVE",
					"COPY",
					"REPORT",
					"SEARCH",
				},
				AllowedHeaders: []string{
					"Origin",
					"Accept",
					"Content-Type",
					"Depth",
					"Authorization",
					"Ocs-Apirequest",
					"If-None-Match",
					"If-Match",
					"Destination",
					"Overwrite",
					"X-Request-Id",
					"X-Requested-With",
					"Tus-Resumable",
					"Tus-Checksum-Algorithm",
					"Upload-Concat",
					"Upload-Length",
					"Upload-Metadata",
					"Upload-Defer-Length",
					"Upload-Expires",
					"Upload-Checksum",
					"Upload-Offset",
					"X-HTTP-Method-Override",
					"Cache-Control",
				},
				AllowCredentials: true,
			},
		},
		GRPC: config.GRPCConfig{
			Addr:      "127.0.0.1:9282",
			Namespace: "com.owncloud.api",
		},
		Reva: shared.DefaultRevaConfig(),
		Service: config.Service{
			Name: "ocm",
		},
		ScienceMesh: config.ScienceMesh{
			Prefix: "sciencemesh",
		},
		OCMD: config.OCMD{
			Prefix: "ocm",
		},
		OCMInviteManager: config.OCMInviteManager{
			Driver: "json",
			Drivers: config.OCMInviteManagerDrivers{
				JSON: config.OCMInviteManagerJSONDriver{
					File: filepath.Join(defaults.BaseDataPath(), "storage", "ocminvites.json"),
				},
			},
			Insecure: false,
		},
		OCMProviderAuthorizerDriver: "json",
		OCMProviderAuthorizerDrivers: config.OCMProviderAuthorizerDrivers{
			JSON: config.OCMProviderAuthorizerJSONDriver{
				Providers: filepath.Join(defaults.BaseDataPath(), "storage", "ocmproviders.json"),
			},
		},
		OCMShareProvider: config.OCMShareProvider{
			Driver: "json",
			Drivers: config.OCMShareProviderDrivers{
				JSON: config.OCMShareProviderJSONDriver{
					File: filepath.Join(defaults.BaseDataPath(), "storage", "ocmshares.json"),
				},
			},
			Insecure: false,
		},
		OCMCore: config.OCMCore{
			Driver: "json",
			Drivers: config.OCMCoreDrivers{
				JSON: config.OCMCoreJSONDriver{
					File: filepath.Join(defaults.BaseDataPath(), "storage", "ocmshares.json"),
				},
			},
		},
	}
}

// EnsureDefaults ensures the config contains default values
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

	if cfg.GRPCClientTLS == nil && cfg.Commons != nil {
		cfg.GRPCClientTLS = structs.CopyOrZeroValue(cfg.Commons.GRPCClientTLS)
	}

	if cfg.GRPC.TLS == nil && cfg.Commons != nil {
		cfg.GRPC.TLS = structs.CopyOrZeroValue(cfg.Commons.GRPCServiceTLS)
	}
}

// Sanitize sanitizes the config
func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
