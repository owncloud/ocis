package service

import (
	"context"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/eventhistory/pkg/config"
	"go-micro.dev/v4/store"
)

// EventHistoryService is the service responsible for event history
type EventHistoryService struct {
	ch    <-chan events.Event
	store store.Store
	cfg   *config.Config
	log   log.Logger
}

// NewEventHistoryService returns an EventHistory service
func NewEventHistoryService(cfg *config.Config, consumer events.Consumer, store store.Store, log log.Logger) (*EventHistoryService, error) {
	if consumer == nil || store == nil {
		return nil, fmt.Errorf("Need non nil consumer (%v) and store (%v) to work properly", consumer, store)
	}

	ch, err := events.ConsumeAll(consumer, "evhistory")
	if err != nil {
		return nil, err
	}

	eh := &EventHistoryService{ch: ch, store: store, cfg: cfg, log: log}
	go eh.StoreEvents()

	return eh, nil
}

// StoreEvents consumes all events and stores them in the store. Will block
func (eh *EventHistoryService) StoreEvents() {
	for event := range eh.ch {
		if err := eh.store.Write(&store.Record{
			Key:    event.ID,
			Value:  event.Event.([]byte),
			Expiry: eh.cfg.Store.RecordExpiry,
			Metadata: map[string]interface{}{
				"type": event.Type,
			},
		}); err != nil {
			// we can't store. That's it for us.
			eh.log.Error().Err(err).Str("eventid", event.ID).Msg("could not store event")
			return
		}
	}
}

// GetEvents allows to retrieve events from the eventstore by id
func (eh *EventHistoryService) GetEvents(ctx context.Context, req *ehsvc.GetEventsRequest, resp *ehsvc.GetEventsResponse) error {
	for _, id := range req.Ids {
		evs, err := eh.store.Read(id)
		if err != nil {
			if err != store.ErrNotFound {
				eh.log.Error().Err(err).Str("eventid", id).Msg("could not read event")
			}
			continue
		}

		resp.Events = append(resp.Events, &ehmsg.Event{
			Id:    id,
			Event: evs[0].Value,
			Type:  evs[0].Metadata["type"].(string),
		})
	}

	return nil
}
