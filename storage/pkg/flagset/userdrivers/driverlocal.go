package userdrivers

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
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.Local.Root, path.Join(defaults.BaseDataPath(), "storage", "local", "users")),
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_LOCAL_ROOT"},
			Destination: &cfg.Reva.UserStorage.Local.Root,
		},
		&cli.StringFlag{
			Name:        "storage-local-share-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.Local.ShareFolder, "/Shares"),
			Usage:       "the path to the local share folder",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_LOCAL_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.Local.ShareFolder,
		},
		&cli.StringFlag{
			Name:        "storage-local-user-layout",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.Local.UserLayout, "{{.Username}}"),
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.UsernameLower}} and {{.Provider}} also supports prefixing dirs: "{{.UsernamePrefixCount.2}}/{{.UsernameLower}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_LOCAL_USER_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.Local.UserLayout,
		},
	}
}
