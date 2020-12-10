package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// DriverOCISWithConfig applies cfg to the root flagset
func DriverOCISWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-ocis-root",
			Value:       "./data/storage/ocis/users",
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"STORAGE_DRIVER_OCIS_ROOT"},
			Destination: &cfg.Reva.Storages.Common.Root,
		},
		&cli.BoolFlag{
			Name:        "storage-ocis-enable-home",
			Value:       false,
			Usage:       "enable the creation of home storages",
			EnvVars:     []string{"STORAGE_DRIVER_OCIS_ENABLE_HOME"},
			Destination: &cfg.Reva.Storages.Common.EnableHome,
		},
		&cli.StringFlag{
			Name:        "storage-ocis-layout",
			Value:       "{{.Id.OpaqueId}}",
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_DRIVER_OCIS_LAYOUT"},
			Destination: &cfg.Reva.Storages.Common.UserLayout,
		},
	}
}
