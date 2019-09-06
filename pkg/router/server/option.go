package server

// Option configures an assets option.
type Option func(*server)

// WithRoot returns an option to set a root.
func WithRoot(val string) Option {
	return func(s *server) {
		s.root = val
	}
}

// WithPath returns an option to set a path.
func WithPath(val string) Option {
	return func(s *server) {
		s.path = val
	}
}

// WithCustom returns an option to set a path.
func WithCustom(val string) Option {
	return func(s *server) {
		s.custom = val
	}
}

// WithServer returns an option to set a path.
func WithServer(val string) Option {
	return func(s *server) {
		s.server = val
	}
}

// WithTheme returns an option to set a path.
func WithTheme(val string) Option {
	return func(s *server) {
		s.theme = val
	}
}

// WithVersion returns an option to set a path.
func WithVersion(val string) Option {
	return func(s *server) {
		s.version = val
	}
}

// WithClient returns an option to set a path.
func WithClient(val string) Option {
	return func(s *server) {
		s.client = val
	}
}

// WithApps returns an option to set a path.
func WithApps(val []string) Option {
	return func(s *server) {
		s.apps = val
	}
}
