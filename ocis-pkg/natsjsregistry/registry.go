// Package natsjsregistry implements a registry using natsjs object store
package natsjsregistry

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	natsjs "github.com/go-micro/plugins/v4/store/nats-js"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/store"
	"go-micro.dev/v4/util/cmd"
)

var _registryName = "natsjs"

func init() {
	cmd.DefaultRegistries[_registryName] = NewRegistry
}

// NewRegistry returns a new natsjs registry
func NewRegistry(opts ...registry.Option) registry.Registry {
	options := registry.Options{
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	exp, _ := options.Context.Value(expiryKey{}).(time.Duration)
	return &storeregistry{
		opts:   options,
		store:  natsjs.NewStore(storeOptions(options)...),
		typ:    _registryName,
		expiry: exp,
	}
}

type storeregistry struct {
	opts   registry.Options
	store  store.Store
	typ    string
	expiry time.Duration
}

// Init inits the registry
func (n *storeregistry) Init(opts ...registry.Option) error {
	for _, o := range opts {
		o(&n.opts)
	}
	return n.store.Init(storeOptions(n.opts)...)
}

// Options returns the configured options
func (n *storeregistry) Options() registry.Options {
	return n.opts
}

// Register adds a service to the registry
func (n *storeregistry) Register(s *registry.Service, _ ...registry.RegisterOption) error {
	if s == nil {
		return errors.New("wont store nil service")
	}
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return n.store.Write(&store.Record{
		Key:    s.Name,
		Value:  b,
		Expiry: n.expiry,
	})
}

// Deregister removes a service from the registry
func (n *storeregistry) Deregister(s *registry.Service, _ ...registry.DeregisterOption) error {
	return n.store.Delete(s.Name)
}

// GetService gets a specific service from the registry
func (n *storeregistry) GetService(s string, _ ...registry.GetOption) ([]*registry.Service, error) {
	recs, err := n.store.Read(s)
	if err != nil {
		return nil, err
	}
	svcs := make([]*registry.Service, 0, len(recs))
	for _, rec := range recs {
		var s registry.Service
		if err := json.Unmarshal(rec.Value, &s); err != nil {
			return nil, err
		}
		svcs = append(svcs, &s)
	}
	return svcs, nil
}

// ListServices lists all registered services
func (n *storeregistry) ListServices(...registry.ListOption) ([]*registry.Service, error) {
	keys, err := n.store.List()
	if err != nil {
		return nil, err
	}

	var svcs []*registry.Service
	for _, k := range keys {
		s, err := n.GetService(k)
		if err != nil {
			// TODO: continue ?
			return nil, err
		}
		svcs = append(svcs, s...)

	}
	return svcs, nil
}

// Watch allowes following the changes in the registry if it would be implemented
func (n *storeregistry) Watch(...registry.WatchOption) (registry.Watcher, error) {
	return nil, errors.New("watcher not implemented")
}

// String returns the name of the registry
func (n *storeregistry) String() string {
	return n.typ
}

func storeOptions(opts registry.Options) []store.Option {
	storeoptions := []store.Option{
		store.Nodes(opts.Addrs...),
		store.Database("service-registry"),
		store.Table("service-registry"),
	}
	if so, ok := opts.Context.Value(storeOptionsKey{}).([]store.Option); ok {
		storeoptions = append(storeoptions, so...)
	}
	return storeoptions
}
