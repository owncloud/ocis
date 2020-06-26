package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis-reva/pkg/config"
)

// StorageHomeWithConfig applies cfg to the root flagset
func StorageHomeWithConfig(cfg *config.Config) []cli.Flag {
	flags := commonTracingWithConfig(cfg)

	flags = append(flags, commonSecretWithConfig(cfg)...)

	flags = append(flags, commonDebugWithConfig(cfg)...)

	flags = append(flags, storageDriversWithConfig(cfg)...)

	flags = append(flags,

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the reva service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"REVA_STORAGE_HOME_NETWORK"},
			Destination: &cfg.Reva.StorageHome.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for reva service, can be 'http' or 'grpc'",
			EnvVars:     []string{"REVA_STORAGE_HOME_PROTOCOL"},
			Destination: &cfg.Reva.StorageHome.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9154",
			Usage:       "Address to bind reva service",
			EnvVars:     []string{"REVA_STORAGE_HOME_ADDR"},
			Destination: &cfg.Reva.StorageHome.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9155",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"REVA_STORAGE_HOME_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageHome.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9154",
			Usage:       "URL to use for the reva service",
			EnvVars:     []string{"REVA_STORAGE_HOME_URL"},
			Destination: &cfg.Reva.StorageHome.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("storageprovider"),
			Usage:   "--service storageprovider [--service otherservice]",
			EnvVars: []string{"REVA_STORAGE_HOME_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       "owncloud",
			Usage:       "storage driver, eg. local, eos, owncloud or s3",
			EnvVars:     []string{"REVA_STORAGE_HOME_DRIVER"},
			Destination: &cfg.Reva.StorageHome.Driver,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Value:       "/home",
			Usage:       "mount path",
			EnvVars:     []string{"REVA_STORAGE_HOME_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageHome.MountPath,
		},
		&cli.StringFlag{
			Name: "mount-id",
			// This is tho mount id of the /oc storage
			// set it to 1284d238-aa92-42ce-bdc4-0b0000009158 for /eos
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009162",
			Usage:       "mount id",
			EnvVars:     []string{"REVA_STORAGE_HOME_MOUNT_ID"},
			Destination: &cfg.Reva.StorageHome.MountID,
		},
		&cli.BoolFlag{
			Name:        "expose-data-server",
			Value:       false,
			Usage:       "exposes a dedicated data server",
			EnvVars:     []string{"REVA_STORAGE_HOME_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StorageHome.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "http://localhost:9156/data",
			Usage:       "data server url",
			EnvVars:     []string{"REVA_STORAGE_HOME_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageHome.DataServerURL,
		},
	)

	return flags
}
