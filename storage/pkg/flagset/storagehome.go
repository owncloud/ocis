package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// StorageHomeWithConfig applies cfg to the root flagset
func StorageHomeWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9155",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageHome.DebugAddr,
		},

		// Services

		// Storage home

		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_NETWORK"},
			Destination: &cfg.Reva.StorageHome.Network,
		},
		&cli.StringFlag{
			Name:        "protocol",
			Value:       "grpc",
			Usage:       "protocol for storage service, can be 'http' or 'grpc'",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_PROTOCOL"},
			Destination: &cfg.Reva.StorageHome.Protocol,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9154",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_ADDR"},
			Destination: &cfg.Reva.StorageHome.Addr,
		},
		&cli.StringFlag{
			Name:        "url",
			Value:       "localhost:9154",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_URL"},
			Destination: &cfg.Reva.StorageHome.URL,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("storageprovider"),
			Usage:   "--service storageprovider [--service otherservice]",
			EnvVars: []string{"STORAGE_STORAGE_HOME_SERVICES"},
		},

		&cli.StringFlag{
			Name:        "driver",
			Value:       "owncloud",
			Usage:       "storage driver for home mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_DRIVER"},
			Destination: &cfg.Reva.StorageHome.Driver,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Value:       "/home",
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageHome.MountPath,
		},
		&cli.StringFlag{
			Name: "mount-id",
			// This is the mount id of the storage provider using the same storage driver
			// as /home but withoud home enabled. Set it to
			// 1284d238-aa92-42ce-bdc4-0b0000009158 for /eos
			// 1284d238-aa92-42ce-bdc4-0b0000009162 for /oc
			Value:       "1284d238-aa92-42ce-bdc4-0b0000009162", // /oc
			Usage:       "mount id",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_MOUNT_ID"},
			Destination: &cfg.Reva.StorageHome.MountID,
		},
		&cli.BoolFlag{
			Name:        "expose-data-server",
			Value:       false,
			Usage:       "exposes a dedicated data server",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StorageHome.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       "http://localhost:9156/data",
			Usage:       "data server url",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageHome.DataServerURL,
		},
		&cli.BoolFlag{
			Name:        "enable-home",
			Value:       true,
			Usage:       "enable the creation of home directories",
			EnvVars:     []string{"STORAGE_STORAGE_HOME_ENABLE_HOME"},
			Destination: &cfg.Reva.Storages.Home.EnableHome,
		},

		// User provider

		&cli.StringFlag{
			Name:        "users-url",
			Value:       "localhost:9144",
			Usage:       "URL to use for the storage service",
			EnvVars:     []string{"STORAGE_USERS_URL"},
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
