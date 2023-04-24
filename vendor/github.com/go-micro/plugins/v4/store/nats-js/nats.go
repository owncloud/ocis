package natsjs

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"go-micro.dev/v4/store"
	"go-micro.dev/v4/util/cmd"
)

var (
	ErrBucketNotFound = errors.New("Bucket (database) not found")
)

type natsStore struct {
	sync.Once
	sync.RWMutex

	ttl         time.Duration
	storageType nats.StorageType
	description string

	opts            store.Options
	nopts           nats.Options
	jsopts          []nats.JSOpt
	objStoreConfigs []*nats.ObjectStoreConfig

	conn    *nats.Conn
	js      nats.JetStreamContext
	buckets map[string]nats.ObjectStore
}

func init() {
	cmd.DefaultStores["natsjs"] = NewStore
}

// NewStore will create a new NATS JetStream Object Store
func NewStore(opts ...store.Option) store.Store {
	options := store.Options{
		Nodes:    []string{},
		Database: "default",
		Table:    "",
		Context:  context.Background(),
	}

	n := &natsStore{
		description:     "Object storage administered by go-micro store plugin",
		opts:            options,
		jsopts:          []nats.JSOpt{},
		objStoreConfigs: []*nats.ObjectStoreConfig{},
		buckets:         map[string]nats.ObjectStore{},
		storageType:     nats.FileStorage,
	}

	n.setOption(opts...)

	return n
}

// Init initialises the store. It must perform any required setup on the backing storage implementation and check that it is ready for use, returning any errors.
func (n *natsStore) Init(opts ...store.Option) error {
	n.setOption(opts...)

	// Connect to NATS servers
	conn, err := n.nopts.Connect()
	if err != nil {
		return errors.Wrap(err, "Failed to connect to NATS Server")
	}
	n.conn = conn

	// Create JetStream context
	js, err := conn.JetStream(n.jsopts...)
	if err != nil {
		return errors.Wrap(err, "Failed to create JetStream context")
	}
	n.js = js

	// Create default config if no configs present
	if len(n.objStoreConfigs) == 0 {
		n.objStoreConfigs = append(n.objStoreConfigs, &nats.ObjectStoreConfig{
			Bucket:      n.opts.Database,
			Description: n.description,
			TTL:         n.ttl,
			Storage:     n.storageType,
		})
	}

	// Create objest store buckets
	for _, cfg := range n.objStoreConfigs {
		store, err := js.CreateObjectStore(cfg)
		if err == nats.ErrStreamNameAlreadyInUse {
			store, err = n.js.ObjectStore(cfg.Bucket)
		}
		if err != nil {
			return errors.Wrapf(err, "Failed to create bucket (%s)", cfg.Bucket)
		}
		n.buckets[cfg.Bucket] = store
	}

	return nil
}

func (n *natsStore) setOption(opts ...store.Option) {
	for _, o := range opts {
		o(&n.opts)
	}

	n.Once.Do(func() {
		n.nopts = nats.GetDefaultOptions()
	})

	// Extract options from context
	if nopts, ok := n.opts.Context.Value(natsOptionsKey{}).(nats.Options); ok {
		n.nopts = nopts
	}

	if jsopts, ok := n.opts.Context.Value(jsOptionsKey{}).([]nats.JSOpt); ok {
		n.jsopts = append(n.jsopts, jsopts...)
	}

	if cfg, ok := n.opts.Context.Value(objOptionsKey{}).([]*nats.ObjectStoreConfig); ok {
		n.objStoreConfigs = append(n.objStoreConfigs, cfg...)
	}

	if ttl, ok := n.opts.Context.Value(ttlOptionsKey{}).(time.Duration); ok {
		n.ttl = ttl
	}

	if sType, ok := n.opts.Context.Value(memoryOptionsKey{}).(nats.StorageType); ok {
		n.storageType = sType
	}

	if text, ok := n.opts.Context.Value(descriptionOptionsKey{}).(string); ok {
		n.description = text
	}

	// Assign store option server addresses to nats options
	if len(n.opts.Nodes) > 0 {
		n.nopts.Url = ""
		n.nopts.Servers = n.opts.Nodes
	}

	if len(n.nopts.Servers) == 0 && n.nopts.Url == "" {
		n.nopts.Url = nats.DefaultURL
	}
}

// Options allows you to view the current options.
func (n *natsStore) Options() store.Options {
	return n.opts
}

