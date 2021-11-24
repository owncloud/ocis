package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `ocisConfig:"addr"`
	Token  string `ocisConfig:"token"`
	Pprof  bool   `ocisConfig:"pprof"`
	Zpages bool   `ocisConfig:"zpages"`
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr      string `ocisConfig:"addr"`
	Namespace string `ocisConfig:"namespace"`
}

// Service provides configuration options for the service
type Service struct {
	Name    string `ocisConfig:"name"`
	Version string `ocisConfig:"version"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `ocisConfig:"enabled"`
	Type      string `ocisConfig:"type"`
	Endpoint  string `ocisConfig:"endpoint"`
	Collector string `ocisConfig:"collector"`
	Service   string `ocisConfig:"service"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons

	File      string      `ocisConfig:"file"`
	Log       *shared.Log `ocisConfig:"log"`
	Debug     Debug       `ocisConfig:"debug"`
	GRPC      GRPC        `ocisConfig:"grpc"`
	Service   Service     `ocisConfig:"service"`
	Tracing   Tracing     `ocisConfig:"tracing"`
	Thumbnail Thumbnail   `ocisConfig:"thumbnail"`

	Context    context.Context
	Supervised bool
}

// FileSystemStorage defines the available filesystem storage configuration.
type FileSystemStorage struct {
	RootDirectory string `ocisConfig:"root_directory"`
}

// FileSystemSource defines the available filesystem source configuration.
type FileSystemSource struct {
	BasePath string `ocisConfig:"base_path"`
}

// Thumbnail defines the available thumbnail related configuration.
type Thumbnail struct {
	Resolutions         []string          `ocisConfig:"resolutions"`
	FileSystemStorage   FileSystemStorage `ocisConfig:"filesystem_storage"`
	WebdavAllowInsecure bool              `ocisConfig:"webdav_allow_insecure"`
	CS3AllowInsecure    bool              `ocisConfig:"cs3_allow_insecure"`
	RevaGateway         string            `ocisConfig:"reva_gateway"`
	WebdavNamespace     string            `ocisConfig:"webdav_namespace"`
	FontMapFile         string            `ocisConfig:"font_map_file"`
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
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
			WebdavNamespace:     "/home",
			CS3AllowInsecure:    false,
		},
	}
}
