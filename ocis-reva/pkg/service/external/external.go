package external

import (
	"context"
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/v2/registry"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// RegisterGRPCEndpoint publishes an arbitrary endpoint to the service-registry. This allows to query nodes of
// non-micro GRPC-services like reva. No health-checks are done, thus the caller is responsible for canceling.
//
func RegisterGRPCEndpoint(ctx context.Context, serviceID string, uuid string, addr string, logger log.Logger) error {
	node := &registry.Node{
		Id:       serviceID + "-" + uuid,
		Address:  addr,
		Metadata: make(map[string]string),
	}
	node.Metadata["broker"] = broker.String()
	node.Metadata["registry"] = registry.String()
	node.Metadata["server"] = "grpc"
	node.Metadata["transport"] = "grpc"
	node.Metadata["protocol"] = "grpc"

	service := &registry.Service{
		Name:      serviceID,
		Version:   "",
		Nodes:     []*registry.Node{node},
		Endpoints: make([]*registry.Endpoint, 0),
	}

	rOpts := []registry.RegisterOption{registry.RegisterTTL(time.Minute)}

	logger.Info().Msgf("Registering external service %v@%v", node.Id, node.Address)

	if err := registry.Register(service, rOpts...); err != nil {
		logger.Fatal().Err(err).Msgf("Registration error for external service %v", serviceID)
	}

	t := time.NewTicker(time.Second * 30)

	go func() {
		for {
			select {
			case <-t.C:
				logger.Debug().Interface("service", service).Msg("Refreshing external service-registration")
				err := registry.Register(service, rOpts...)
				if err != nil {
					logger.Error().Err(err).Msgf("Registration error for external service %v", serviceID)
				}
			case <-ctx.Done():
				logger.Debug().Interface("service", service).Msg("Unregistering")
				t.Stop()
				err := registry.Deregister(service)
				if err != nil {
					logger.Err(err).Msgf("Error unregistering external service %v", serviceID)
				}

			}
		}
	}()

	return nil
}
