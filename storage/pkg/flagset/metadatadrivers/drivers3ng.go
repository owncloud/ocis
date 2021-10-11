package metadatadrivers

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverS3NGWithConfig applies cfg to the root flagset
func DriverS3NGWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-s3ng-root",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.S3NG.Root, "/var/tmp/ocis/storage/metadata"),
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_ROOT"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.Root,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-layout",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.S3NG.UserLayout, "{{.Id.OpaqueId}}"),
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_LAYOUT"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.UserLayout,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-region",
			Value:       "default",
			Usage:       `"the s3 region" `,
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_REGION"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.Region,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-accesskey",
			Value:       "",
			Usage:       `"the s3 access key" `,
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_ACCESS_KEY"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.AccessKey,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-secretkey",
			Value:       "",
			Usage:       `"the secret s3 api key" `,
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_SECRET_KEY"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.SecretKey,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-endpoint",
			Value:       "",
			Usage:       `"s3 compatible API endpoint" `,
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_ENDPOINT"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage-s3ng-bucket",
			Value:       "",
			Usage:       `"bucket where the data will be stored in`,
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_S3NG_BUCKET"},
			Destination: &cfg.Reva.MetadataStorage.S3NG.Bucket,
		},
	}
}
