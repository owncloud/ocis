package flagset

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// SharingSQLWithConfig applies the Sharing SQL driver cfg to the flagset
func SharingSQLWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-mount-id",
			Value:       flags.OverrideDefaultString(cfg.Reva.Sharing.UserStorageMountId, "1284d238-aa92-42ce-bdc4-0b0000009157"),
			Usage:       "mount id of the storage that is used for accessing the shares",
			EnvVars:     []string{"STORAGE_SHARING_USER_STORAGE_MOUNT_ID"},
			Destination: &cfg.Reva.Sharing.UserStorageMountId,
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
	}
}
