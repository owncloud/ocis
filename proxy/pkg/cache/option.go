package cache

import "time"

// Options are all the possible options.
type Options struct {
	size int
	ttl  time.Duration
}

// Option mutates option
type Option func(*Options)

// Size configures the size of the cache in items.
func Size(s int) Option {
	return func(o *Options) {
		o.size = s
	}
}

// TTL rebuilds the cache after the configured duration.
func TTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.ttl = ttl
	}
}

func newOptions(opts ...Option) Options {
	o := Options{}

	for _, v := range opts {
		v(&o)
	}

	return o
}
