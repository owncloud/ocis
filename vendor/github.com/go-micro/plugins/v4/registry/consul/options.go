package consul

import (
	"context"
	"fmt"
	"time"

	consul "github.com/hashicorp/consul/api"
	"go-micro.dev/v4/registry"
)

// Connect specifies services should be registered as Consul Connect services.
func Connect() registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_connect", true)
	}
}

func Config(c *consul.Config) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_config", c)
	}
}

// AllowStale sets whether any Consul server (non-leader) can service
// a read. This allows for lower latency and higher throughput
// at the cost of potentially stale data.
// Works similar to Consul DNS Config option [1].
// Defaults to true.
//
// [1] https://www.consul.io/docs/agent/options.html#allow_stale
func AllowStale(v bool) registry.Option {
	return func(o *registry.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_allow_stale", v)
	}
}

// QueryOptions specifies the QueryOptions to be used when calling
// Consul. See `Consul API` for more information [1].
//
// [1] https://godoc.org/github.com/hashicorp/consul/api#QueryOptions
func QueryOptions(q *consul.QueryOptions) registry.Option {
	return func(o *registry.Options) {
		if q == nil {
			return
		}
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_query_options", q)
	}
}

// TCPCheck will tell the service provider to check the service address
// and port every `t` interval. It will enabled only if `t` is greater than 0.
// See `TCP + Interval` for more information [1].
//
// [1] https://www.consul.io/docs/agent/checks.html
func TCPCheck(t time.Duration) registry.Option {
	return func(o *registry.Options) {
		if t <= time.Duration(0) {
			return
		}
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, "consul_tcp_check", t)
	}
}

// HTTPCheck will tell the service provider to invoke the health check endpoint
// with an interval and timeout. It will be enabled only if interval and
// timeout are greater than 0.
// See `HTTP + Interval` for more information [1].
//
// [1] https://www.consul.io/docs/agent/checks.html
func HTTPCheck(protocol, port, httpEndpoint string, interval, timeout time.Duration) registry.Option {
	return func(o *registry.Options) {
		if interval <= time.Duration(0) || timeout <= time.Duration(0) {
			return
		}
		if o.Context == nil {
			o.Context = context.Background()
		}
		check := consul.AgentServiceCheck{
			HTTP:     fmt.Sprintf("%s://{host}:%s%s", protocol, port, httpEndpoint),
			Interval: fmt.Sprintf("%v", interval),
			Timeout:  fmt.Sprintf("%v", timeout),
		}
		o.Context = context.WithValue(o.Context, "consul_http_check_config", check)
	}
}
