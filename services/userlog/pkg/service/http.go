package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
)

// ServeHTTP fulfills Handler interface
func (ul *UserlogService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ul.m.ServeHTTP(w, r)
}

// HandleGetEvents is the GET handler for events
func (ul *UserlogService) HandleGetEvents(w http.ResponseWriter, r *http.Request) {
	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		ul.log.Error().Int("returned statuscode", http.StatusUnauthorized).Msg("user unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	evs, err := ul.GetEvents(r.Context(), u.GetId().GetOpaqueId())
	if err != nil {
		ul.log.Error().Err(err).Int("returned statuscode", http.StatusInternalServerError).Msg("get events failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := GetEventResponseOC10{}
	for _, e := range evs {
		noti, err := ul.convertEvent(r.Context(), e)
		if err != nil {
			ul.log.Error().Err(err).Str("eventid", e.Id).Str("eventtype", e.Type).Msg("failed to convert event")
			continue
		}

		resp.OCS.Data = append(resp.OCS.Data, noti)
	}

	resp.OCS.Meta.StatusCode = http.StatusOK
	b, _ := json.Marshal(resp)
	w.Write(b)
}

// HandleDeleteEvents is the DELETE handler for events
func (ul *UserlogService) HandleDeleteEvents(w http.ResponseWriter, r *http.Request) {
	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		ul.log.Error().Int("returned statuscode", http.StatusUnauthorized).Msg("user unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var ids []string
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		ul.log.Error().Err(err).Int("returned statuscode", http.StatusBadRequest).Msg("request body is malformed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ul.DeleteEvents(u.GetId().GetOpaqueId(), ids); err != nil {
		ul.log.Error().Err(err).Int("returned statuscode", http.StatusInternalServerError).Msg("delete events failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ul *UserlogService) convertEvent(ctx context.Context, event *ehmsg.Event) (OC10Notification, error) {
	etype, ok := ul.registeredEvents[event.Type]
	if !ok {
		// this should not happen
		return OC10Notification{}, errors.New("eventtype not registered")
	}

	einterface, err := etype.Unmarshal(event.Event)
	if err != nil {
		// this shouldn't happen either
		return OC10Notification{}, errors.New("cant unmarshal event")
	}

	noti := OC10Notification{
		EventID:   event.Id,
		Service:   "userlog",
		Timestamp: time.Now().Format(time.RFC3339Nano),
	}

	// TODO: strange bug with getting space -> fix postponed to make master panic-free
	var space storageprovider.StorageSpace
	switch ev := einterface.(type) {
	// file related
	case events.UploadReady:
		noti.UserID = ev.ExecutingUser.GetId().GetOpaqueId()
		noti.Subject = "File uploaded"
		noti.Message = fmt.Sprintf("File '%s' was uploaded to space '%s' by user '%s'", ev.Filename, space.GetName(), ev.ExecutingUser.GetUsername())
	case events.ContainerCreated:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Folder created"
		noti.Message = fmt.Sprintf("Folder '%s' was created", ev.Ref.GetPath())
	case events.FileTouched:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "File touched"
		noti.Message = fmt.Sprintf("File '%s' was touched", ev.Ref.GetPath())
	case events.FileDownloaded:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "File downloaded"
		noti.Message = fmt.Sprintf("File '%s' was downloaded", ev.Ref.GetPath())
	case events.FileVersionRestored:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "File version restored"
		noti.Message = fmt.Sprintf("An older version of file '%s' was restored", ev.Ref.GetPath())
	case events.ItemMoved:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "File moved"
		noti.Message = fmt.Sprintf("File '%s' was moved from '%s'", ev.Ref.GetPath(), ev.OldReference.GetPath())
	case events.ItemTrashed:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "File trashed"
		noti.Message = fmt.Sprintf("File '%s' was trashed", ev.Ref.GetPath())
	case events.ItemPurged:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "File purged"
		noti.Message = fmt.Sprintf("File '%s' was purged", ev.Ref.GetPath())
	case events.ItemRestored:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "File restored"
		noti.Message = fmt.Sprintf("File '%s' was restored", ev.Ref.GetPath())

	// space related
	case events.SpaceCreated:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Space created"
		noti.Message = fmt.Sprintf("Space '%s' was created", ev.Name)
	case events.SpaceRenamed:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Space renamed"
		noti.Message = fmt.Sprintf("Space '%s' was renamed", ev.Name)
	case events.SpaceEnabled:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Space enabled"
		noti.Message = fmt.Sprintf("Space '%s' was renamed", space.Name)
	case events.SpaceDisabled:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Space disabled"
		noti.Message = fmt.Sprintf("Space '%s' was disabled", space.Name)
	case events.SpaceDeleted:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Space deleted"
		noti.Message = fmt.Sprintf("Space '%s' was deleted", space.Name)
	case events.SpaceShared:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Space shared"
		noti.Message = fmt.Sprintf("Space '%s' was shared", space.Name)
	case events.SpaceUnshared:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Space unshared"
		noti.Message = fmt.Sprintf("Space '%s' was unshared", space.Name)
	case events.SpaceUpdated:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Space updated"
		noti.Message = fmt.Sprintf("Space '%s' was updated", space.Name)
	case events.SpaceMembershipExpired:
		noti.UserID = ""
		noti.Subject = "Space membership expired"
		noti.Message = fmt.Sprintf("A spacemembership for space '%s' has expired", space.Name)

	// share related
	case events.ShareCreated:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Share received"
		noti.Message = fmt.Sprintf("A file was shared in space %s", space.Name)
	case events.ShareUpdated:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Share updated"
		noti.Message = fmt.Sprintf("A share was updated in space %s", space.Name)
	case events.ShareExpired:
		noti.Subject = "Share expired"
		noti.Message = fmt.Sprintf("A share has expired in space %s", space.Name)
	case events.LinkCreated:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Share received"
		noti.Message = fmt.Sprintf("A link was created in space %s", space.Name)
	case events.LinkUpdated:
		noti.UserID = ev.Executant.GetOpaqueId()
		noti.Subject = "Share received"
		noti.Message = fmt.Sprintf("A link was updated in space %s", space.Name)
	}

	return noti, nil
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
