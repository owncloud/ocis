package config

import (
	"context"
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `mapstructure:"addr"`
	Token  string `mapstructure:"token"`
	Pprof  bool   `mapstructure:"pprof"`
	Zpages bool   `mapstructure:"zpages"`
}

// Server defines the available server configuration.
type Server struct {
	Name      string `mapstructure:"name"`
	Namespace string `mapstructure:"namespace"`
	Address   string `mapstructure:"address"`
	Version   string `mapstructure:"version"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"enabled"`
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	Collector string `mapstructure:"collector"`
	Service   string `mapstructure:"service"`
}

// Config combines all available configuration parts.
type Config struct {
	File      string     `mapstructure:"file"`
	Log       shared.Log `mapstructure:"log"`
	Debug     Debug      `mapstructure:"debug"`
	Server    Server     `mapstructure:"server"`
	Tracing   Tracing    `mapstructure:"tracing"`
	Thumbnail Thumbnail  `mapstructure:"thumbnail"`

	Context    context.Context
	Supervised bool
}

// FileSystemStorage defines the available filesystem storage configuration.
type FileSystemStorage struct {
	RootDirectory string `mapstructure:"root_directory"`
}

// FileSystemSource defines the available filesystem source configuration.
type FileSystemSource struct {
	BasePath string `mapstructure:"base_path"`
}

// Thumbnail defines the available thumbnail related configuration.
type Thumbnail struct {
	Resolutions         []string          `mapstructure:"resolutions"`
	FileSystemStorage   FileSystemStorage `mapstructure:"filesystem_storage"`
	WebdavAllowInsecure bool              `mapstructure:"webdav_allow_insecure"`
	CS3AllowInsecure    bool `mapstructure:"cs3_allow_insecure"`
	RevaGateway         string            `mapstructure:"reva_gateway"`
	WebdavNamespace     string            `mapstructure:"webdav_namespace"`
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		Log: shared.Log{},
		Debug: Debug{
			Addr:   "127.0.0.1:9189",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		Server: Server{
			Name:      "thumbnails",
			Namespace: "com.owncloud.api",
			Address:   "127.0.0.1:9185",
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
			CS3AllowInsecure: false,
		},
	}
}

// GetEnv fetches a list of known env variables for this extension. It is to be used by gookit, as it provides a list
// with all the environment variables an extension supports.
func GetEnv() []string {
	var r = make([]string, len(structMappings(&Config{})))
	for i := range structMappings(&Config{}) {
		r = append(r, structMappings(&Config{})[i].EnvVars...)
	}

	return r
}
