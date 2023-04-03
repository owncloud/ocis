package service

import (
	"context"
	"fmt"
	"regexp"

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

// GetEventsForUser allows to retrieve events from the eventstore by userID
// This function will match all events that contains the user ID between two non-word characters.
// The reasoning behind this is that events put the userID in many different fields, which can differ
// per event type. This function will match all events that contain the userID by using a regex.
// This should also cover future events that might contain the userID in a different field.
func (eh *EventHistoryService) GetEventsForUser(ctx context.Context, req *ehsvc.GetEventsForUserRequest, resp *ehsvc.GetEventsResponse) error {
	idx, err := eh.store.List(store.ListPrefix(""))
	if err != nil {
		eh.log.Error().Err(err).Msg("could not list events")
		return err
	}

	// Match all events that contains the user ID between two non-word characters.
	userID, err := regexp.Compile(fmt.Sprintf(`\W%s\W`, req.UserID))
	if err != nil {
		eh.log.Error().Err(err).Str("userID", req.UserID).Msg("could not compile regex")
		return err
	}

	for _, i := range idx {
		e, err := eh.store.Read(i)
		if err != nil {
			eh.log.Error().Err(err).Str("eventid", i).Msg("could not read event")
			continue
		}

		if userID.Match(e[0].Value) {
			resp.Events = append(resp.Events, &ehmsg.Event{
				Id:    i,
				Event: e[0].Value,
				Type:  e[0].Metadata["type"].(string),
			})
		}
	}

	return nil
}
