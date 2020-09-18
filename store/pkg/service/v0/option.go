package service

import (
	"github.com/owncloud/ocis-pkg/v2/log"
	"github.com/owncloud/ocis-store/pkg/config"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger log.Logger
	Config *config.Config

	Database, Table string
	Nodes           []string
}

func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// Database configures the database option.
func Database(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// Table configures the Table option.
func Table(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// Nodes configures the Nodes option.
func Nodes(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// Config configures the Config option.
func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}
