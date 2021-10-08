package userdrivers

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverLocalWithConfig applies cfg to the root flagset
func DriverLocalWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-local-root",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.Local.Root, "/var/tmp/ocis/storage/local"),
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_LOCAL_ROOT"},
			Destination: &cfg.Reva.UserStorage.Local.Root,
		},
	}
}
