package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// StoragePublicLink applies cfg to the root flagset
func StoragePublicLink(cfg *config.Config) []cli.Flag {
	flags := commonTracingWithConfig(cfg)

	flags = append(flags, commonGatewayWithConfig(cfg)...)

	flags = append(flags, commonSecretWithConfig(cfg)...)

	flags = append(flags, commonDebugWithConfig(cfg)...)

	flags = append(flags, storageDriversWithConfig(cfg)...)

	flags = append(flags,
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
			Value:       "0.0.0.0:9170",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9171",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_DEBUG_ADDR"},
			Destination: &cfg.Reva.StoragePublicLink.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9170",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_URL"},
			Destination: &cfg.Reva.StoragePublicLink.URL,
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
			Name:        "mount-path",
			Value:       "/public/",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_MOUNT_PATH"},
			Destination: &cfg.Reva.StoragePublicLink.MountPath,
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
			Value:       false,
			Usage:       "exposes a dedicated data server",
			EnvVars:     []string{"REVA_STORAGE_PUBLIC_LINK_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StoragePublicLink.ExposeDataServer,
		},
		// has no data provider, only redirects to the actual storage
	)

	return flags
}
