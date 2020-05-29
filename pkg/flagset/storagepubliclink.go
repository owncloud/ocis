package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// StoragePublicLink applies cfg to the root flagset
func StoragePublicLink(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"REVA_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       "jaeger",
			Usage:       "Tracing backend type",
			EnvVars:     []string{"REVA_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       "",
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"REVA_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       "",
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"REVA_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       "reva",
			Usage:       "Service name for tracing",
			EnvVars:     []string{"REVA_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:10053",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_DEBUG_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       "",
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"REVA_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"REVA_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"REVA_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Shared jwt secret for reva service communication",
			EnvVars:     []string{"REVA_JWT_SECRET"},
			Destination: &cfg.Reva.JWTSecret,
		},
		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_NETWORK"},
			Destination: &cfg.Reva.StoragePublicLink.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_PROTOCOL"},
			Destination: &cfg.Reva.StoragePublicLink.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "localhost:10054",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.Addr,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("storageprovider"),
			Usage:   "--service storageprovider [--service otherservice]",
			EnvVars: []string{"REVA_STORAGE_PUBLIC_LINK_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "public_share_provider_addr",
			Value:       "localhost:9150",
			Usage:       "public share provider service address",
			EnvVars:     []string{"REVA_STORAGE_PUBLICLINK_PUBLIC_SHARE_PROVIDER_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.PublicShareProviderAddr,
		},
		&cli.StringFlag{
			Name:        "user_provider_addr",
			Value:       "localhost:9144",
			Usage:       "user provider service address",
			EnvVars:     []string{"REVA_STORAGE_PUBLICLINK_USER_PROVIDER_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.UserProviderAddr,
		},
		&cli.StringFlag{
			Name:        "storage_provider_addr",
			Value:       "localhost:9154",
			Usage:       "storage provider service address",
			EnvVars:     []string{"REVA_STORAGE_PUBLICLINK_STORAGE_PROVIDER_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.StorageProviderAddr,
		},
		&cli.StringFlag{
			Name:        "driver",
			Value:       "owncloud",
			Usage:       "storage driver, eg. local, eos, owncloud or s3",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_DRIVER"},
			Destination: &cfg.Reva.StoragePublicLink.Driver,
		},
		&cli.StringFlag{
			Name:        "mount-id",
			Value:       "e1a73ede-549b-4226-abdf-40e69ca8230d",
			Usage:       "mount id",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_MOUNT_ID"},
			Destination: &cfg.Reva.StoragePublicLink.MountID,
		},
		&cli.BoolFlag{
			Name:        "expose-data-server",
			Value:       true,
			Usage:       "exposes a dedicated data server",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StoragePublicLink.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "http://localhost:9156/data",
			Usage:       "data server url",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StoragePublicLink.DataServerURL,
		},
		&cli.BoolFlag{
			Name:        "enable-home-creation",
			Value:       true,
			Usage:       "if enabled home dirs will be automatically created",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_ENABLE_HOME_CREATION"},
			Destination: &cfg.Reva.StoragePublicLink.EnableHomeCreation,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Value:       "/public/",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_MOUNT_PATH"},
			Destination: &cfg.Reva.StoragePublicLink.MountPath,
		},
		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       "localhost:9142",
			Usage:       "URL to use for the reva gateway service",
			EnvVars:     []string{"REVA_GATEWAY_URL"},
			Destination: &cfg.Reva.Gateway.URL,
		},
	}
}
