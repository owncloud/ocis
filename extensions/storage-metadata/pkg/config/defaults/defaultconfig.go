package defaults

import (
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/extensions/storage-metadata/pkg/config"
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
			Addr:   "127.0.0.1:9217",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:     "127.0.0.1:9215",
			Protocol: "tcp",
		},
		HTTP: config.HTTPConfig{
			Addr:     "127.0.0.1:9216",
			Protocol: "tcp",
		},
		Service: config.Service{
			Name: "storage-metadata",
		},
		Reva: &config.Reva{
			Address: "127.0.0.1:9142",
		},
		TempFolder:    filepath.Join(defaults.BaseDataPath(), "tmp", "metadata"),
		DataServerURL: "http://localhost:9216/data",
		Driver:        "ocis",
		Drivers: config.Drivers{
			EOS: config.EOSDriver{
				Root:                "/eos/dockertest/reva",
				UserLayout:          "{{substr 0 1 .Username}}/{{.Username}}",
				ShadowNamespace:     "",
				UploadsNamespace:    "",
				EosBinary:           "/usr/bin/eos",
				XrdcopyBinary:       "/usr/bin/xrdcopy",
				MasterURL:           "root://eos-mgm1.eoscluster.cern.ch:1094",
				GRPCURI:             "",
				SlaveURL:            "root://eos-mgm1.eoscluster.cern.ch:1094",
				CacheDirectory:      os.TempDir(),
				EnableLogging:       false,
				ShowHiddenSysFiles:  false,
				ForceSingleUserMode: false,
				UseKeytab:           false,
				SecProtocol:         "",
				Keytab:              "",
				SingleUsername:      "",
				GatewaySVC:          "127.0.0.1:9142",
			},
			Local: config.LocalDriver{
				Root: filepath.Join(defaults.BaseDataPath(), "storage", "local", "metadata"),
			},
			S3: config.S3Driver{
				Region: "default",
			},
			S3NG: config.S3NGDriver{
				Root:                filepath.Join(defaults.BaseDataPath(), "storage", "metadata"),
				UserLayout:          "{{.Id.OpaqueId}}",
				Region:              "default",
				PermissionsEndpoint: "127.0.0.1:9191",
			},
			OCIS: config.OCISDriver{
				Root:                filepath.Join(defaults.BaseDataPath(), "storage", "metadata"),
				UserLayout:          "{{.Id.OpaqueId}}",
				PermissionsEndpoint: "127.0.0.1:9191",
			},
		},
	}
}

func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for BindEnv.
	if cfg.Logging == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Logging = &config.Logging{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Logging == nil {
		cfg.Logging = &config.Logging{}
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

	if cfg.Reva == nil && cfg.Commons != nil && cfg.Commons.Reva != nil {
		cfg.Reva = &config.Reva{
			Address: cfg.Commons.Reva.Address,
		}
	} else if cfg.Reva == nil {
		cfg.Reva = &config.Reva{}
	}

	if cfg.TokenManager == nil && cfg.Commons != nil && cfg.Commons.TokenManager != nil {
		cfg.TokenManager = &config.TokenManager{
			JWTSecret: cfg.Commons.TokenManager.JWTSecret,
		}
	} else if cfg.TokenManager == nil {
		cfg.TokenManager = &config.TokenManager{}
	}
}

func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
