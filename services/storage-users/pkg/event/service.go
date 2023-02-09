package event

import (
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/task"
)

const (
	consumerGroup = "storage-users"
)

// Service wraps all common logic that is needed to react to incoming events.
type Service struct {
	gateway      gateway.GatewayAPIClient
	stream       events.Stream
	logger       log.Logger
	clientSecret string
}

// NewService prepares and returns a Service implementation.
func NewService(gw gateway.GatewayAPIClient, s events.Stream, l log.Logger, clientSecret string) (Service, error) {
	svc := Service{
		gateway:      gw,
		stream:       s,
		logger:       l,
		clientSecret: clientSecret,
	}

	return svc, nil
}

// Run to fulfil Runner interface
func (s Service) Run() error {
	ch, err := events.Consume(s.stream, consumerGroup, PurgeTrashBin{})
	if err != nil {
		return err
	}

	for e := range ch {
		var err error

		switch ev := e.(type) {
		case PurgeTrashBin:
			err = task.PurgeTrashBin(s.gateway, s.clientSecret, ev.ExecutantID, ev.RemoveBefore)
		}

		if err != nil {
			s.logger.Error().Err(err).Interface("event", e)
		}
	}

	return nil
}
