package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// StorageHomeWithConfig applies cfg to the root flagset
func StorageHomeWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.DebugAddr, "0.0.0.0:9156"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_HOME_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageHome.DebugAddr,
		},

		// Services

		// Storage home

		&cli.StringFlag{
			Name:        "grpc-network",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_HOME_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageHome.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.GRPCAddr, "0.0.0.0:9154"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_HOME_GRPC_ADDR"},
			Destination: &cfg.Reva.StorageHome.GRPCAddr,
		},
		&cli.StringFlag{
			Name:        "http-network",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.HTTPNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_HOME_HTTP_NETWORK"},
			Destination: &cfg.Reva.StorageHome.HTTPNetwork,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.HTTPAddr, "0.0.0.0:9155"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_HOME_HTTP_ADDR"},
			Destination: &cfg.Reva.StorageHome.HTTPAddr,
		},

		// TODO allow disabling grpc / http services
		/*
			&cli.StringSliceFlag{
				Name:    "grpc-service",
				Value:   cli.NewStringSlice("storageprovider"),
				Usage:   "--service storageprovider [--service otherservice]",
				EnvVars: []string{"STORAGE_HOME_GRPC_SERVICES"},
			},
			&cli.StringSliceFlag{
				Name:    "http-service",
				Value:   cli.NewStringSlice("dataprovider"),
				Usage:   "--service dataprovider [--service otherservice]",
				EnvVars: []string{"STORAGE_HOME_HTTP_SERVICES"},
			},
		*/

		&cli.StringFlag{
			Name:        "driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.Driver, "ocis"),
			Usage:       "storage driver for home mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"STORAGE_HOME_DRIVER"},
			Destination: &cfg.Reva.StorageHome.Driver,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.MountPath, "/home"),
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_HOME_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageHome.MountPath,
		},
		&cli.StringFlag{
			Name: "mount-id",
			// This is the mount id of the storage provider using the same storage driver
			// as /home but without home enabled.
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.MountID, "1284d238-aa92-42ce-bdc4-0b0000009157"),
			Usage:       "mount id",
			EnvVars:     []string{"STORAGE_HOME_MOUNT_ID"},
			Destination: &cfg.Reva.StorageHome.MountID,
		},
		&cli.BoolFlag{
			Name:        "expose-data-server",
			Value:       flags.OverrideDefaultBool(cfg.Reva.StorageHome.ExposeDataServer, false),
			Usage:       "exposes a dedicated data server",
			EnvVars:     []string{"STORAGE_HOME_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StorageHome.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.DataServerURL, "http://localhost:9155/data"),
			Usage:       "data server url",
			EnvVars:     []string{"STORAGE_HOME_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageHome.DataServerURL,
		},
		&cli.StringFlag{
			Name:        "http-prefix",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.HTTPPrefix, "data"),
			Usage:       "prefix for the http endpoint, without leading slash",
			EnvVars:     []string{"STORAGE_HOME_HTTP_PREFIX"},
			Destination: &cfg.Reva.StorageHome.HTTPPrefix,
		},
		&cli.StringFlag{
			Name:        "tmp-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageHome.TempFolder, "/var/tmp/ocis/tmp/home"),
			Usage:       "path to tmp folder",
			EnvVars:     []string{"STORAGE_HOME_TMP_FOLDER"},
			Destination: &cfg.Reva.StorageHome.TempFolder,
		},
		&cli.BoolFlag{
			Name:        "enable-home",
			Value:       flags.OverrideDefaultBool(cfg.Reva.Storages.Home.EnableHome, true),
			Usage:       "enable the creation of home directories",
			EnvVars:     []string{"STORAGE_HOME_ENABLE_HOME"},
			Destination: &cfg.Reva.Storages.Home.EnableHome,
		},

		// some drivers need to look up users at the gateway

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "localhost:9142"),
			Usage:       "endpoint to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},
		// User provider

		&cli.StringFlag{
			Name:        "users-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Users.Endpoint, "localhost:9144"),
			Usage:       "endpoint to use for the storage service",
			EnvVars:     []string{"STORAGE_USERPROVIDER_ENDPOINT"},
			Destination: &cfg.Reva.Users.Endpoint,
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
