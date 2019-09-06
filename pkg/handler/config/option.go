package config

// Option configures an assets option.
type Option func(*config)

// WithServer returns an option to set server.
func WithServer(val string) Option {
	return func(c *config) {
		c.server = val
	}
}

// WithTheme returns an option to set theme.
func WithTheme(val string) Option {
	return func(c *config) {
		c.theme = val
	}
}

// WithVersion returns an option to set version.
func WithVersion(val string) Option {
	return func(c *config) {
		c.version = val
	}
}

// WithClient returns an option to set client id.
func WithClient(val string) Option {
	return func(c *config) {
		c.client = val
	}
}

// WithApps returns an option to set apps.
func WithApps(val []string) Option {
	return func(c *config) {
		c.apps = val
	}
}

// WithCustom returns an option to set custom config.
func WithCustom(val string) Option {
	return func(c *config) {
		c.custom = val
	}
}