// Read takes a single key name and optional ReadOptions. It returns matching []*Record or an error.
func (n *natsStore) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	if n.conn == nil {
		if err := n.Init(); err != nil {
			return nil, err
		}
	}

	opt := store.ReadOptions{}

	for _, o := range opts {
		o(&opt)
	}

	if opt.Database == "" {
		opt.Database = n.opts.Database
	}
	if opt.Table == "" {
		opt.Table = n.opts.Table
	}

	bucket, ok := n.buckets[opt.Database]
	if !ok {
		return nil, ErrBucketNotFound
	}

	var keys []string
	objects, err := bucket.List()
	if err == nats.ErrNoObjectsFound {
		return []*store.Record{}, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Failed to list objects")
	}

	for _, obj := range objects {
		name := obj.Name
		if (!opt.Prefix && !opt.Suffix) && getKey(key, opt.Table) != name {
			continue
		}
		if opt.Prefix && !strings.HasPrefix(name, getKey(key, opt.Table)) {
			continue
		}

		if opt.Suffix && !strings.HasSuffix(name, key) {
			continue
		}
		keys = append(keys, name)
	}

	records := []*store.Record{}
	for _, key := range keys {
		obj, err := bucket.Get(key)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to get object from bucket")
		}

		b, err := io.ReadAll(obj)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to read returned bytes")
		}

		info, err := obj.Info()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to fetch record info")
		}

		metadata := map[string]interface{}{}
		for key, value := range info.Headers {
			var val interface{}
			if err := json.Unmarshal([]byte(value[0]), &val); err != nil {
				return nil, errors.Wrap(err, "Failed to JSON unmarshal metadata")
			}
			metadata[key] = val
		}

		records = append(records, &store.Record{
			Key:      key,
			Value:    b,
			Metadata: metadata,
		})

		// Why is there a close method?
		obj.Close()
	}

	if opt.Limit > 0 {
		return records[opt.Offset : opt.Offset+opt.Limit], nil
	}
	if opt.Offset > 0 {
		return records[opt.Offset:], nil
	}
	return records, nil
}

// Write writes a record to the store, and returns an error if the record was not written.
func (n *natsStore) Write(r *store.Record, opts ...store.WriteOption) error {
	if n.conn == nil {
		if err := n.Init(); err != nil {
			return err
		}
	}

	opt := store.WriteOptions{}

	for _, o := range opts {
		o(&opt)
	}

	if opt.Database == "" {
		opt.Database = n.opts.Database
	}
	if opt.Table == "" {
		opt.Table = n.opts.Table
	}

	store, ok := n.buckets[opt.Database]

	// Create new bucket if not exists
	if !ok {
		var err error
		store, err = n.createNewBucket(opt.Database)
		if err != nil {
			return err
		}
	}

	header := nats.Header{}
	for key, value := range r.Metadata {
		val, err := json.Marshal(value)
		if err != nil {
			return errors.Wrap(err, "Failed to JSON marshal metadata")
		}
		header.Set(key, string(val))
	}

	_, err := store.Put(&nats.ObjectMeta{
		Name:        getKey(r.Key, opt.Table),
		Description: "Store managed by go-micro",
		Headers:     header,
	}, bytes.NewReader(r.Value))

	if err != nil {
		return errors.Wrap(err, "Failed to store data in bucket")
	}

	return nil
}

// Delete removes the record with the corresponding key from the store.
func (n *natsStore) Delete(key string, opts ...store.DeleteOption) error {
	if n.conn == nil {
		if err := n.Init(); err != nil {
			return err
		}
	}

	opt := store.DeleteOptions{}

	for _, o := range opts {
		o(&opt)
	}

	if opt.Database == "" {
		opt.Database = n.opts.Database
	}
	if opt.Table == "" {
		opt.Table = n.opts.Table
	}

	if opt.Table == "DELETE_BUCKET" {
		delete(n.buckets, key)
		if err := n.js.DeleteObjectStore(key); err != nil {
			return errors.Wrap(err, "Failed to delete bucket")
		}
		return nil
	}

	store, ok := n.buckets[opt.Database]
	if !ok {
		return ErrBucketNotFound
	}

	if err := store.Delete(getKey(key, opt.Table)); err != nil {
		return errors.Wrap(err, "Failed to delete data")
	}
	return nil
}

// List returns any keys that match, or an empty list with no error if none matched.
func (n *natsStore) List(opts ...store.ListOption) ([]string, error) {
	if n.conn == nil {
		if err := n.Init(); err != nil {
			return nil, err
		}
	}

	opt := store.ListOptions{}

	for _, o := range opts {
		o(&opt)
	}

	if opt.Database == "" {
		opt.Database = n.opts.Database
	}
	if opt.Table == "" {
		opt.Table = n.opts.Table
	}

	store, ok := n.buckets[opt.Database]
	if !ok {
		return nil, ErrBucketNotFound
	}

	objects, err := store.List()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list keys in bucket")
	}

	var keys []string
	for _, obj := range objects {
		key := obj.Name

		if !strings.HasPrefix(key, getKey(opt.Prefix, opt.Table)) {
			continue
		}

		if !strings.HasSuffix(key, opt.Suffix) {
			continue
		}
		keys = append(keys, key)
	}

	if opt.Limit > 0 {
		return keys[opt.Offset : opt.Offset+opt.Limit], nil
	}
	if opt.Offset > 0 {
		return keys[opt.Offset:], nil
	}
	return keys, nil
}

// Close the store
func (n *natsStore) Close() error {
	n.conn.Close()
	return nil
}

// String returns the name of the implementation.
func (n *natsStore) String() string {
	return "NATS JetStream ObjectStore"
}

func getKey(key, table string) string {
	if table != "" {
		key = table + "_" + key
	}
	return key
}

func (n *natsStore) createNewBucket(name string) (nats.ObjectStore, error) {
	store, err := n.js.CreateObjectStore(&nats.ObjectStoreConfig{
		Bucket:      name,
		Description: n.description,
		TTL:         n.ttl,
		Storage:     n.storageType,
	})
	if err == nats.ErrStreamNameAlreadyInUse {
		store, err = n.js.ObjectStore(name)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create new bucket (%s)", name)
	}
	n.buckets[name] = store
	return store, err
}
