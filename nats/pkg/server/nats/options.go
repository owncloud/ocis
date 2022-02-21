package nats

import (
	nserver "github.com/nats-io/nats-server/v2/server"
)

// Option configures the nats server
type Option func(*nserver.Options)

// Host sets the host URL for the nats server
func Host(url string) Option {
	return func(o *nserver.Options) {
		o.Host = url
	}
}

// Port sets the host URL for the nats server
func Port(port int) Option {
	return func(o *nserver.Options) {
		o.Port = port
	}
}

// ClusterID sets the name for the nats cluster
func ClusterID(clusterID string) Option {
	return func(o *nserver.Options) {
		o.Cluster.Name = clusterID
	}
}
