package defaults

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/structs"
	"github.com/owncloud/ocis/v2/services/settings/pkg/config"
	rdefaults "github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	"github.com/pkg/errors"
)

// FullDefaultConfig returns a fully initialized default configuration
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
			CORS: config.CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With", "X-Request-Id"},
				AllowCredentials: true,
			},
		},
		GRPC: config.GRPCConfig{
			Addr:      "127.0.0.1:9191",
			Namespace: "com.owncloud.api",
		},
		SetupDefaultAssignments: false,
		Metadata: config.Metadata{
			GatewayAddress: "com.owncloud.api.storage-system",
			StorageAddress: "com.owncloud.api.storage-system",
			SystemUserIDP:  "internal",
			Cache: &config.Cache{
				Store:          "memory",
				Nodes:          []string{"127.0.0.1:9233"},
				Database:       "settings-cache",
				FileTable:      "settings_files",
				DirectoryTable: "settings_dirs",
				TTL:            time.Minute * 10,
			},
		},
		BundlesPath:       "",
		Bundles:           nil,
		ServiceAccountIDs: []string{"service-user-id"},
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

	if cfg.GRPCClientTLS == nil && cfg.Commons != nil {
		cfg.GRPCClientTLS = structs.CopyOrZeroValue(cfg.Commons.GRPCClientTLS)
	}
	if cfg.GRPC.TLS == nil && cfg.Commons != nil {
		cfg.GRPC.TLS = structs.CopyOrZeroValue(cfg.Commons.GRPCServiceTLS)
	}

	if cfg.Commons != nil {
		cfg.HTTP.TLS = cfg.Commons.HTTPServiceTLS
	}
}

// Sanitize sanitized the configuration
func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}

// LoadBundles loads setting bundles from a file or from defaults
func LoadBundles(cfg *config.Config) error {
	if cfg.BundlesPath != "" {
		data, _ := os.ReadFile(cfg.BundlesPath)
		err := json.Unmarshal(data, &cfg.Bundles)
		if err != nil {
			return errors.Wrapf(err, "Could not load bundles from path %s", cfg.BundlesPath)
		}
	}
	if len(cfg.Bundles) == 0 {
		cfg.Bundles = rdefaults.GenerateBundlesDefaultRoles()
	}
	return nil
}
