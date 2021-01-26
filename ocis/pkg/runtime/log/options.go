package log

import "github.com/rs/zerolog"

// Options are the configurable options for a Controller.
type Options struct {
	Level  zerolog.Level
	Pretty bool
}

// Option represents an option.
type Option func(o *Options)

// NewOptions returns a new Options struct.
func NewOptions() *Options {
	return &Options{
		Level: zerolog.DebugLevel,
	}
}

// WithPretty sets the pretty option.
func WithPretty(pretty bool) Option {
	return func(o *Options) {
		o.Pretty = pretty
	}
}
