package config

import (
	"context"
	"fmt"
	"reflect"

	gofig "github.com/gookit/config/v2"
	"github.com/owncloud/ocis/ocis-pkg/shared"
)

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string `mapstructure:"addr"`
	Token  string `mapstructure:"token"`
	Pprof  bool   `mapstructure:"pprof"`
	Zpages bool   `mapstructure:"zpages"`
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string `mapstructure:"addr"`
	Root      string `mapstructure:"root"`
	Namespace string `mapstructure:"namespace"`
}

// Server configures a server.
type Server struct {
	Version string `mapstructure:"version"`
	Name    string `mapstructure:"name"`
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool   `mapstructure:"enabled"`
	Type      string `mapstructure:"type"`
	Endpoint  string `mapstructure:"endpoint"`
	Collector string `mapstructure:"collector"`
	Service   string `mapstructure:"service"`
}

// GraphExplorer defines the available graph-explorer configuration.
type GraphExplorer struct {
	ClientID     string `mapstructure:"client_id"`
	Issuer       string `mapstructure:"issuer"`
	GraphURLBase string `mapstructure:"graph_url_base"`
	GraphURLPath string `mapstructure:"graph_url_path"`
}

// Config combines all available configuration parts.
type Config struct {
	File          string        `mapstructure:"file"`
	Log           shared.Log    `mapstructure:"log"`
	Debug         Debug         `mapstructure:"debug"`
	HTTP          HTTP          `mapstructure:"http"`
	Server        Server        `mapstructure:"server"`
	Tracing       Tracing       `mapstructure:"tracing"`
	GraphExplorer GraphExplorer `mapstructure:"graph_explorer"`

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

// DefaultConfig provides with a working version of a config.
func DefaultConfig() *Config {
	return &Config{
		Log: shared.Log{},
		Debug: Debug{
			Addr:   "127.0.0.1:9136",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr:      "127.0.0.1:9135",
			Root:      "/graph-explorer",
			Namespace: "com.owncloud.web",
		},
		Server: Server{},
		Tracing: Tracing{
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "graph-explorer",
		},
		GraphExplorer: GraphExplorer{
			ClientID:     "ocis-explorer.js",
			Issuer:       "https://localhost:9200",
			GraphURLBase: "https://localhost:9200",
			GraphURLPath: "/graph",
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
