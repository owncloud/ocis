package bleve

// GetIndexOption is a function that sets some option for the GetIndex method.
type GetIndexOption func(o *GetIndexOptions)

// GetIndexOptions contains the options for the GetIndex method.
type GetIndexOptions struct {
	ReadOnly bool
}

// ReadOnly is an option to opens the index in read-only mode.
// This option should allow running multiple read-only operations in parallel.
// The behavior of write operations is not defined when this option is used.
func ReadOnly(b bool) GetIndexOption {
	return func(o *GetIndexOptions) {
		o.ReadOnly = b
	}
}

// newGetIndexOptions creates a new GetIndexOptions with the given options.
func newGetIndexOptions(opts ...GetIndexOption) GetIndexOptions {
	o := GetIndexOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
