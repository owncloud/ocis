package metadatadrivers

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverLocalWithConfig applies cfg to the root flagset
func DriverLocalWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-local-root",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.Local.Root, path.Join(defaults.BaseDataPath(), "storage", "local", "metadata")),
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_LOCAL_ROOT"},
			Destination: &cfg.Reva.MetadataStorage.Local.Root,
		},
	}
}
