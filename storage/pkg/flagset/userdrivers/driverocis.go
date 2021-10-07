package userdrivers

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverOCISWithConfig applies cfg to the root flagset
func DriverOCISWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-ocis-root",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.OCIS.Root, path.Join(defaults.BaseDataPath(), "storage", "users")),
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_ROOT"},
			Destination: &cfg.Reva.UserStorage.OCIS.Root,
		},
		&cli.StringFlag{
			Name:        "storage-ocis-layout",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.OCIS.UserLayout, "{{.Id.OpaqueId}}"),
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.OCIS.UserLayout,
		},
		&cli.StringFlag{
			Name:        "storage-ocis-share-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.OCIS.ShareFolder, "/Shares"),
			Usage:       "name of the shares folder",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.OCIS.ShareFolder,
		},
		&cli.StringFlag{
			Name:        "service-user-uuid",
			Value:       "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
			Usage:       "uuid of the internal service user",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_OCIS_SERVICE_USER_UUID"},
			Destination: &cfg.Reva.UserStorage.OCIS.ServiceUserUUID,
		},
	}
}
