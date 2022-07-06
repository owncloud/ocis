package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	GRPC GRPC `yaml:"grpc"`
	HTTP HTTP `yaml:"http"`

	Thumbnail Thumbnail `yaml:"thumbnail"`

	Context context.Context `yaml:"-"`
}

// FileSystemStorage defines the available filesystem storage configuration.
type FileSystemStorage struct {
	RootDirectory string `yaml:"root_directory" env:"THUMBNAILS_FILESYSTEMSTORAGE_ROOT" desc:"The directory where the filesystem storage will store the thumbnails."`
}

// Thumbnail defines the available thumbnail related configuration.
type Thumbnail struct {
	Resolutions         []string          `yaml:"resolutions" env:"THUMBNAILS_RESOLUTIONS" desc:"The supported target resolutions in the format WidthxHeight e.g. 32x32. You can define any resolution as required and separate multiple resolutions by blank or comma."`
	FileSystemStorage   FileSystemStorage `yaml:"filesystem_storage"`
	WebdavAllowInsecure bool              `yaml:"webdav_allow_insecure" env:"OCIS_INSECURE;THUMBNAILS_WEBDAVSOURCE_INSECURE" desc:"Ignore untrusted SSL certificates when connecting to the webdav source."`
	CS3AllowInsecure    bool              `yaml:"cs3_allow_insecure" env:"OCIS_INSECURE;THUMBNAILS_CS3SOURCE_INSECURE" desc:"Ignore untrusted SSL certificates when connecting to the CS3 source."`
	RevaGateway         string            `yaml:"reva_gateway" env:"REVA_GATEWAY" desc:"The CS3 gateway endpoint."` //TODO: use REVA config
	FontMapFile         string            `yaml:"font_map_file" env:"THUMBNAILS_TXT_FONTMAP_FILE" desc:"The path to a font file for txt thumbnails."`
	TransferSecret      string            `yaml:"transfer_secret" env:"THUMBNAILS_TRANSFER_TOKEN" desc:"The secret to sign JWT to download the actual thumbnail file."`
	DataEndpoint        string            `yaml:"data_endpoint" env:"THUMBNAILS_DATA_ENDPOINT" desc:"The HTTP endpoint where the actual thumbnail file can be downloaded."`
}
