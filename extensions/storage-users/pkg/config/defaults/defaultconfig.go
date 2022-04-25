package defaults

import (
	"os"
	"path/filepath"

	"github.com/owncloud/ocis/extensions/storage-users/pkg/config"
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
			Addr:   "127.0.0.1:9159",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: config.GRPCConfig{
			Addr:     "127.0.0.1:9157",
			Protocol: "tcp",
		},
		HTTP: config.HTTPConfig{
			Addr:     "127.0.0.1:9158",
			Protocol: "tcp",
			Prefix:   "data",
		},
		Service: config.Service{
			Name: "storage-users",
		},
		GatewayEndpoint: "127.0.0.1:9142",
		JWTSecret:       "Pive-Fumkiu4",
		TempFolder:      filepath.Join(defaults.BaseDataPath(), "tmp", "users"),
		DataServerURL:   "http://localhost:9158/data",
		MountID:         "1284d238-aa92-42ce-bdc4-0b0000009157",
		Driver:          "ocis",
		Drivers: config.Drivers{
			EOS: config.EOSDriver{
				Root:             "/eos/dockertest/reva",
				ShareFolder:      "/Shares",
				UserLayout:       "{{substr 0 1 .Username}}/{{.Username}}",
				ShadowNamespace:  "",
				UploadsNamespace: "",
				EosBinary:        "/usr/bin/eos",
				XrdcopyBinary:    "/usr/bin/xrdcopy",
				MasterURL:        "root://eos-mgm1.eoscluster.cern.ch:1094",
				GRPCURI:          "",
				SlaveURL:         "root://eos-mgm1.eoscluster.cern.ch:1094",
				CacheDirectory:   os.TempDir(),
				GatewaySVC:       "127.0.0.1:9142",
			},
			Local: config.LocalDriver{
				Root:        filepath.Join(defaults.BaseDataPath(), "storage", "local", "users"),
				ShareFolder: "/Shares",
				UserLayout:  "{{.Username}}",
			},
			OwnCloudSQL: config.OwnCloudSQLDriver{
				Root:          filepath.Join(defaults.BaseDataPath(), "storage", "owncloud"),
				ShareFolder:   "/Shares",
				UserLayout:    "{{.Username}}",
				UploadInfoDir: filepath.Join(defaults.BaseDataPath(), "storage", "uploadinfo"),
				DBUsername:    "owncloud",
				DBPassword:    "owncloud",
				DBHost:        "",
				DBPort:        3306,
				DBName:        "owncloud",
			},
			S3: config.S3Driver{
				Region: "default",
			},
			S3NG: config.S3NGDriver{
				Root:                       filepath.Join(defaults.BaseDataPath(), "storage", "users"),
				ShareFolder:                "/Shares",
				UserLayout:                 "{{.Id.OpaqueId}}",
				Region:                     "default",
				PersonalSpaceAliasTemplate: "{{.SpaceType}}/{{.User.Username | lower}}",
				GeneralSpaceAliasTemplate:  "{{.SpaceType}}/{{.SpaceName | replace \" \" \"-\" | lower}}",
			},
			OCIS: config.OCISDriver{
				Root:                       filepath.Join(defaults.BaseDataPath(), "storage", "users"),
				ShareFolder:                "/Shares",
				UserLayout:                 "{{.Id.OpaqueId}}",
				PersonalSpaceAliasTemplate: "{{.SpaceType}}/{{.User.Username | lower}}",
				GeneralSpaceAliasTemplate:  "{{.SpaceType}}/{{.SpaceName | replace \" \" \"-\" | lower}}",
			},
		},
		Events: config.Events{
			Addr:      "127.0.0.1:9233",
			ClusterID: "ocis-cluster",
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
}

func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
