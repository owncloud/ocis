package nats

import (
	nserver "github.com/nats-io/nats-server/v2/server"
)

// NatsOption configures the nats server
type NatsOption func(*nserver.Options)

// Host sets the host URL for the nats server
func Host(url string) NatsOption {
	return func(o *nserver.Options) {
		o.Host = url
	}
}

// Port sets the host URL for the nats server
func Port(port int) NatsOption {
	return func(o *nserver.Options) {
		o.Port = port
	}
}

// ClusterID sets the name for the nats cluster
func ClusterID(clusterID string) NatsOption {
	return func(o *nserver.Options) {
		o.Cluster.Name = clusterID
	}
}

// NatsOption configures the nats server
type JetStreamOption func(*nserver.JetStreamConfig)

// ClusterID sets the name for the nats cluster
func JetStreamStoreDir(StoreDir string) JetStreamOption {
	return func(o *nserver.JetStreamConfig) {
		o.StoreDir = StoreDir
	}
}
