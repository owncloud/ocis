// Package redis is a redis backed store implementation
package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/store"
	"go-micro.dev/v4/util/cmd"
)

// DefaultDatabase is the namespace that the store
// will use if no namespace is provided.
var (
	DefaultDatabase = "micro"
	DefaultTable    = "micro"
)

type rkv struct {
	ctx     context.Context
	options store.Options
	Client  redis.UniversalClient
}

func init() {
	cmd.DefaultStores["redis"] = NewStore
}

func (r *rkv) Init(opts ...store.Option) error {
	for _, o := range opts {
		o(&r.options)
	}

	return r.configure()
}

func (r *rkv) Close() error {
	return r.Client.Close()
}

func (r *rkv) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	options := store.ReadOptions{
		Table: r.options.Table,
	}

	for _, o := range opts {
		o(&options)
	}

	var keys []string

	var rkey string

	switch {
	case options.Prefix:
		rkey = fmt.Sprintf("%s%s*", options.Table, key)
	case options.Suffix:
		rkey = fmt.Sprintf("%s*%s", options.Table, key)
	default:
		keys = []string{fmt.Sprintf("%s%s", options.Table, key)}
	}

	if len(keys) == 0 {
		cursor := uint64(options.Offset)
		count := int64(options.Limit)

		for {
			var err error

			var ks []string

			ks, cursor, err = r.Client.Scan(r.ctx, cursor, rkey, count).Result()
			if err != nil {
				return nil, err
			}

			keys = append(keys, ks...)

			if cursor == 0 {
				break
			}
		}
	}

	records := make([]*store.Record, 0, len(keys))

	// read all keys, continue on error
	var val []byte

	var d time.Duration

	var err error

	for _, rkey = range keys {
		val, err = r.Client.Get(r.ctx, rkey).Bytes()
		if err != nil || val == nil {
			continue
		}

		d, err = r.Client.TTL(r.ctx, rkey).Result()
		if err != nil {
			continue
		}

		records = append(records, &store.Record{
			Key:    key,
			Value:  val,
			Expiry: d,
		})
	}

	if len(keys) == 1 {
		if errors.Is(err, redis.Nil) {
			return records, store.ErrNotFound
		}
		return records, err
	}

	// keys might have vanished since we scanned them, ignore errors
	return records, nil
}

func (r *rkv) Delete(key string, opts ...store.DeleteOption) error {
	options := store.DeleteOptions{
		Table: r.options.Table,
	}

	for _, o := range opts {
		o(&options)
	}

	rkey := fmt.Sprintf("%s%s", options.Table, key)

	return r.Client.Del(r.ctx, rkey).Err()
}

func (r *rkv) Write(record *store.Record, opts ...store.WriteOption) error {
	options := store.WriteOptions{
		Table: r.options.Table,
	}

	for _, o := range opts {
		o(&options)
	}

	rkey := fmt.Sprintf("%s%s", options.Table, record.Key)

	return r.Client.Set(r.ctx, rkey, record.Value, record.Expiry).Err()
}

func (r *rkv) List(opts ...store.ListOption) ([]string, error) {
	options := store.ListOptions{
		Table: r.options.Table,
	}

	for _, o := range opts {
		o(&options)
	}

	key := fmt.Sprintf("%s%s*%s", options.Table, options.Prefix, options.Suffix)

	cursor := uint64(options.Offset)

	count := int64(options.Limit)

	var allKeys []string

	var keys []string

	var err error

	for {
		keys, cursor, err = r.Client.Scan(r.ctx, cursor, key, count).Result()
		if err != nil {
			return nil, err
		}

		for i, key := range keys {
			keys[i] = strings.TrimPrefix(key, options.Table)
		}

		allKeys = append(allKeys, keys...)

		if cursor == 0 {
			break
		}
	}

	return allKeys, nil
}

func (r *rkv) Options() store.Options {
	return r.options
}

func (r *rkv) String() string {
	return "redis"
}

// NewStore returns a redis store.
func NewStore(opts ...store.Option) store.Store {
	options := store.Options{
		Database: DefaultDatabase,
		Table:    DefaultTable,
		Logger:   logger.DefaultLogger,
	}

	for _, o := range opts {
		o(&options)
	}

	s := &rkv{
		ctx:     context.Background(),
		options: options,
	}

	if err := s.configure(); err != nil {
		s.options.Logger.Log(logger.ErrorLevel, "Error configuring store ", err)
	}

	return s
}

func (r *rkv) configure() error {
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			return err
		}
	}

	r.Client = newUniversalClient(r.options)

	return nil
}
