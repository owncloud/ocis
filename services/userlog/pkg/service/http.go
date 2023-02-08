package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
)

// BuildMux builds the and configures the muxer
func (ul *UserlogService) BuildMux() {
	m := chi.NewMux()
	m.Use(middleware.ExtractAccountUUID())

	m.Route("/", func(r chi.Router) {
		r.Get("/*", ul.HandleGetEvents)
		r.Delete("/*", ul.HandleDeleteEvents)
	})

	ul.m = m
}

// ServeHTTP fulfills Handler interface
func (ul *UserlogService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ul.m.ServeHTTP(w, r)
}

// HandleGetEvents is the GET handler for events
func (ul *UserlogService) HandleGetEvents(w http.ResponseWriter, r *http.Request) {
	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID := u.GetId().GetOpaqueId()

	evs, err := ul.GetEvents(r.Context(), userID)
	if err != nil {
		return
	}

	resp := GetEventResponseOC10{}
	for _, e := range evs {
		resp.OCS.Data = append(resp.OCS.Data, ul.convertEvent(e))
	}

	resp.OCS.Meta.StatusCode = http.StatusOK
	b, _ := json.Marshal(resp)
	w.Write(b)
}

// HandleDeleteEvents is the DELETE handler for events
func (ul *UserlogService) HandleDeleteEvents(w http.ResponseWriter, r *http.Request) {
	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ul.DeleteEvents(u.GetId().GetOpaqueId(), ids); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ul *UserlogService) convertEvent(event *ehmsg.Event) OC10Notification {
	etype, ok := ul.registeredEvents[event.Type]
	if !ok {
		// this should not happen
		return OC10Notification{}
	}

	einterface, err := etype.Unmarshal(event.Event)
	if err != nil {
		// this shouldn't happen either
		return OC10Notification{}
	}

	noti := OC10Notification{
		EventID:   event.Id,
		Service:   "userlog",
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}

	switch ev := einterface.(type) {
	case events.UploadReady:
		ctx, _, _ := utils.Impersonate(ev.SpaceOwner, ul.gwClient, ul.cfg.MachineAuthAPIKey)
		space, _ := ul.getSpace(ctx, ev.FileRef.GetResourceId().GetSpaceId())
		noti.UserID = ev.ExecutingUser.GetId().GetOpaqueId()
		noti.Subject = "File uploaded"
		noti.Message = fmt.Sprintf("File %s was uploaded to space %s by user %s", ev.Filename, space.GetName(), ev.ExecutingUser.GetUsername())
	}

	return noti
}

// OC10Notification is the oc10 style representation of an event
// some fields are left out for simplicity
type OC10Notification struct {
	EventID   string `json:"notification_id"`
	Service   string `json:"app"`
	Timestamp string `json:"datetime"`
	UserID    string `json:"user"`
	Subject   string `json:"subject"`
	Message   string `json:"message"`
}

// GetEventResponseOC10 is the response from GET events endpoint in oc10 style
type GetEventResponseOC10 struct {
	OCS struct {
		Meta struct {
			Message    string `json:"message"`
			Status     string `json:"status"`
			StatusCode int    `json:"statuscode"`
		} `json:"meta"`
		Data []OC10Notification `json:"data"`
	} `json:"ocs"`
}
