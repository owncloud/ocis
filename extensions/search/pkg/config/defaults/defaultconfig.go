package defaults

import (
	"path"

	"github.com/owncloud/ocis/extensions/search/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:  "127.0.0.1:9224",
			Token: "",
		},
		GRPC: config.GRPC{
			Addr:      "127.0.0.1:9220",
			Namespace: "com.owncloud.api",
		},
		Service: config.Service{
			Name: "search",
		},
		Datapath: path.Join(defaults.BaseDataPath(), "search"),
		Reva: config.Reva{
			Address: "127.0.0.1:9142",
		},
		TokenManager: config.TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		Events: config.Events{
			Endpoint:      "127.0.0.1:9233",
			Cluster:       "ocis-cluster",
			ConsumerGroup: "search",
		},
		MachineAuthAPIKey: "change-me-please",
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
}

func Sanitize(cfg *config.Config) {
	// no http endpoint to be sanitized
}
