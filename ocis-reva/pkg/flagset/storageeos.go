package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-reva/pkg/config"
)

// StorageEOSWithConfig applies cfg to the root flagset
func StorageEOSWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9159",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_STORAGE_EOS_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageEOS.DebugAddr,
		},

		// Storage eos

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_STORAGE_EOS_NETWORK"},
			Destination: &cfg.Reva.StorageEOS.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_STORAGE_EOS_PROTOCOL"},
			Destination: &cfg.Reva.StorageEOS.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9158",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_STORAGE_EOS_ADDR"},
			Destination: &cfg.Reva.StorageEOS.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9158",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_EOS_URL"},
			Destination: &cfg.Reva.StorageEOS.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("storageprovider"),
			Usage:   "--service storageprovider [--service otherservice]",
			EnvVars: []string{"REVA_STORAGE_EOS_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       "eos",
			Usage:       "storage driver for eos mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"REVA_STORAGE_EOS_DRIVER"},
			Destination: &cfg.Reva.StorageEOS.Driver,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Value:       "/eos",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_EOS_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageEOS.MountPath,
		},
		&cli.StringFlag{
			Name:        "mount-id",
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009158",
			Usage:       "mount id",
			EnvVars:     []string{"REVA_STORAGE_EOS_MOUNT_ID"},
			Destination: &cfg.Reva.StorageEOS.MountID,
		},
		&cli.BoolFlag{
			Name:        "expose-data-server",
			Value:       false,
			Usage:       "exposes a dedicated data server",
			EnvVars:     []string{"REVA_STORAGE_EOS_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StorageEOS.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "http://localhost:9160/data",
			Usage:       "data server url",
			EnvVars:     []string{"REVA_STORAGE_EOS_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageEOS.DataServerURL,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, DriverEOSWithConfig(cfg)...)
	flags = append(flags, DriverLocalWithConfig(cfg)...)
	flags = append(flags, DriverOwnCloudWithConfig(cfg)...)
	flags = append(flags, DriverOCISWithConfig(cfg)...)

	return flags
}
