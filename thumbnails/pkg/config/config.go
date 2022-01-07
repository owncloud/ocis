package config

import (
	"context"

	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	Service Service

	Tracing Tracing `ocisConfig:"tracing"`
	Log     *Log    `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	GRPC GRPC `ocisConfig:"grpc"`

	Thumbnail Thumbnail `ocisConfig:"thumbnail"`

	Context context.Context
}

// FileSystemStorage defines the available filesystem storage configuration.
type FileSystemStorage struct {
	RootDirectory string `ocisConfig:"root_directory" env:"THUMBNAILS_FILESYSTEMSTORAGE_ROOT"`
}

// FileSystemSource defines the available filesystem source configuration.
type FileSystemSource struct {
	BasePath string `ocisConfig:"base_path"`
}

// Thumbnail defines the available thumbnail related configuration.
type Thumbnail struct {
	Resolutions         []string          `ocisConfig:"resolutions"`
	FileSystemStorage   FileSystemStorage `ocisConfig:"filesystem_storage"`
	WebdavAllowInsecure bool              `ocisConfig:"webdav_allow_insecure" env:"OCIS_INSECURE;THUMBNAILS_WEBDAVSOURCE_INSECURE"`
	CS3AllowInsecure    bool              `ocisConfig:"cs3_allow_insecure" env:"OCIS_INSECURE;THUMBNAILS_CS3SOURCE_INSECURE"`
	RevaGateway         string            `ocisConfig:"reva_gateway" env:"REVA_GATEWAY"`
	WebdavNamespace     string            `ocisConfig:"webdav_namespace" env:"STORAGE_WEBDAV_NAMESPACE"`
	FontMapFile         string            `ocisConfig:"font_map_file" env:"THUMBNAILS_TXT_FONTMAP_FILE"`
}
