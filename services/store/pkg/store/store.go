package store

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	storemsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/store/v0"
	storesvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/store/v0"
	"go-micro.dev/v4/client"
	merrors "go-micro.dev/v4/errors"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/store"
	"go-micro.dev/v4/util/cmd"
)

// DefaultDatabase is the namespace that the store
// will use if no namespace is provided.
var (
	DefaultDatabase = "proxy"
	DefaultTable    = "signing-keys"
)

type oss struct {
	ctx     context.Context
	options store.Options
	svc     storesvc.StoreService
}

func init() {
	cmd.DefaultStores["ocisstoreservice"] = NewStore
}

// NewStore returns a micro store.Store wrapper to access the micro store service.
// It only implements the minimal Read and Write options that are used by the proxy and ocs services
// Deprecated: use a different micro.Store implementation like nats-js-ks
func NewStore(opts ...store.Option) store.Store {
	options := store.Options{
		Context:  context.Background(),
		Database: DefaultDatabase,
		Table:    DefaultTable,
		Logger:   logger.DefaultLogger,
		Nodes:    []string{"com.owncloud.api.store"},
	}

	for _, o := range opts {
		o(&options)
	}

	c, ok := options.Context.Value(grpcClientContextKey{}).(client.Client)
	if !ok {
		var err error
		c, err = grpc.NewClient()
		if err != nil {
			options.Logger.Fields(map[string]interface{}{"err": err}).Log(logger.FatalLevel, "ocisstoreservice could not create new grpc client")
		}
	}
	svc := storesvc.NewStoreService(options.Nodes[0], c)

	s := &oss{
		ctx:     context.Background(),
		options: options,
		svc:     svc,
	}

	return s
}

// Init initializes the store by configuring a storeservice and initializing
// a grpc client if it has not been passed as a context option.
func (s *oss) Init(opts ...store.Option) error {
	for _, o := range opts {
		o(&s.options)
	}
	return s.configure()
}

func (s *oss) configure() error {
	c, ok := s.options.Context.Value(grpcClientContextKey{}).(client.Client)
	if !ok {
		var err error
		c, err = grpc.NewClient()
		if err != nil {
			logger.Fatal("ocisstoreservice could not create new grpc client:", err)
		}
	}
	if len(s.options.Nodes) < 1 {
		return errors.New("no node configured")
	}
	s.svc = storesvc.NewStoreService(s.options.Nodes[0], c)
	return nil
}

// Options allows you to view the current options.
func (s *oss) Options() store.Options {
	return s.options
}

// Read takes a single key name and optional ReadOptions. It returns matching []*Record or an error.
// Only database and table options are used.
func (s *oss) Read(key string, opts ...store.ReadOption) ([]*store.Record, error) {
	options := store.ReadOptions{
		Database: s.options.Database,
		Table:    s.options.Table,
	}

	for _, o := range opts {
		o(&options)
	}

	res, err := s.svc.Read(context.Background(), &storesvc.ReadRequest{
		Options: &storemsg.ReadOptions{
			Database: options.Database,
			Table:    options.Table,
			// Other options ignored
		},
		Key: key,
	})

	if err != nil {
		e := merrors.Parse(err.Error())
		if e.Code == http.StatusNotFound {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	records := make([]*store.Record, 0, len(res.Records))
	for _, record := range res.Records {
		r := &store.Record{
			Key:      record.Key,
			Value:    record.Value,
			Metadata: map[string]interface{}{},
			Expiry:   time.Duration(record.Expiry),
		}
		for k, v := range record.Metadata {
			r.Metadata[k] = v.Value // we only support string
		}
		records = append(records, r)
	}
	return records, nil
}

// Write() writes a record to the store, and returns an error if the record was not written.
func (s *oss) Write(r *store.Record, opts ...store.WriteOption) error {
	options := store.WriteOptions{
		Database: s.options.Database,
		Table:    s.options.Table,
	}

	for _, o := range opts {
		o(&options)
	}
	_, err := s.svc.Write(context.Background(), &storesvc.WriteRequest{
		Options: &storemsg.WriteOptions{
			Database: options.Database,
			Table:    options.Table,
		},
		Record: &storemsg.Record{
			Key:   r.Key,
			Value: r.Value,
			// No expiry supported
		},
	})

	return err
}

// Delete is not implemented for the ocis service
func (s *oss) Delete(key string, opts ...store.DeleteOption) error {
	return errors.ErrUnsupported
}

// List is not implemented for the ocis service
func (s *oss) List(opts ...store.ListOption) ([]string, error) {
	return nil, errors.ErrUnsupported
}

// Close does nothing
func (s *oss) Close() error {
	return nil
}

// String returns the name of the implementation.
func (s *oss) String() string {
	return "ocisstoreservice"
}
