package bleve

type GetIndexOption func(o *GetIndexOptions)

type GetIndexOptions struct {
	ReadOnly bool
}

func ReadOnly(b bool) GetIndexOption {
	return func(o *GetIndexOptions) {
		o.ReadOnly = b
	}
}

func newGetIndexOptions(opts ...GetIndexOption) GetIndexOptions {
	o := GetIndexOptions{}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
