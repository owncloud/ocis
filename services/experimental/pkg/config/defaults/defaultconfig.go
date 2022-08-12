package defaults

import (
	"github.com/owncloud/ocis/v2/services/experimental/pkg/config"
	"strings"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

	return cfg
}

// DefaultConfig sets default service configuration.
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:  "127.0.0.1:9181",
			Token: "",
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9180",
			Namespace: "com.owncloud.experimental",
			Root:      "/experimental",
		},
		Service: config.Service{
			Name: "experimental",
		},
	}
}

// EnsureDefaults ensures that all default values are applied.
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
}

// Sanitize config
func Sanitize(cfg *config.Config) {
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}
