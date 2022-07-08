package defaults

import (
	"github.com/owncloud/ocis/v2/services/postprocessing/pkg/config"
)

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// NOTE: Most of this configuration is not needed to keep it as simple as possible
// TODO: Clean up unneeded configuration

func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr: "127.0.0.1:9174",
		},
		Service: config.Service{
			Name: "postprocessing",
		},
		Postprocessing: config.Postprocessing{
			Events: config.Events{
				Endpoint:      "127.0.0.1:9233",
				Cluster:       "ocis-cluster",
				ConsumerGroup: "postprocessing",
			},
		},
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
}

func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
