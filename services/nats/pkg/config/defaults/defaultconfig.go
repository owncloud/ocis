package defaults

import (
	"path/filepath"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/nats/pkg/config"
)

// DefaultConfig returns a basic default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr: "127.0.0.1:9234",
		},
		Service: config.Service{
			Name: "nats",
		},
		Nats: config.Nats{
			Host:      "127.0.0.1",
			Port:      9233,
			ClusterID: "ocis-cluster",
			StoreDir:  filepath.Join(defaults.BaseDataPath(), "nats"),
			TLSCert:   filepath.Join(defaults.BaseDataPath(), "nats/tls.crt"),
			TLSKey:    filepath.Join(defaults.BaseDataPath(), "nats/tls.key"),
			EnableTLS: false,
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
}
