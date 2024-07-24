package natsjsregistry

import (
	"encoding/json"
	"errors"

	natsjskv "github.com/go-micro/plugins/v4/store/nats-js-kv"
	"github.com/nats-io/nats.go"
	"go-micro.dev/v4/registry"
)

// NatsWatcher is the watcher of the nats interface
type NatsWatcher interface {
	WatchAll(bucket string, opts ...nats.WatchOpt) (<-chan *natsjskv.StoreUpdate, func() error, error)
}

// Watcher is used to keep track of changes in the registry
type Watcher struct {
	updates <-chan *natsjskv.StoreUpdate
	stop    func() error
	reg     *storeregistry
}

// NewWatcher returns a new watcher
func NewWatcher(s *storeregistry) (*Watcher, error) {
	w, ok := s.store.(NatsWatcher)
	if !ok {
		return nil, errors.New("store does not implement watcher interface")
	}

	watcher, stop, err := w.WatchAll("service-registry")
	if err != nil {
		return nil, err
	}

	return &Watcher{
		updates: watcher,
		stop:    stop,
		reg:     s,
	}, nil
}

// Next returns the next result. It is a blocking call
func (w *Watcher) Next() (*registry.Result, error) {
	kve := <-w.updates
	if kve == nil {
		return nil, errors.New("watcher stopped")
	}

	var svc *registry.Service
	if err := json.Unmarshal(kve.Value.Data, svc); err != nil {
		return nil, err
	}

	return &registry.Result{
		Service: svc,
		Action:  kve.Action,
	}, nil
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	_ = w.stop()
}
