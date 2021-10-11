package userdrivers

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverS3NGWithConfig applies cfg to the root flagset
func DriverS3WithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-s3-region",
			Value:       flags.OverrideDefaultString(cfg.Reva.UserStorage.S3.Region, "default"),
			Usage:       `"the s3 region" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_REGION"},
			Destination: &cfg.Reva.UserStorage.S3.Region,
		},
		&cli.StringFlag{
			Name:        "storage-s3-accesskey",
			Value:       "",
			Usage:       `"the s3 access key" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_ACCESS_KEY"},
			Destination: &cfg.Reva.UserStorage.S3.AccessKey,
		},
		&cli.StringFlag{
			Name:        "storage-s3-secretkey",
			Value:       "",
			Usage:       `"the secret s3 api key" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_SECRET_KEY"},
			Destination: &cfg.Reva.UserStorage.S3.SecretKey,
		},
		&cli.StringFlag{
			Name:        "storage-s3-endpoint",
			Value:       "",
			Usage:       `"s3 compatible API endpoint" `,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_ENDPOINT"},
			Destination: &cfg.Reva.UserStorage.S3.Endpoint,
		},
		&cli.StringFlag{
			Name:        "storage-s3-bucket",
			Value:       "",
			Usage:       `"bucket where the data will be stored in`,
			EnvVars:     []string{"STORAGE_USERS_DRIVER_S3_BUCKET"},
			Destination: &cfg.Reva.UserStorage.S3.Bucket,
		},
	}
}
