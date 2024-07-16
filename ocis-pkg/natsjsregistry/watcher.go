package natsjsregistry

import (
	"errors"

	"github.com/nats-io/nats.go"
	"go-micro.dev/v4/registry"
)

// NatsWatcher is the watcher of the nats interface
type NatsWatcher interface {
	Watch(bucket string) (nats.KeyWatcher, error)
}

// Watcher is used to keep track of changes in the registry
type Watcher struct {
	watch   nats.KeyWatcher
	updates <-chan nats.KeyValueEntry
	reg     *storeregistry
}

// NewWatcher returns a new watcher
func NewWatcher(s *storeregistry) (*Watcher, error) {
	w, ok := s.store.(NatsWatcher)
	if !ok {
		return nil, errors.New("store does not implement watcher interface")
	}

	watcher, err := w.Watch("service-registry")
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watch:   watcher,
		updates: watcher.Updates(),
		reg:     s,
	}, nil
}

// Next returns the next result. It is a blocking call
func (w *Watcher) Next() (*registry.Result, error) {
	kve := <-w.updates
	if kve == nil {
		return nil, errors.New("watcher stopped")
	}

	service, err := w.reg.getService(kve.Key())
	if err != nil {
		return nil, err
	}

	var action string
	switch kve.Operation() {
	default:
		action = "create"
	case nats.KeyValuePut:
		action = "create"
	case nats.KeyValueDelete:
		action = "delete"
	case nats.KeyValuePurge:
		action = "delete"
	}

	return &registry.Result{
		Service: service,
		Action:  action,
	}, nil
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	_ = w.watch.Stop()
}
