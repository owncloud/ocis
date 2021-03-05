package command

// Option defines a single modifier to an Options attribute.
type Option func(o *Options)

type Options struct {
	// LogPretty toggles pretty logging lines.
	LogPretty bool

	// LogColor toggles colored output.
	LogColor bool

	// LogLevel raises / decreases logging levels.
	LogLevel string
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// WithLogPretty toggles pretty output for a storage logger.
func WithLogPretty(v bool) Option {
	return func(o *Options) {
		o.LogPretty = v
	}
}

// WithLogColor toggles colored output for a storage logger.
func WithLogColor(v bool) Option {
	return func(o *Options) {
		o.LogColor = v
	}
}

// WithLogLevel toggles colored output for a storage logger.
func WithLogLevel(v string) Option {
	return func(o *Options) {
		o.LogLevel = v
	}
}
