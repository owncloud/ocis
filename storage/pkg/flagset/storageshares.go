package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// StorageShares applies cfg to the root flagset
func StorageShares(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageShares.DebugAddr, "0.0.0.0:9179"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_SHARES_DEBUG_ADDR"},
			Destination: &cfg.Reva.StorageShares.DebugAddr,
		},

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageShares.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_SHARES_GRPC_NETWORK"},
			Destination: &cfg.Reva.StorageShares.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageShares.GRPCAddr, "0.0.0.0:9182"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_SHARES_GRPC_ADDR"},
			Destination: &cfg.Reva.StorageShares.GRPCAddr,
		},

		&cli.StringFlag{
			Name:        "mount-path",
			Value:       flags.OverrideDefaultString(cfg.Reva.StorageShares.MountPath, "/home/Shares"),
			Usage:       "mount path",
			EnvVars:     []string{"STORAGE_SHARES_MOUNT_PATH"},
			Destination: &cfg.Reva.StorageShares.MountPath,
		},

		&cli.StringFlag{
			Name:        "gateway-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "localhost:9142"),
			Usage:       "endpoint to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},

		&cli.StringFlag{
			Name:        "user-driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserDriver, "json"),
			Usage:       "driver to use for the UserShareProvider",
			EnvVars:     []string{"STORAGE_SHARING_USER_DRIVER"},
			Destination: &cfg.Reva.Sharing.UserDriver,
		},
		&cli.StringFlag{
			Name:        "user-json-file",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserJSONFile, "/var/tmp/ocis/storage/shares.json"),
			Usage:       "file used to persist shares for the UserShareProvider",
			EnvVars:     []string{"STORAGE_SHARING_USER_JSON_FILE"},
			Destination: &cfg.Reva.Sharing.UserJSONFile,
		},
		&cli.StringFlag{
			Name:        "user-sql-username",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserSQLUsername, ""),
			Usage:       "Username to be used to connect to the SQL database",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_USERNAME"},
			Destination: &cfg.Reva.Sharing.UserSQLUsername,
		},
		&cli.StringFlag{
			Name:        "user-sql-password",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserSQLPassword, ""),
			Usage:       "Password to be used to connect to the SQL database",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_PASSWORD"},
			Destination: &cfg.Reva.Sharing.UserSQLPassword,
		},
		&cli.StringFlag{
			Name:        "user-sql-host",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserSQLHost, ""),
			Usage:       "Hostname of the SQL database",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_HOST"},
			Destination: &cfg.Reva.Sharing.UserSQLHost,
		},
		&cli.IntFlag{
			Name:        "user-sql-port",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Sharing.UserSQLPort, 1433),
			Usage:       "The port on which the SQL database is exposed",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_PORT"},
			Destination: &cfg.Reva.Sharing.UserSQLPort,
		},
		&cli.StringFlag{
			Name:        "user-sql-name",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserSQLName, ""),
			Usage:       "Name of the SQL database",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_NAME"},
			Destination: &cfg.Reva.Sharing.UserSQLName,
		},
		&cli.IntFlag{
			Name:        "public-password-hash-cost",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Sharing.PublicPasswordHashCost, 11),
			Usage:       "the cost of hashing the public shares passwords",
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_PASSWORD_HASH_COST"},
			Destination: &cfg.Reva.Sharing.PublicPasswordHashCost,
		},
		&cli.BoolFlag{
			Name:        "public-enable-expired-shares-cleanup",
			Value:       flags.OverrideDefaultBool(cfg.Reva.Sharing.PublicEnableExpiredSharesCleanup, true),
			Usage:       "whether to periodically delete expired public shares",
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_ENABLE_EXPIRED_SHARES_CLEANUP"},
			Destination: &cfg.Reva.Sharing.PublicEnableExpiredSharesCleanup,
		},
		&cli.IntFlag{
			Name:        "public-janitor-run-interval",
			Value:       flags.OverrideDefaultInt(cfg.Reva.Sharing.PublicJanitorRunInterval, 60),
			Usage:       "the time period in seconds after which to start a janitor run",
			EnvVars:     []string{"STORAGE_SHARING_PUBLIC_JANITOR_RUN_INTERVAL"},
			Destination: &cfg.Reva.Sharing.PublicJanitorRunInterval,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
