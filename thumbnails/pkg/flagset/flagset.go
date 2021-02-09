package flagset

import (
	"os"
	"path/filepath"

	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Value:       true,
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Value:       true,
			Usage:       "Enable colored logging",
			EnvVars:     []string{"OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9189",
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
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"THUMBNAILS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"THUMBNAILS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"THUMBNAILS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"THUMBNAILS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"THUMBNAILS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "thumbnails",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"THUMBNAILS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9189",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"THUMBNAILS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
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
			Value:       "thumbnails",
			Usage:       "Name of the service",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       "0.0.0.0:9185",
			Usage:       "Address to bind grpc server",
			EnvVars:     []string{"THUMBNAILS_GRPC_ADDR"},
			Destination: &cfg.Server.Address,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAMESPACE"},
			Destination: &cfg.Server.Namespace,
		},
		&cli.StringFlag{
			Name:        "filesystemstorage-root",
			Value:       filepath.Join(os.TempDir(), "ocis-thumbnails/"),
			Usage:       "Root path of the filesystem storage directory",
			EnvVars:     []string{"THUMBNAILS_FILESYSTEMSTORAGE_ROOT"},
			Destination: &cfg.Thumbnail.FileSystemStorage.RootDirectory,
		},
		&cli.StringFlag{
			Name:        "webdavsource-baseurl",
			Value:       "https://localhost:9200/remote.php/webdav/",
			Usage:       "Base url for a webdav api",
			EnvVars:     []string{"THUMBNAILS_WEBDAVSOURCE_BASEURL"},
			Destination: &cfg.Thumbnail.WebDavSource.BaseURL,
		},
		&cli.BoolFlag{
			Name:        "webdavsource-insecure",
			Value:       true,
			Usage:       "Whether to skip certificate checks",
			EnvVars:     []string{"THUMBNAILS_WEBDAVSOURCE_INSECURE"},
			Destination: &cfg.Thumbnail.WebDavSource.Insecure,
		},
		&cli.StringSliceFlag{
			Name:    "thumbnail-resolution",
			Value:   cli.NewStringSlice("16x16", "32x32", "64x64", "128x128", "1920x1080", "3840x2160", "7680x4320"),
			Usage:   "--thumbnail-resolution 16x16 [--thumbnail-resolution 32x32]",
			EnvVars: []string{"THUMBNAILS_RESOLUTIONS"},
		},
	}
}

// ListThumbnailsWithConfig applies the config to the flagset for listing thumbnails services.
func ListThumbnailsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-name",
			Value:       "thumbnails",
			Usage:       "Name of the service",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"THUMBNAILS_GRPC_NAMESPACE"},
			Destination: &cfg.Server.Namespace,
		},
	}
}
