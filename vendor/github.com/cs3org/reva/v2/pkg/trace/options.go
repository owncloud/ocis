package trace

import "google.golang.org/grpc/credentials"

// Options for trace
type Options struct {
	Enabled              bool
	Insecure             bool
	Exporter             string
	Collector            string
	Endpoint             string
	ServiceName          string
	TransportCredentials credentials.TransportCredentials
}

// Option for trace
type Option func(o *Options)

// WithEnabled option
func WithEnabled() Option {
	return func(o *Options) {
		o.Enabled = true
	}
}

// WithExporter option
func WithExporter(v string) Option {
	return func(o *Options) {
		o.Exporter = v
	}
}

// WithInsecure option
func WithInsecure() Option {
	return func(o *Options) {
		o.Insecure = true
	}
}

// WithCollector option
func WithCollector(v string) Option {
	return func(o *Options) {
		o.Collector = v
	}
}

// WithEndpoint option
func WithEndpoint(v string) Option {
	return func(o *Options) {
		o.Endpoint = v
	}
}

// WithServiceName option
func WithServiceName(v string) Option {
	return func(o *Options) {
		o.ServiceName = v
	}
}
