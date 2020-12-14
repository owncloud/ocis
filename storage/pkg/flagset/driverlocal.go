package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// DriverLocalWithConfig applies cfg to the root flagset
func DriverLocalWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-local-root",
			Value:       "/var/tmp/ocis/storage/local",
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"STORAGE_DRIVER_LOCAL_ROOT"},
			Destination: &cfg.Reva.Storages.Local.Root,
		},
	}
}
