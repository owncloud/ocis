package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// StorageRootWithConfig applies cfg to the root flagset
func StorageRootWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9153",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_STORAGE_ROOT_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageRoot.DebugAddr,
		},

		// Services

		// Storage root

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_STORAGE_ROOT_NETWORK"},
			Destination: &cfg.Reva.StorageRoot.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_STORAGE_ROOT_PROTOCOL"},
			Destination: &cfg.Reva.StorageRoot.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9152",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_STORAGE_ROOT_ADDR"},
			Destination: &cfg.Reva.StorageRoot.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9152",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_ROOT_URL"},
			Destination: &cfg.Reva.StorageRoot.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("storageprovider"),
			Usage:   "--service storageprovider [--service otherservice]",
			EnvVars: []string{"REVA_STORAGE_ROOT_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       "local",
			Usage:       "storage driver, eg. local, eos, owncloud or s3",
			EnvVars:     []string{"REVA_STORAGE_ROOT_DRIVER"},
			Destination: &cfg.Reva.StorageRoot.Driver,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Value:       "/",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_ROOT_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageRoot.MountPath,
		},
		&cli.StringFlag{
			Name:        "mount-id",
			Value:       "123e4567-e89b-12d3-a456-426655440001",
			Usage:       "mount id",
			EnvVars:     []string{"REVA_STORAGE_ROOT_MOUNT_ID"},
			Destination: &cfg.Reva.StorageRoot.MountID,
		},
		&cli.BoolFlag{
			Name:        "expose-data-server",
			Usage:       "exposes a dedicated data server",
			EnvVars:     []string{"REVA_STORAGE_ROOT_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StorageRoot.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "",
			Usage:       "data server url",
			EnvVars:     []string{"REVA_STORAGE_ROOT_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageRoot.DataServerURL,
		},
		&cli.BoolFlag{
			Name:        "enable-home-creation",
			Usage:       "if enabled home dirs will be automatically created",
			EnvVars:     []string{"REVA_STORAGE_HOME_ENABLE_HOME_CREATION"},
			Destination: &cfg.Reva.StorageHome.EnableHomeCreation,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, DriverEOSWithConfig(cfg)...)
	flags = append(flags, DriverLocalWithConfig(cfg)...)
	flags = append(flags, DriverOwnCloudWithConfig(cfg)...)

	return flags
}
