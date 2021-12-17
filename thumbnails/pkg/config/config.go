package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr" env:"THUMBNAILS_DEBUG_ADDR"`
	Token  string `ocisConfig:"token" env:"THUMBNAILS_DEBUG_TOKEN"`
	Pprof  bool   `ocisConfig:"pprof" env:"THUMBNAILS_DEBUG_PPROF"`
	Zpages bool   `ocisConfig:"zpages" env:"THUMBNAILS_DEBUG_ZPAGES"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `ocisConfig:"addr" env:"THUMBNAILS_GRPC_ADDR"`
	Namespace string
}

// Service defines the available service configuration.
type Service struct {
	Name    string
	Version string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled" env:"OCIS_TRACING_ENABLED;THUMBNAILS_TRACING_ENABLED"`
	Type      string `ocisConfig:"type" env:"OCIS_TRACING_TYPE;THUMBNAILS_TRACING_TYPE"`
	Endpoint  string `ocisConfig:"endpoint" env:"OCIS_TRACING_ENDPOINT;THUMBNAILS_TRACING_ENDPOINT"`
	Collector string `ocisConfig:"collector" env:"OCIS_TRACING_COLLECTOR;THUMBNAILS_TRACING_COLLECTOR"`
	Service   string `ocisConfig:"service" env:"THUMBNAILS_TRACING_SERVICE"` //TODO: should this be an ID? or the same as Service.Name?
}

// Log defines the available log configuration.
type Log struct {
	Level  string `mapstructure:"level" env:"OCIS_LOG_LEVEL;THUMBNAILS_LOG_LEVEL"`
	Pretty bool   `mapstructure:"pretty" env:"OCIS_LOG_PRETTY;THUMBNAILS_LOG_PRETTY"`
	Color  bool   `mapstructure:"color" env:"OCIS_LOG_COLOR;THUMBNAILS_LOG_COLOR"`
	File   string `mapstructure:"file" env:"OCIS_LOG_FILE;THUMBNAILS_LOG_FILE"`
}

// Config combines all available configuration parts.
type Config struct {
	Service Service `ocisConfig:"service"`

	Tracing Tracing `ocisConfig:"tracing"`
	Log     Log     `ocisConfig:"log"`
	Debug   Debug   `ocisConfig:"debug"`

	GRPC GRPC `ocisConfig:"grpc"`

	Thumbnail Thumbnail `ocisConfig:"thumbnail"`

	Context    context.Context
	Supervised bool
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
	Resolutions         []string          `ocisConfig:"resolutions"` // TODO: how to configure
	FileSystemStorage   FileSystemStorage `ocisConfig:"filesystem_storage"`
	WebdavAllowInsecure bool              `ocisConfig:"webdav_allow_insecure" env:"OCIS_INSECURE;THUMBNAILS_WEBDAVSOURCE_INSECURE"`
	CS3AllowInsecure    bool              `ocisConfig:"cs3_allow_insecure" env:"OCIS_INSECURE;THUMBNAILS_CS3SOURCE_INSECURE"`
	RevaGateway         string            `ocisConfig:"reva_gateway" env:"REVA_GATEWAY"`
	FontMapFile         string            `ocisConfig:"font_map_file" env:"THUMBNAILS_TXT_FONTMAP_FILE"`
}

func DefaultConfig() *Config {
	return &Config{
		Debug: Debug{
			Addr:   "127.0.0.1:9189",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: GRPC{
			Addr:      "127.0.0.1:9185",
			Namespace: "com.owncloud.api",
		},
		Service: Service{
			Name: "thumbnails",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "thumbnails",
		},
		Thumbnail: Thumbnail{
			Resolutions: []string{"16x16", "32x32", "64x64", "128x128", "1920x1080", "3840x2160", "7680x4320"},
			FileSystemStorage: FileSystemStorage{
				RootDirectory: path.Join(defaults.BaseDataPath(), "thumbnails"),
			},
			WebdavAllowInsecure: true,
			RevaGateway:         "127.0.0.1:9142",
			CS3AllowInsecure:    false,
		},
	}
}
