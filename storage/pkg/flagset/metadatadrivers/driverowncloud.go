package metadatadrivers

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverOwnCloudWithConfig applies cfg to the root flagset
func DriverOwnCloudWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-owncloud-datadir",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.OwnCloud.Root, "/var/tmp/ocis/storage/owncloud"),
			Usage:       "the path to the owncloud data directory",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OWNCLOUD_DATADIR"},
			Destination: &cfg.Reva.MetadataStorage.OwnCloud.Root,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-uploadinfo-dir",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.OwnCloud.UploadInfoDir, "/var/tmp/ocis/storage/uploadinfo"),
			Usage:       "the path to the tus upload info directory",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OWNCLOUD_UPLOADINFO_DIR"},
			Destination: &cfg.Reva.MetadataStorage.OwnCloud.UploadInfoDir,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-share-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.OwnCloud.ShareFolder, "/Shares"),
			Usage:       "name of the shares folder",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OWNCLOUD_SHARE_FOLDER"},
			Destination: &cfg.Reva.MetadataStorage.OwnCloud.ShareFolder,
		},
		&cli.BoolFlag{
			Name:        "storage-owncloud-scan",
			Value:       flags.OverrideDefaultBool(cfg.Reva.MetadataStorage.OwnCloud.Scan, true),
			Usage:       "scan files on startup to add fileids",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OWNCLOUD_SCAN"},
			Destination: &cfg.Reva.MetadataStorage.OwnCloud.Scan,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-redis",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.OwnCloud.Redis, ":6379"),
			Usage:       "the address of the redis server",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OWNCLOUD_REDIS_ADDR"},
			Destination: &cfg.Reva.MetadataStorage.OwnCloud.Redis,
		},
		&cli.StringFlag{
			Name:        "storage-owncloud-layout",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.OwnCloud.UserLayout, "{{.Id.OpaqueId}}"),
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OWNCLOUD_LAYOUT"},
			Destination: &cfg.Reva.MetadataStorage.OwnCloud.UserLayout,
		},
	}
}
