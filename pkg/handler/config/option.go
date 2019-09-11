package config

// Option configures an assets option.
type Option func(*config)

// WithConfig returns an option to set config.
func WithConfig(val string) Option {
	return func(c *config) {
		c.file = val
	}
}
