package runtime

import (
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/urfave/cli/v2"
)

// Options is a runtime option
type Options struct {
	Services []string
	Logger   log.Logger
	Context  *cli.Context
}

// Option undocumented
type Option func(o *Options)

// Services option
func Services(s []string) Option {
	return func(o *Options) {
		o.Services = append(o.Services, s...)
	}
}

// Context option
func Context(c *cli.Context) Option {
	return func(o *Options) {
		o.Context = c
	}
}
