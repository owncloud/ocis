package registry

import (
	"context"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	mRegistry "go-micro.dev/v4/registry"
)

// RegisterService publishes an arbitrary endpoint to the service-registry. This allows querying nodes of
// non-micro services like reva. No health-checks are done, thus the caller is responsible for canceling.
func RegisterService(ctx context.Context, service *mRegistry.Service, logger log.Logger) error {
	registry := GetRegistry()
	node := service.Nodes[0]

	logger.Info().Msgf("registering external service %v@%v", node.Id, node.Address)

	rOpts := []mRegistry.RegisterOption{mRegistry.RegisterTTL(GetRegisterTTL())}
	if err := registry.Register(service, rOpts...); err != nil {
		logger.Fatal().Err(err).Msgf("Registration error for external service %v", service.Name)
	}

	t := time.NewTicker(GetRegisterInterval())

	go func() {
		for {
			select {
			case <-t.C:
				logger.Debug().Interface("service", service).Msg("refreshing external service-registration")
				err := registry.Register(service, rOpts...)
				if err != nil {
					logger.Error().Err(err).Msgf("registration error for external service %v", service.Name)
				}
			case <-ctx.Done():
				logger.Debug().Interface("service", service).Msg("unregistering")
				t.Stop()
				err := registry.Deregister(service)
				if err != nil {
					logger.Err(err).Msgf("Error unregistering external service %v", service.Name)
				}
				return
			}
		}
	}()

	return nil
}
