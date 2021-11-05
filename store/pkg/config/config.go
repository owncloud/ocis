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
	Level  string
	Pretty bool
	Color  bool
	File   string
}

// Debug defines the available debug configuration.
type Debug struct {
	Addr   string
	Token  string
	Pprof  bool
	Zpages bool
}

// GRPC defines the available grpc configuration.
type GRPC struct {
	Addr string
	Root string
}

// Service defines the available service configuration.
type Service struct {
	Name      string
	Namespace string
	Version   string
}

// Tracing defines the available tracing configuration.
type Tracing struct {
	Enabled   bool
	Type      string
	Endpoint  string
	Collector string
	Service   string
}

// Config combines all available configuration parts.
type Config struct {
	File     string
	Log      Log
	Debug    Debug
	GRPC     GRPC
	Tracing  Tracing
	Datapath string
	Service  Service

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

func DefaultConfig() *Config {
	return &Config{
		Log: Log{},
		Debug: Debug{
			Addr:   "127.0.0.1:9464",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		GRPC: GRPC{
			Addr: "127.0.0.1:9460",
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "store",
		},
		Datapath: path.Join(defaults.BaseDataPath(), "store"),
		Service: Service{
			Name:      "store",
			Namespace: "com.owncloud.api",
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
