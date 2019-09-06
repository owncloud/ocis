package assets

// Option configures an assets option.
type Option func(*assets)

// WithPath returns an option to set custom assets path.
func WithPath(val string) Option {
	return func(a *assets) {
		a.path = val
	}
}
