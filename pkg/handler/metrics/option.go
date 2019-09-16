package metrics

// Option configures an assets option.
type Option func(*metrics)

// WithToken returns an option to set a token.
func WithToken(val string) Option {
	return func(m *metrics) {
		m.token = val
	}
}
