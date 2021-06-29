package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// SharingSQLWithConfig applies the Shring SQL driver cfg to the flagset
func SharingSQLWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
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
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode",
		},
	}
}
