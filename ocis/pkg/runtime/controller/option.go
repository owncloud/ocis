package controller

import (
	"github.com/owncloud/ocis/ocis/pkg/runtime/config"
	"github.com/rs/zerolog"
)

// Options are the configurable options for a Controller.
type Options struct {
	Bin     string
	Restart bool
	Config  *config.Config
	Log     *zerolog.Logger
}

// Option represents an option.
type Option func(o *Options)

// NewOptions returns a new Options struct.
func NewOptions() Options {
	return Options{}
}

// WithConfig sets Controller config.
func WithConfig(cfg *config.Config) Option {
	return func(o *Options) {
		o.Config = cfg
	}
}

// WithLog sets Controller config.
func WithLog(l *zerolog.Logger) Option {
	return func(o *Options) {
		o.Log = l
	}
}
