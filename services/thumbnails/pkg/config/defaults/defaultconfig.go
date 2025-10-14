package defaults

import (
	"path"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/ocis-pkg/structs"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/config"
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
			Addr:   "127.0.0.1:9189",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:                  "127.0.0.1:9185",
			Namespace:             "com.owncloud.api",
			MaxConcurrentRequests: 0,
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9186",
			Root:      "/thumbnails",
			Namespace: "com.owncloud.web",
			CORS: config.CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With", "X-Request-Id", "Cache-Control"},
				AllowCredentials: true,
			},
		},
		Service: config.Service{
			Name: "thumbnails",
		},
		Thumbnail: config.Thumbnail{
			Resolutions: []string{"16x16", "32x32", "64x64", "128x128", "1080x1920", "1920x1080", "2160x3840", "3840x2160", "4320x7680", "7680x4320"},
			FileSystemStorage: config.FileSystemStorage{
				RootDirectory: path.Join(defaults.BaseDataPath(), "thumbnails"),
			},
			WebdavAllowInsecure:   false,
			RevaGateway:           shared.DefaultRevaConfig().Address,
			CS3AllowInsecure:      false,
			DataEndpoint:          "http://127.0.0.1:9186/thumbnails/data",
			MaxInputWidth:         7680,
			MaxInputHeight:        7680,
			MaxInputImageFileSize: "50MB",
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
	// nothing to sanitize here atm
	if len(cfg.Thumbnail.Resolutions) == 1 && strings.Contains(cfg.Thumbnail.Resolutions[0], ",") {
		cfg.Thumbnail.Resolutions = strings.Split(cfg.Thumbnail.Resolutions[0], ",")
	}
}
