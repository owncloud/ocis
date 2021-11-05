package config

import (
	"context"
	"fmt"
	"reflect"

	gofig "github.com/gookit/config/v2"
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

// CORS defines the available cors configuration.
type CORS struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// HTTP defines the available http configuration.
type HTTP struct {
	Addr string
	Root string
	CORS CORS
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

// TokenManager is the config for using the reva token manager
type TokenManager struct {
	JWTSecret string
}

// IdentityManagement keeps track of the OIDC address. This is because Reva requisite of uniqueness for users
// is based in the combination of IDP hostname + UserID. For more information see:
// https://github.com/cs3org/reva/blob/4fd0229f13fae5bc9684556a82dbbd0eced65ef9/pkg/storage/utils/decomposedfs/node/node.go#L856-L865
type IdentityManagement struct {
	Address string
}

// Config combines all available configuration parts.
type Config struct {
	File               string
	Log                Log
	Debug              Debug
	HTTP               HTTP
	Tracing            Tracing
	TokenManager       TokenManager
	Service            Service
	AccountBackend     string
	RevaAddress        string
	StorageUsersDriver string
	MachineAuthAPIKey  string
	IdentityManagement IdentityManagement

	Context    context.Context
	Supervised bool
}

// New initializes a new configuration with or without defaults.
func New() *Config {
	return &Config{}
}

// DefaultConfig provides default values for a config struct.
func DefaultConfig() *Config {
	return &Config{
		Log: Log{},
		Debug: Debug{
			Addr:   "127.0.0.1:9114",
			Token:  "",
			Pprof:  false,
			Zpages: false,
		},
		HTTP: HTTP{
			Addr: "127.0.0.1:9110",
			Root: "/ocs",
			CORS: CORS{
				AllowedOrigins:   []string{"*"},
				AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"},
				AllowCredentials: true,
			},
		},
		Tracing: Tracing{
			Enabled:   false,
			Type:      "jaeger",
			Endpoint:  "",
			Collector: "",
			Service:   "ocs",
		},
		TokenManager: TokenManager{
			JWTSecret: "Pive-Fumkiu4",
		},
		Service: Service{
			Name:      "ocs",
			Namespace: "com.owncloud.web",
		},
		AccountBackend:     "accounts",
		RevaAddress:        "127.0.0.1:9142",
		StorageUsersDriver: "ocis",
		MachineAuthAPIKey:  "change-me-please",
		IdentityManagement: IdentityManagement{
			Address: "https://localhost:9200",
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
