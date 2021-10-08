package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset/userdrivers"
	"github.com/urfave/cli/v2"
)

// StorageUsersWithConfig applies cfg to the root flagset
func StorageUsersWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.DebugAddr, "0.0.0.0:9159"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_USERS_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageUsers.DebugAddr,
		},

		// Services

		// Storage home

		&cli.StringFlag{
			Name:        "grpc-network",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the users storage, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_USERS_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageUsers.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.GRPCAddr, "0.0.0.0:9157"),
			Usage:       "GRPC Address to bind users storage",
			EnvVars:     []string{"STORAGE_USERS_GRPC_ADDR"},
			Destination: &cfg.Reva.StorageUsers.GRPCAddr,
		},
		&cli.StringFlag{
			Name:        "http-network",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.HTTPNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_USERS_HTTP_NETWORK"},
			Destination: &cfg.Reva.StorageUsers.HTTPNetwork,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.HTTPAddr, "0.0.0.0:9158"),
			Usage:       "HTTP Address to bind users storage",
			EnvVars:     []string{"STORAGE_USERS_HTTP_ADDR"},
			Destination: &cfg.Reva.StorageUsers.HTTPAddr,
		},
		// TODO allow disabling grpc / http services
		/*
			&cli.StringSliceFlag{
				Name:    "grpc-service",
				Value:   cli.NewStringSlice("storageprovider"),
				Usage:   "--service storageprovider [--service otherservice]",
				EnvVars: []string{"STORAGE_USERS_GRPC_SERVICES"},
			},
			&cli.StringSliceFlag{
				Name:    "http-service",
				Value:   cli.NewStringSlice("dataprovider"),
				Usage:   "--service dataprovider [--service otherservice]",
				EnvVars: []string{"STORAGE_USERS_HTTP_SERVICES"},
			},
		*/

		&cli.StringFlag{
			Name:        "driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.Driver, "ocis"),
			Usage:       "storage driver for users mount: eg. local, eos, owncloud, ocis or s3",
			EnvVars:     []string{"STORAGE_USERS_DRIVER"},
			Destination: &cfg.Reva.StorageUsers.Driver,
		},
		&cli.BoolFlag{
			Name:        "read-only",
			Value:       flags.OverrideDefaultBool(cfg.Reva.StorageUsers.ReadOnly, false),
			Usage:       "use storage driver in read-only mode",
			EnvVars:     []string{"STORAGE_USERS_READ_ONLY", "OCIS_STORAGE_READ_ONLY"},
			Destination: &cfg.Reva.StorageUsers.ReadOnly,
		},
		&cli.StringFlag{
			Name:        "mount-path",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.MountPath, "/users"),
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_USERS_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageUsers.MountPath,
		},
		&cli.StringFlag{
			Name:        "mount-id",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.MountID, "1284d238-aa92-42ce-bdc4-0b0000009157"),
			Usage:       "mount id",
			EnvVars:     []string{"STORAGE_USERS_MOUNT_ID"},
			Destination: &cfg.Reva.StorageUsers.MountID,
		},
		&cli.BoolFlag{
			Name:        "expose-data-server",
			Value:       flags.OverrideDefaultBool(cfg.Reva.StorageUsers.ExposeDataServer, false),
			Usage:       "exposes a dedicated data server",
			EnvVars:     []string{"STORAGE_USERS_EXPOSE_DATA_SERVER"},
			Destination: &cfg.Reva.StorageUsers.ExposeDataServer,
		},
		&cli.StringFlag{
			Name:        "data-server-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.DataServerURL, "http://localhost:9158/data"),
			Usage:       "data server url",
			EnvVars:     []string{"STORAGE_USERS_DATA_SERVER_URL"},
			Destination: &cfg.Reva.StorageUsers.DataServerURL,
		},
		&cli.StringFlag{
			Name:        "http-prefix",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.HTTPPrefix, "data"),
			Usage:       "prefix for the http endpoint, without leading slash",
			EnvVars:     []string{"STORAGE_USERS_HTTP_PREFIX"},
			Destination: &cfg.Reva.StorageUsers.HTTPPrefix,
		},
		&cli.StringFlag{
			Name:        "tmp-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageUsers.TempFolder, "/var/tmp/ocis/tmp/users"),
			Usage:       "path to tmp folder",
			EnvVars:     []string{"STORAGE_USERS_TMP_FOLDER"},
			Destination: &cfg.Reva.StorageUsers.TempFolder,
		},

		// some drivers need to look up users at the gateway

		// Gateway

		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY"},
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
	flags = append(flags, userdrivers.DriverEOSWithConfig(cfg)...)
	flags = append(flags, userdrivers.DriverLocalWithConfig(cfg)...)
	flags = append(flags, userdrivers.DriverOwnCloudWithConfig(cfg)...)
	flags = append(flags, userdrivers.DriverOwnCloudSQLWithConfig(cfg)...)
	flags = append(flags, userdrivers.DriverOCISWithConfig(cfg)...)
	flags = append(flags, userdrivers.DriverS3NGWithConfig(cfg)...)

	return flags
}
