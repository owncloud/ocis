package service

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/ocis-pkg/messaging"
)

//go:embed l10n/locale
var _translationFS embed.FS

// HeaderAcceptLanguage is the header where the client can set the locale
var HeaderAcceptLanguage = "Accept-Language"

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

	filesystem, _ := fs.Sub(_translationFS, "l10n/locale")
	if ul.cfg.TranslationPath != "" {
		filesystem = os.DirFS(ul.cfg.TranslationPath)
	}

	conv := messaging.NewConverter(r.Header.Get(HeaderAcceptLanguage), ul.gwClient, ul.cfg.MachineAuthAPIKey, ul.cfg.Service.Name, filesystem)

	resp := GetEventResponseOC10{}
	for _, e := range evs {
		etype, ok := ul.registeredEvents[e.Type]
		if !ok {
			// this should not happen
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ev, err := etype.Unmarshal(e.Event)
		if err != nil {
			// this shouldn't happen either
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		noti, err := conv.ConvertEvent(e.Id, ev)
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

	var req DeleteEventsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ul.log.Error().Err(err).Int("returned statuscode", http.StatusBadRequest).Msg("request body is malformed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ul.DeleteEvents(u.GetId().GetOpaqueId(), req.IDs); err != nil {
		ul.log.Error().Err(err).Int("returned statuscode", http.StatusInternalServerError).Msg("delete events failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetEventResponseOC10 is the response from GET events endpoint in oc10 style
type GetEventResponseOC10 struct {
	OCS struct {
		Meta struct {
			Message    string `json:"message"`
			Status     string `json:"status"`
			StatusCode int    `json:"statuscode"`
		} `json:"meta"`
		Data []messaging.OC10Notification `json:"data"`
	} `json:"ocs"`
}

// DeleteEventsRequest is the expected body for the delete request
type DeleteEventsRequest struct {
	IDs []string `json:"ids"`
}
