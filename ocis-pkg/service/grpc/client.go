package grpc

import (
	"sync"
	"time"

	mgrpcc "github.com/go-micro/plugins/v4/client/grpc"
	mbreaker "github.com/go-micro/plugins/v4/wrapper/breaker/gobreaker"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"go-micro.dev/v4/client"
)

// DefaultClient is a custom oCIS grpc configured client.
var (
	defaultClient client.Client
	once          sync.Once
)

// DefaultClient returns the default grpc client.
func DefaultClient() client.Client {
	return getDefaultGrpcClient()
}

func getDefaultGrpcClient() client.Client {
	once.Do(func() {
		reg := registry.GetRegistry()

		defaultClient = mgrpcc.NewClient(
			client.Registry(reg),
			client.Wrap(mbreaker.NewClientWrapper()),
		)
	})
	return defaultClient
}

// Client returns a configurable client.
func Client(opts ...ClientOption) client.Client {
	copts := newClientOptions(opts...)
	reg := registry.GetRegistry()

	return mgrpcc.NewClient(
		client.Registry(reg),
		client.Wrap(mbreaker.NewClientWrapper()),
		client.RequestTimeout(copts.RequestTimeout),
	)
}

// ClientOption defines a single client option function.
type ClientOption func(o *ClientOptions)

// ClientOptions defines the available options for this package.
type ClientOptions struct {
	RequestTimeout time.Duration
}

// newClientOptions initializes the available client options.
func newClientOptions(opts ...ClientOption) ClientOptions {
	opt := ClientOptions{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// RequestTimeout provides a function to set the request timeout option.
func RequestTimeout(d time.Duration) ClientOption {
	return func(o *ClientOptions) {
		o.RequestTimeout = d
	}
}
