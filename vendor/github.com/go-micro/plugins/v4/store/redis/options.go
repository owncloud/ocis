package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4/store"
)

type redisOptionsContextKey struct{}

// WithRedisOptions sets advanced options for redis.
func WithRedisOptions(options redis.UniversalOptions) store.Option {
	return func(o *store.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}

		o.Context = context.WithValue(o.Context, redisOptionsContextKey{}, options)
	}
}

func newUniversalClient(o store.Options) redis.UniversalClient {
	if o.Context == nil {
		o.Context = context.Background()
	}

	opts, ok := o.Context.Value(redisOptionsContextKey{}).(redis.UniversalOptions)
	if !ok && len(o.Nodes) <= 1 {
		addr := "redis://127.0.0.1:6379"
		if len(o.Nodes) > 0 {
			addr = o.Nodes[0]
		}

		redisOptions, err := redis.ParseURL(addr)
		if err != nil {
			redisOptions = &redis.Options{Addr: addr}
		}

		return redis.NewClient(redisOptions)
	}

	if len(opts.Addrs) == 0 && len(o.Nodes) > 0 {
		opts.Addrs = o.Nodes
	}

	return redis.NewUniversalClient(&opts)
}
