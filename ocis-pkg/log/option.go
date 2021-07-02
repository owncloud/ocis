package log

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Name   string
	Level  string
	Pretty bool
	Color  bool
	File   string
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{
		Name:   "ocis",
		Level:  "info",
		Pretty: true,
		Color:  true,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Name provides a function to set the name option.
func Name(val string) Option {
	return func(o *Options) {
		o.Name = val
	}
}

// Level provides a function to set the level option.
func Level(val string) Option {
	return func(o *Options) {
		o.Level = val
	}
}

// Pretty provides a function to set the pretty option.
func Pretty(val bool) Option {
	return func(o *Options) {
		o.Pretty = val
	}
}

// Color provides a function to set the color option.
func Color(val bool) Option {
	return func(o *Options) {
		o.Color = val
	}
}

// File provides a function to set the color option.
func File(val string) Option {
	return func(o *Options) {
		o.File = val
	}
}
