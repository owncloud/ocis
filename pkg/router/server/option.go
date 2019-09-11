package server

// Option configures an assets option.
type Option func(*server)

// WithRoot returns an option to set root.
func WithRoot(val string) Option {
	return func(s *server) {
		s.root = val
	}
}

// WithPath returns an option to set path.
func WithPath(val string) Option {
	return func(s *server) {
		s.path = val
	}
}

// WithConfig returns an option to set config.
func WithConfig(val string) Option {
	return func(s *server) {
		s.config = val
	}
}
