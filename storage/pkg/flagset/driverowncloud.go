package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// DriverOwnCloudWithConfig applies cfg to the root flagset
func DriverOwnCloudWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-owncloud-datadir",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloud.Root, "/var/tmp/ocis/storage/owncloud"),
			Usage:       "the path to the owncloud data directory",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUD_DATADIR"},
			Destination: &cfg.Reva.Storages.OwnCloud.Root,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-uploadinfo-dir",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloud.UploadInfoDir, "/var/tmp/ocis/storage/uploadinfo"),
			Usage:       "the path to the tus upload info directory",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUD_UPLOADINFO_DIR"},
			Destination: &cfg.Reva.Storages.OwnCloud.UploadInfoDir,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-share-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloud.ShareFolder, "/Shares"),
			Usage:       "name of the shares folder",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUD_SHARE_FOLDER"},
			Destination: &cfg.Reva.Storages.OwnCloud.ShareFolder,
		},
		&cli.BoolFlag{
			Name:        "storage-owncloud-scan",
			Value:       flags.OverrideDefaultBool(cfg.Reva.Storages.OwnCloud.Scan, true),
			Usage:       "scan files on startup to add fileids",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUD_SCAN"},
			Destination: &cfg.Reva.Storages.OwnCloud.Scan,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-redis",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloud.Redis, ":6379"),
			Usage:       "the address of the redis server",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUD_REDIS_ADDR"},
			Destination: &cfg.Reva.Storages.OwnCloud.Redis,
		},
		&cli.BoolFlag{
			Name:        "storage-owncloud-enable-home",
			Value:       flags.OverrideDefaultBool(cfg.Reva.Storages.OwnCloud.EnableHome, false),
			Usage:       "enable the creation of home storages",
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUD_ENABLE_HOME"},
			Destination: &cfg.Reva.Storages.OwnCloud.EnableHome,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-layout",
			Value:       flags.OverrideDefaultString(cfg.Reva.Storages.OwnCloud.UserLayout, "{{.Id.OpaqueId}}"),
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_DRIVER_OWNCLOUD_LAYOUT"},
			Destination: &cfg.Reva.Storages.OwnCloud.UserLayout,
		},
	}
}
