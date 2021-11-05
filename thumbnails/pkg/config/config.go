package config

import (
	"context"
	"fmt"
	"path"
	"reflect"

	gofig "github.com/gookit/config/v2"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
)

// Log defines the available logging configuration.
type Log struct {
	Level  string `mapstructure:"level"`
	Pretty bool   `mapstructure:"pretty"`
	Color  bool   `mapstructure:"color"`
	File   string `mapstructure:"file"`
}

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
	File      string    `mapstructure:"file"`
	Log       Log       `mapstructure:"log"`
	Debug     Debug     `mapstructure:"debug"`
	Server    Server    `mapstructure:"server"`
	Tracing   Tracing   `mapstructure:"tracing"`
	Thumbnail Thumbnail `mapstructure:"thumbnail"`

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
	RevaGateway         string            `mapstructure:"reva_gateway"`
	WebdavNamespace     string            `mapstructure:"webdav_namespace"`
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		Log: Log{},
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

// UnmapEnv loads values from the gooconf.Config argument and sets them in the expected destination.
func (c *Config) UnmapEnv(gooconf *gofig.Config) error {
	vals := structMappings(c)
	for i := range vals {
		for j := range vals[i].EnvVars {
			// we need to guard against v != "" because this is the condition that checks that the value is set from the environment.
			// the `ok` guard is not enough, apparently.
			if v, ok := gooconf.GetValue(vals[i].EnvVars[j]); ok && v != "" {

				// get the destination type from destination
				switch reflect.ValueOf(vals[i].Destination).Type().String() {
				case "*bool":
					r := gooconf.Bool(vals[i].EnvVars[j])
					*vals[i].Destination.(*bool) = r
				case "*string":
					r := gooconf.String(vals[i].EnvVars[j])
					*vals[i].Destination.(*string) = r
				case "*int":
					r := gooconf.Int(vals[i].EnvVars[j])
					*vals[i].Destination.(*int) = r
				case "*float64":
					// defaults to float64
					r := gooconf.Float(vals[i].EnvVars[j])
					*vals[i].Destination.(*float64) = r
				default:
					// it is unlikely we will ever get here. Let this serve more as a runtime check for when debugging.
					return fmt.Errorf("invalid type for env var: `%v`", vals[i].EnvVars[j])
				}
			}
		}
	}

	return nil
}
