package defaults

import (
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/authz/pkg/config"
	"strings"
	"time"
)

// DefaultConfig returns the default config
func DefaultConfig() *config.Config {
	return &config.Config{
		Service: config.Service{
			Name: "authz",
		},
		GRPC: config.GRPC{
			Addr:      "127.0.0.1:9180",
			Namespace: "com.owncloud.api",
		},
		HTTP: config.HTTP{
			Addr:      "127.0.0.1:9181",
			Namespace: "com.owncloud.web",
			Root:      "/",
		},
		Reva: shared.DefaultRevaConfig(),
		Events: config.Events{
			Endpoint:  "127.0.0.1:9233",
			Cluster:   "ocis-cluster",
			EnableTLS: false,
		},
		OPA: config.OPA{
			Enabled: true,
			Policies: []string{
				"services/authz/pkg/config/policies/ocis.authz.rego",
			},
			Timeout: 5,
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
	if cfg.MachineAuthAPIKey == "" && cfg.Commons != nil && cfg.Commons.MachineAuthAPIKey != "" {
		cfg.MachineAuthAPIKey = cfg.Commons.MachineAuthAPIKey
	}

	if cfg.Reva == nil && cfg.Commons != nil && cfg.Commons.Reva != nil {
		cfg.Reva = &shared.Reva{
			Address: cfg.Commons.Reva.Address,
			TLS:     cfg.Commons.Reva.TLS,
		}
	} else if cfg.Reva == nil {
		cfg.Reva = &shared.Reva{}
	}
}

func Sanitize(cfg *config.Config) {
	// sanitize config
	if cfg.HTTP.Root != "/" {
		cfg.HTTP.Root = strings.TrimSuffix(cfg.HTTP.Root, "/")
	}

	// convert timeout to millisecond
	// the config is in seconds, therefore we need multiply it.
	cfg.OPA.Timeout = cfg.OPA.Timeout * int(time.Second)
}
