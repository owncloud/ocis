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

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func WithJWTSecret(val string) Option {
	return func(o *Options) {
		o.JWTSecret = val
	}
}

func WithDataURL(val string) Option {
	return func(o *Options) {
		o.DataURL = val
	}
}

func WithDataPrefix(val string) Option {
	return func(o *Options) {
		o.DataPrefix = val
	}
}

func WithEntityDirName(val string) Option {
	return func(o *Options) {
		o.EntityDirName = val
	}
}

func WithDataDir(val string) Option {
	return func(o *Options) {
		o.DataDir = val
	}
}

func WithTypeName(val string) Option {
	return func(o *Options) {
		o.TypeName = val
	}
}

func WithIndexBy(val string) Option {
	return func(o *Options) {
		o.IndexBy = val
	}
}

func WithIndexBaseDir(val string) Option {
	return func(o *Options) {
		o.IndexBaseDir = val
	}
}

func WithFilesDir(val string) Option {
	return func(o *Options) {
		o.FilesDir = val
	}
}
func WithProviderAddr(val string) Option {
	return func(o *Options) {
		o.ProviderAddr = val
	}
}
