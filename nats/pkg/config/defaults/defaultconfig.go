package defaults

import (
	"path"

	"github.com/owncloud/ocis/nats/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// NOTE: Most of this configuration is not needed to keep it as simple as possible
// TODO: Clean up unneeded configuration

func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()

	EnsureDefaults(cfg)
	Sanitize(cfg)

	return cfg
}

func DefaultConfig() *config.Config {
	return &config.Config{
		Service: config.Service{
			Name: "nats",
		},
		Nats: config.Nats{
			Host:      "127.0.0.1",
			Port:      9233,
			ClusterID: "ocis-cluster",
			StoreDir:  path.Join(defaults.BaseDataPath(), "nats"),
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
