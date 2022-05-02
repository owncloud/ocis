package defaults

import (
	"github.com/owncloud/ocis/extensions/notifications/pkg/config"
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
		Service: config.Service{
			Name: "notifications",
		},
		Notifications: config.Notifications{
			SMTP: config.SMTP{
				Host:     "127.0.0.1",
				Port:     "1025",
				Sender:   "god@example.com",
				Password: "godisdead",
			},
			Events: config.Events{
				Endpoint:      "127.0.0.1:9233",
				Cluster:       "ocis-cluster",
				ConsumerGroup: "notifications",
			},
			RevaGateway: "127.0.0.1:9142",
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

	if cfg.Notifications.MachineAuthAPIKey == "" && cfg.Commons != nil && cfg.Commons.MachineAuthAPIKey != "" {
		cfg.Notifications.MachineAuthAPIKey = cfg.Commons.MachineAuthAPIKey
	}
}

func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
