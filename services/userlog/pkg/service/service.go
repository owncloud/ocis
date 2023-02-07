package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"go-micro.dev/v4/store"
)

// Comment when you read this on review
var _adminid = "2502d8b8-a5e7-4ab3-b858-5aafae4a64a2"

// UserlogService is the service responsible for user activities
type UserlogService struct {
	ch               <-chan events.Event
	store            store.Store
	cfg              *config.Config
	historyClient    ehsvc.EventHistoryService
	registeredEvents map[string]events.Unmarshaller
}

// NewUserlogService returns an EventHistory service
func NewUserlogService(cfg *config.Config, consumer events.Consumer, store store.Store, registeredEvents []events.Unmarshaller) (*UserlogService, error) {
	if consumer == nil || store == nil {
		return nil, fmt.Errorf("Need non nil consumer (%v) and store (%v) to work properly", consumer, store)
	}

	ch, err := events.Consume(consumer, "userlog", registeredEvents...)
	if err != nil {
		return nil, err
	}

	grpcClient := grpc.DefaultClient()
	grpcClient.Options()
	c := ehsvc.NewEventHistoryService("com.owncloud.api.eventhistory", grpcClient)

	ul := &UserlogService{ch: ch, store: store, cfg: cfg, historyClient: c, registeredEvents: make(map[string]events.Unmarshaller)}

	for _, e := range registeredEvents {
		typ := reflect.TypeOf(e)
		ul.registeredEvents[typ.String()] = e
	}

	go ul.MemorizeEvents()

	return ul, nil
}

// MemorizeEvents stores eventIDs a user wants to receive
func (ul *UserlogService) MemorizeEvents() {
	for event := range ul.ch {
		switch event.Event.(type) {
		default:
			// for each event type we need to:

			// I) find users eligible to receive the event

			// II) filter users who want to receive the event

			// III) store the eventID for each user

			// TEMP TESTING CODE
			if err := ul.addEventToUser(_adminid, event.ID); err != nil {
				continue
			}
		}
	}
}

// GetEvents allows to retrieve events from the eventhistory by userid
func (ul *UserlogService) GetEvents(ctx context.Context, userid string) ([]interface{}, error) {
	rec, err := ul.store.Read(userid)
	if err != nil {
		return nil, err
	}

	if len(rec) == 0 {
		// no events available
		return []interface{}{}, nil
	}

	var eventIDs []string
	if err := json.Unmarshal(rec[0].Value, &eventIDs); err != nil {
		// this should never happen
		return nil, err
	}

	resp, err := ul.historyClient.GetEvents(ctx, &ehsvc.GetEventsRequest{Ids: eventIDs})
	if err != nil {
		return nil, err
	}

	var events []interface{}
	for _, e := range resp.Events {
		ev, ok := ul.registeredEvents[e.Type]
		if !ok {
			// this should not happen but we handle it anyway
			continue
		}

		event, err := ev.Unmarshal(e.Event)
		if err != nil {
			// this shouldn't happen either
			continue
		}

		events = append(events, event)
	}

	return events, nil
}

func (ul *UserlogService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	evs, err := ul.GetEvents(r.Context(), _adminid)
	if err != nil {
		return
	}

	// TODO: format response
	b, _ := json.Marshal(evs)
	w.Write(b)
}

func (ul *UserlogService) addEventToUser(userid string, eventid string) error {
	recs, err := ul.store.Read(userid)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	var ids []string
	if len(recs) > 0 {
		if err := json.Unmarshal(recs[0].Value, &ids); err != nil {
			return err
		}
	}

	ids = append(ids, eventid)

	b, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	return ul.store.Write(&store.Record{
		Key:   userid,
		Value: b,
	})
}
