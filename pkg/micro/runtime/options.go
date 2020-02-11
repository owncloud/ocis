package runtime

import (
	gorun "github.com/micro/go-micro/v2/runtime"
	"github.com/owncloud/ocis-pkg/v2/log"
)

// Options is a runtime option
type Options struct {
	Services     []string
	Logger       log.Logger
	MicroRuntime *gorun.Runtime
}

// Option undocummented
type Option func(o *Options)

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Services option
func Services(s []string) Option {
	return func(o *Options) {
		o.Services = append(o.Services, s...)
	}
}

// Logger option
func Logger(l log.Logger) Option {
	return func(o *Options) {
		o.Logger = l
	}
}

// MicroRuntime option
func MicroRuntime(rt *gorun.Runtime) Option {
	return func(o *Options) {
		o.MicroRuntime = rt
	}
}
