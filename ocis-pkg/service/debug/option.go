package debug

import (
	"net/http"

	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger               log.Logger
	Name                 string
	Version              string
	Address              string
	Token                string
	Pprof                bool
	Zpages               bool
	Health               func(http.ResponseWriter, *http.Request)
	Ready                func(http.ResponseWriter, *http.Request)
	ConfigDump           func(http.ResponseWriter, *http.Request)
	CorsAllowedOrigins   []string
	CorsAllowedMethods   []string
	CorsAllowedHeaders   []string
	CorsAllowCredentials bool
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the logger option.
func Logger(l log.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// Name provides a function to set the name option.
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Version provides a function to set the version option.
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Address provides a function to set the address option.
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// Token provides a function to set the token option.
func Token(t string) Option {
	return func(o *Options) {
		o.Token = t
	}
}

// Pprof provides a function to set the pprof option.
func Pprof(p bool) Option {
	return func(o *Options) {
		o.Pprof = p
	}
}

// Zpages provides a function to set the zpages option.
func Zpages(z bool) Option {
	return func(o *Options) {
		o.Zpages = z
	}
}

// Health provides a function to set the health option.
func Health(h func(http.ResponseWriter, *http.Request)) Option {
	return func(o *Options) {
		o.Health = h
	}
}

// Ready provides a function to set the ready option.
func Ready(r func(http.ResponseWriter, *http.Request)) Option {
	return func(o *Options) {
		o.Ready = r
	}
}

// ConfigDump to be documented.
func ConfigDump(r func(http.ResponseWriter, *http.Request)) Option {
	return func(o *Options) {
		o.ConfigDump = r
	}
}

// CorsAllowedOrigins provides a function to set the CorsAllowedOrigin option.
func CorsAllowedOrigins(origins []string) Option {
	return func(o *Options) {
		o.CorsAllowedOrigins = origins
	}
}

// CorsAllowedMethods provides a function to set the CorsAllowedMethods option.
func CorsAllowedMethods(methods []string) Option {
	return func(o *Options) {
		o.CorsAllowedMethods = methods
	}
}

// CorsAllowedHeaders provides a function to set the CorsAllowedHeaders option.
func CorsAllowedHeaders(headers []string) Option {
	return func(o *Options) {
		o.CorsAllowedHeaders = headers
	}
}

// CorsAllowCredentials provides a function to set the CorsAllowAllowCredential option.
func CorsAllowCredentials(allow bool) Option {
	return func(o *Options) {
		o.CorsAllowCredentials = allow
	}
}
