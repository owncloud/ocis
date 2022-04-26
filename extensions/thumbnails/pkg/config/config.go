package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"-"`

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing,omitempty"`
	Log     *Log     `yaml:"log,omitempty"`
	Debug   Debug    `yaml:"debug,omitempty"`

	GRPC GRPC `yaml:"grpc,omitempty"`
	HTTP HTTP `yaml:"http,omitempty"`

	Thumbnail Thumbnail `yaml:"thumbnail,omitempty"`

	Context context.Context `yaml:"-"`
}

// FileSystemStorage defines the available filesystem storage configuration.
type FileSystemStorage struct {
	RootDirectory string `yaml:"root_directory" env:"THUMBNAILS_FILESYSTEMSTORAGE_ROOT"`
}

// FileSystemSource defines the available filesystem source configuration.
type FileSystemSource struct {
	BasePath string `yaml:"base_path"`
}

// Thumbnail defines the available thumbnail related configuration.
type Thumbnail struct {
	Resolutions         []string          `yaml:"resolutions,omitempty"`
	FileSystemStorage   FileSystemStorage `yaml:"filesystem_storage,omitempty"`
	WebdavAllowInsecure bool              `yaml:"webdav_allow_insecure,omitempty" env:"OCIS_INSECURE;THUMBNAILS_WEBDAVSOURCE_INSECURE"`
	CS3AllowInsecure    bool              `yaml:"cs3_allow_insecure,omitempty" env:"OCIS_INSECURE;THUMBNAILS_CS3SOURCE_INSECURE"`
	RevaGateway         string            `yaml:"reva_gateway,omitempty" env:"REVA_GATEWAY"` //TODO: use REVA config
	FontMapFile         string            `yaml:"font_map_file,omitempty" env:"THUMBNAILS_TXT_FONTMAP_FILE"`
	TransferSecret      string            `yaml:"transfer_secret,omitempty" env:"THUMBNAILS_TRANSFER_TOKEN"`
	DataEndpoint        string            `yaml:"data_endpoint,omitempty" env:"THUMBNAILS_DATA_ENDPOINT"`
}
