package defaults

import (
	"time"

	"github.com/owncloud/ocis/v2/services/antivirus/pkg/config"
)

// FullDefaultConfig returns a fully initialized default configuration which is needed for doc generation.
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig returns the services default config
func DefaultConfig() *config.Config {
	return &config.Config{
		Debug: config.Debug{
			Addr:  "127.0.0.1:9277",
			Token: "",
		},
		Service: config.Service{
			Name: "antivirus",
		},
		Events: config.Events{
			Endpoint: "127.0.0.1:9233",
			Cluster:  "ocis-cluster",
		},
		Workers:              10,
		InfectedFileHandling: "delete",
		Scanner: config.Scanner{
			Type: "clamav",
			ClamAV: config.ClamAV{
				Socket: "/run/clamav/clamd.ctl",
			},
			ICAP: config.ICAP{
				URL:     "icap://127.0.0.1:1344",
				Service: "avscan",
				Timeout: 5 * time.Minute,
			},
		},
	}
}

// EnsureDefaults adds default values to the configuration if they are not set yet
func EnsureDefaults(cfg *config.Config) {
	if cfg.Log == nil {
		cfg.Log = &config.Log{}
	}

	if cfg.Tracing == nil {
		cfg.Tracing = &config.Tracing{}
	}
}

// Sanitize sanitizes the configuration
func Sanitize(cfg *config.Config) {
}
