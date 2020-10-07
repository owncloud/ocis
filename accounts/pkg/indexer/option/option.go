package option

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	// Disk Options
	TypeName      string
	IndexBy       string
	FilesDir      string
	IndexBaseDir  string
	DataDir       string
	EntityDirName string

	// CS3 options
	DataURL      string
	DataPrefix   string
	JWTSecret    string
	ProviderAddr string
}

// WithJWTSecret sets the JWTSecret field.
func WithJWTSecret(val string) Option {
	return func(o *Options) {
		o.JWTSecret = val
	}
}

// WithDataURL sets the DataURl field.
func WithDataURL(val string) Option {
	return func(o *Options) {
		o.DataURL = val
	}
}

// WithDataPrefix sets the DataPrefix field.
func WithDataPrefix(val string) Option {
	return func(o *Options) {
		o.DataPrefix = val
	}
}

// WithEntityDirName sets the EntityDirName field.
func WithEntityDirName(val string) Option {
	return func(o *Options) {
		o.EntityDirName = val
	}
}

// WithDataDir sets the DataDir option.
func WithDataDir(val string) Option {
	return func(o *Options) {
		o.DataDir = val
	}
}

// WithTypeName sets the TypeName option.
func WithTypeName(val string) Option {
	return func(o *Options) {
		o.TypeName = val
	}
}

// WithIndexBy sets the option IndexBy
func WithIndexBy(val string) Option {
	return func(o *Options) {
		o.IndexBy = val
	}
}

// WithFilesDir sets the option FilesDir
func WithFilesDir(val string) Option {
	return func(o *Options) {
		o.FilesDir = val
	}
}

// WithProviderAddr sets the option ProviderAddr
func WithProviderAddr(val string) Option {
	return func(o *Options) {
		o.ProviderAddr = val
	}
}
