package defaults

import (
	"path/filepath"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/ocis-pkg/structs"
	"github.com/owncloud/ocis/v2/services/storage-system/pkg/config"
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
			Addr:   "127.0.0.1:9217",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:      "127.0.0.1:9215",
			Namespace: "com.owncloud.api",
			Protocol:  "tcp",
		},
		HTTP: config.HTTPConfig{
			Addr:      "127.0.0.1:9216",
			Namespace: "com.owncloud.web",
			Protocol:  "tcp",
		},
		Service: config.Service{
			Name: "storage-system",
		},
		Reva:          shared.DefaultRevaConfig(),
		DataServerURL: "http://localhost:9216/data",
		Driver:        "ocis",
		Drivers: config.Drivers{
			OCIS: config.OCISDriver{
				Root:                    filepath.Join(defaults.BaseDataPath(), "storage", "metadata"),
				MaxAcquireLockCycles:    20,
				LockCycleDurationFactor: 30,
			},
		},
		FileMetadataCache: config.Cache{
			Store:    "memory",
			Nodes:    []string{"127.0.0.1:9233"},
			Database: "storage-system",
			TTL:      24 * 60 * time.Second,
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

	if cfg.SystemUserAPIKey == "" && cfg.Commons != nil && cfg.Commons.SystemUserAPIKey != "" {
		cfg.SystemUserAPIKey = cfg.Commons.SystemUserAPIKey
	}

	if cfg.SystemUserID == "" && cfg.Commons != nil && cfg.Commons.SystemUserID != "" {
		cfg.SystemUserID = cfg.Commons.SystemUserID
	}

	if cfg.GRPC.TLS == nil && cfg.Commons != nil {
		cfg.GRPC.TLS = structs.CopyOrZeroValue(cfg.Commons.GRPCServiceTLS)
	}

}

// Sanitize sanitized the configuration
func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
