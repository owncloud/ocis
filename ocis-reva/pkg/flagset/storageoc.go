package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-reva/pkg/config"
)

// StorageOCWithConfig applies cfg to the root flagset
func StorageOCWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9163",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_STORAGE_OC_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageOC.DebugAddr,
		},

		// Services

		// Storage oc

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_STORAGE_OC_NETWORK"},
			Destination: &cfg.Reva.StorageOC.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_STORAGE_OC_PROTOCOL"},
			Destination: &cfg.Reva.StorageOC.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9162",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_STORAGE_OC_ADDR"},
			Destination: &cfg.Reva.StorageOC.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9162",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_OC_URL"},
			Destination: &cfg.Reva.StorageOC.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("storageprovider"),
			Usage:   "--service storageprovider [--service otherservice]",
			EnvVars: []string{"REVA_STORAGE_OC_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       "owncloud",
			Usage:       "storage driver for oc mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"REVA_STORAGE_OC_DRIVER"},
			Destination: &cfg.Reva.StorageOC.Driver,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Value:       "/oc",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_OC_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageOC.MountPath,
		},
		&cli.StringFlag{
			Name:        "mount-id",
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009162",
			Usage:       "mount id",
			EnvVars:     []string{"REVA_STORAGE_OC_MOUNT_ID"},
			Destination: &cfg.Reva.StorageOC.MountID,
		},
		&cli.BoolFlag{
			Name:        "expose-data-server",
			Value:       false,
			Usage:       "exposes a dedicated data server",
			EnvVars:     []string{"REVA_STORAGE_OC_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StorageOC.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "http://localhost:9164/data",
			Usage:       "data server url",
			EnvVars:     []string{"REVA_STORAGE_OC_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageOC.DataServerURL,
		},

		// User provider

		&cli.StringFlag{
			Name:        "users-url",
			Value:       "localhost:9144",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_USERS_URL"},
			Destination: &cfg.Reva.Users.URL,
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
