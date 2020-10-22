package option

// Option defines a single option function.
type Option func(o *Options)

// Bound represents a lower and upper bound range for an index.
// todo: if we would like to provide an upper bound then we would need to deal with ranges, in which case this is why the
// upper bound attribute is here.
type Bound struct {
	Lower, Upper int64
}

// Options defines the available options for this package.
type Options struct {
	CaseInsensitive bool
	Bound           *Bound

	// Disk Options
	TypeName      string
	IndexBy       string
	FilesDir      string
	IndexBaseDir  string
	DataDir       string
	EntityDirName string
	Entity        interface{}

	// CS3 options
	DataURL         string
	DataPrefix      string
	JWTSecret       string
	ProviderAddr    string
	ServiceUserUUID string
	ServiceUserName string
}

// CaseInsensitive sets the CaseInsensitive field.
func CaseInsensitive(val bool) Option {
	return func(o *Options) {
		o.CaseInsensitive = val
	}
}

// WithBounds sets the Bounds field.
func WithBounds(val *Bound) Option {
	return func(o *Options) {
		o.Bound = val
	}
}

// WithEntity sets the Entity field.
func WithEntity(val interface{}) Option {
	return func(o *Options) {
		o.Entity = val
	}
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

// WithIndexBy sets the option IndexBy.
func WithIndexBy(val string) Option {
	return func(o *Options) {
		o.IndexBy = val
	}
}

// WithFilesDir sets the option FilesDir.
func WithFilesDir(val string) Option {
	return func(o *Options) {
		o.FilesDir = val
	}
}

// WithProviderAddr sets the option ProviderAddr.
func WithProviderAddr(val string) Option {
	return func(o *Options) {
		o.ProviderAddr = val
	}
}

// WithServiceUserUUID sets the option ServiceUserUUID.
func WithServiceUserUUID(val string) Option {
	return func(o *Options) {
		o.ServiceUserUUID = val
	}
}

// WithServiceUserName sets the option ServiceUserName.
func WithServiceUserName(val string) Option {
	return func(o *Options) {
		o.ServiceUserName = val
	}
}
