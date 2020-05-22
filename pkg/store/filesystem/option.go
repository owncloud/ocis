package store

// Options are the available configurable options on initialization.
type Options struct {
	UUID string
}

// Option captures configuration behavior.
type Option func(*Options)

// NewOptions build a new Options struct.
func NewOptions(o ...Option) *Options {
	opts := &Options{}

	for _, f := range o {
		f(opts)
	}

	return opts
}

// WithUUID sets UUID option.
func WithUUID(uuid string) Option {
	return func(o *Options) {
		o.UUID = uuid
	}
}
