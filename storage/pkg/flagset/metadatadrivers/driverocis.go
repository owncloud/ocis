package metadatadrivers

import (
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/urfave/cli/v2"
)

// DriverOCISWithConfig applies cfg to the root flagset
func DriverOCISWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-ocis-root",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.OCIS.Root, "/var/tmp/ocis/storage/metadata"),
			Usage:       "the path to the local storage root",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OCIS_ROOT"},
			Destination: &cfg.Reva.MetadataStorage.OCIS.Root,
		},
		&cli.StringFlag{
			Name:        "storage-ocis-layout",
			Value:       flags.OverrideDefaultString(cfg.Reva.MetadataStorage.OCIS.UserLayout, "{{.Id.OpaqueId}}"),
			Usage:       `"layout of the users home dir path on disk, in addition to {{.Username}}, {{.Mail}}, {{.Id.OpaqueId}}, {{.Id.Idp}} also supports prefixing dirs: "{{substr 0 1 .Username}}/{{.Username}}" will turn "Einstein" into "Ei/Einstein" `,
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OCIS_LAYOUT"},
			Destination: &cfg.Reva.MetadataStorage.OCIS.UserLayout,
		},
		&cli.StringFlag{
			Name:        "service-user-uuid",
			Value:       "95cb8724-03b2-11eb-a0a6-c33ef8ef53ad",
			Usage:       "uuid of the internal service user",
			EnvVars:     []string{"STORAGE_METADATA_DRIVER_OCIS_SERVICE_USER_UUID"},
			Destination: &cfg.Reva.MetadataStorage.OCIS.ServiceUserUUID,
		},
	}
}
