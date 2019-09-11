package debug

// Option configures an assets option.
type Option func(*debug)

// WithToken returns an option to set a token.
func WithToken(val string) Option {
	return func(d *debug) {
		d.token = val
	}
}

// WithPprof returns an option to enable pprof.
func WithPprof(val bool) Option {
	return func(d *debug) {
		d.pprof = val
	}
}
