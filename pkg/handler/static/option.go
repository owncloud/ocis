package static

// Option configures an assets option.
type Option func(*static)

// WithRoot returns an option to set a root.
func WithRoot(val string) Option {
	return func(s *static) {
		s.root = val
	}
}

// WithPath returns an option to set a path.
func WithPath(val string) Option {
	return func(s *static) {
		s.path = val
	}
}
