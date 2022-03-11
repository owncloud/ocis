package types

import (
	"github.com/cs3org/reva/v2/pkg/events"
)

// RegisteredEvents returns the events the service is registered for
func RegisteredEvents() []events.Unmarshaller {
	return []events.Unmarshaller{
		events.ShareCreated{},
		events.ShareUpdated{},
		events.LinkCreated{},
		events.LinkUpdated{},
	}
}
