package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"

	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/urfave/cli/v2"
)

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9189"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"THUMBNAILS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-file",
			Usage:       "Enable log to file",
			EnvVars:     []string{"THUMBNAILS_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"THUMBNAILS_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"THUMBNAILS_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"THUMBNAILS_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"THUMBNAILS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"THUMBNAILS_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"THUMBNAILS_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"THUMBNAILS_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"THUMBNAILS_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "thumbnails"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"THUMBNAILS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "0.0.0.0:9189"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"THUMBNAILS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"THUMBNAILS_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"THUMBNAILS_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"THUMBNAILS_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "grpc-name",
			Value:       flags.OverrideDefaultString(cfg.Server.Name, "thumbnails"),
			Usage:       "Name of the service",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       flags.OverrideDefaultString(cfg.Server.Address, "0.0.0.0:9185"),
			Usage:       "Address to bind grpc server",
			EnvVars:     []string{"THUMBNAILS_GRPC_ADDR"},
			Destination: &cfg.Server.Address,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       flags.OverrideDefaultString(cfg.Server.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAMESPACE"},
			Destination: &cfg.Server.Namespace,
		},
		&cli.StringFlag{
			Name:        "filesystemstorage-root",
			Value:       "/var/tmp/ocis/thumbnails",
			Usage:       "Root path of the filesystem storage directory",
			EnvVars:     []string{"THUMBNAILS_FILESYSTEMSTORAGE_ROOT"},
			Destination: &cfg.Thumbnail.FileSystemStorage.RootDirectory,
		},
		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Thumbnail.RevaGateway, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Thumbnail.RevaGateway,
		},
		&cli.BoolFlag{
			Name:        "webdavsource-insecure",
			Value:       flags.OverrideDefaultBool(cfg.Thumbnail.WebdavAllowInsecure, true),
			Usage:       "Whether to skip certificate checks",
			EnvVars:     []string{"THUMBNAILS_WEBDAVSOURCE_INSECURE"},
			Destination: &cfg.Thumbnail.WebdavAllowInsecure,
		},
		&cli.StringSliceFlag{
			Name:    "thumbnail-resolution",
			Value:   cli.NewStringSlice("16x16", "32x32", "64x64", "128x128", "1920x1080", "3840x2160", "7680x4320"),
			Usage:   "--thumbnail-resolution 16x16 [--thumbnail-resolution 32x32]",
			EnvVars: []string{"THUMBNAILS_RESOLUTIONS"},
		},
		&cli.StringFlag{
			Name:        "webdav-namespace",
			Value:       flags.OverrideDefaultString(cfg.Thumbnail.WebdavNamespace, "/home"),
			Usage:       "Namespace prefix for the webdav endpoint",
			EnvVars:     []string{"STORAGE_WEBDAV_NAMESPACE"},
			Destination: &cfg.Thumbnail.WebdavNamespace,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode",
		},
	}
}

// ListThumbnailsWithConfig applies the config to the flagset for listing thumbnails services.
func ListThumbnailsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-name",
			Value:       flags.OverrideDefaultString(cfg.Server.Name, "thumbnails"),
			Usage:       "Name of the service",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       flags.OverrideDefaultString(cfg.Server.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAMESPACE"},
			Destination: &cfg.Server.Namespace,
		},
	}
}
