package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/go-chi/chi/v5"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
)

// ServeHTTP implements the http.Handler interface.
func (s *ActivitylogService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// HandleGetItemActivities handles the request to get the activities of an item.
func (s *ActivitylogService) HandleGetItemActivities(w http.ResponseWriter, r *http.Request) {
	// TODO: Compare driveid with itemid to avoid bad requests
	rid, err := parseIDParam(r, "item-id")
	if err != nil {
		s.log.Info().Err(err).Msg("invalid resource id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	raw, err := s.Activities(&rid)
	if err != nil {
		s.log.Error().Err(err).Msg("error getting activities")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ids := make([]string, 0, len(raw))
	for _, a := range raw {
		// TODO: Filter by depth and timestamp
		ids = append(ids, a.EventID)
	}

	fmt.Println("IDS:", ids)

	evRes, err := s.evHistory.GetEvents(r.Context(), &ehsvc.GetEventsRequest{Ids: ids})
	if err != nil {
		s.log.Error().Err(err).Msg("error getting events")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: compare returned events with initial list and remove missing ones

	fmt.Println("EVENTS:", evRes.GetEvents())

	var acts []Activity
	for _, e := range evRes.GetEvents() {
		// FIXME: Should all users get all events? If not we can filter here

		switch ev := s.unwrapEvent(e).(type) {
		case nil:
			// error already logged in unwrapEvent
			continue
		case events.UploadReady:
			act := UploadReady(e.Id, ev)
			acts = append(acts, act)
		}
	}

	fmt.Println("ACTIVITIES:", acts)

	b, err := json.Marshal(acts)
	if err != nil {
		s.log.Error().Err(err).Msg("error marshalling activities")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

func (s *ActivitylogService) unwrapEvent(e *ehmsg.Event) interface{} {
	etype, ok := s.registeredEvents[e.GetType()]
	if !ok {
		s.log.Error().Str("eventid", e.GetId()).Str("eventtype", e.GetType()).Msg("event not registered")
		return nil
	}

	einterface, err := etype.Unmarshal(e.GetEvent())
	if err != nil {
		s.log.Error().Str("eventid", e.GetId()).Str("eventtype", e.GetType()).Msg("failed to umarshal event")
		return nil
	}

	return einterface
}

// TODO: I found this on graph service. We should move it to `utils` pkg so both services can use it.
func parseIDParam(r *http.Request, param string) (provider.ResourceId, error) {
	driveID, err := url.PathUnescape(chi.URLParam(r, param))
	if err != nil {
		return provider.ResourceId{}, errorcode.New(errorcode.InvalidRequest, err.Error())
	}

	id, err := storagespace.ParseID(driveID)
	if err != nil {
		return provider.ResourceId{}, errorcode.New(errorcode.InvalidRequest, err.Error())
	}
	return id, nil
}
