package defaults

import (
	"github.com/owncloud/ocis/v2/services/hub/pkg/config"
	"strings"
)

// DefaultConfig returns the default config
func DefaultConfig() *config.Config {
	return &config.Config{
		Service: config.Service{
			Name: "hub",
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9180",
			Namespace: "com.owncloud.web",
			Root:      "/",
		},
	}
}

func EnsureDefaults(cfg *config.Config) {
	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}
}

func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}
}
