package event

import (
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"go-micro.dev/v4/events"
)

// NewStream prepares the requested nats stream and returns it.
func NewStream(cfg *config.Config) (events.Stream, error) {
	return stream.NatsFromConfig(cfg.Service.Name, stream.NatsConfig{
		Endpoint:             cfg.Events.Addr,
		Cluster:              cfg.Events.ClusterID,
		EnableTLS:            cfg.Events.EnableTLS,
		TLSInsecure:          cfg.Events.TLSInsecure,
		TLSRootCACertificate: cfg.Events.TLSRootCaCertPath,
	})
}
