package userdrivers

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverS3NGWithConfig applies cfg to the root flagset
func DriverS3NGWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-s3ng-root",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.S3NG.Root, path.Join(defaults.BaseDataPath(), "storage", "users")),
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_ROOT"},
			Destination: &cfg.Reva.UserStorage.S3NG.Root,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-layout",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.S3NG.UserLayout, "{{.Id.OpaqueId}}"),
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_LAYOUT"},
			Destination: &cfg.Reva.UserStorage.S3NG.UserLayout,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-share-folder",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.S3NG.ShareFolder, "/Shares"),
			Usage:       "name of the shares folder",
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_SHARE_FOLDER"},
			Destination: &cfg.Reva.UserStorage.S3NG.ShareFolder,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-region",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.S3NG.Region, "default"),
			Usage:       `"the s3 region" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_REGION"},
			Destination: &cfg.Reva.UserStorage.S3NG.Region,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-accesskey",
			Value:       "",
			Usage:       `"the s3 access key" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_ACCESS_KEY"},
			Destination: &cfg.Reva.UserStorage.S3NG.AccessKey,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-secretkey",
			Value:       "",
			Usage:       `"the secret s3 api key" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_SECRET_KEY"},
			Destination: &cfg.Reva.UserStorage.S3NG.SecretKey,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-endpoint",
			Value:       "",
			Usage:       `"s3 compatible API endpoint" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_ENDPOINT"},
			Destination: &cfg.Reva.UserStorage.S3NG.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-bucket",
			Value:       "",
			Usage:       `"bucket where the data will be stored in`,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3NG_BUCKET"},
			Destination: &cfg.Reva.UserStorage.S3NG.Bucket,
		},
	}
}
