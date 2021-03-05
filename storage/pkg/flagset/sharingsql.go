package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// SharingSQLWithConfig applies the Shring SQL driver cfg to the flagset
func SharingSQLWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "user-sql-username",
			Value:       "",
			Usage:       "Username to be used to connect to the SQL database",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_USERNAME"},
			Destination: &cfg.Reva.Sharing.UserSQLUsername,
		},
		&cli.StringFlag{
			Name:        "user-sql-password",
			Value:       "",
			Usage:       "Password to be used to connect to the SQL database",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_PASSWORD"},
			Destination: &cfg.Reva.Sharing.UserSQLPassword,
		},
		&cli.StringFlag{
			Name:        "user-sql-host",
			Value:       "",
			Usage:       "Hostname of the SQL database",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_HOST"},
			Destination: &cfg.Reva.Sharing.UserSQLHost,
		},
		&cli.IntFlag{
			Name:        "user-sql-port",
			Value:       1433,
			Usage:       "The port on which the SQL database is exposed",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_PORT"},
			Destination: &cfg.Reva.Sharing.UserSQLPort,
		},
		&cli.StringFlag{
			Name:        "user-sql-name",
			Value:       "",
			Usage:       "Name of the SQL database",
			EnvVars:     []string{"STORAGE_SHARING_USER_SQL_NAME"},
			Destination: &cfg.Reva.Sharing.UserSQLName,
		},
	}
}
