package service

import (
	"encoding/json"
	"net/http"

	"github.com/cs3org/reva/v2/pkg/ctx"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	settings "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
)

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

	conv := ul.getConverter(r.Header.Get(HeaderAcceptLanguage))

	resp := GetEventResponseOC10{}
	for _, e := range evs {
		etype, ok := ul.registeredEvents[e.Type]
		if !ok {
			ul.log.Error().Str("eventid", e.Id).Str("eventtype", e.Type).Msg("event not registered")
			continue
		}

		einterface, err := etype.Unmarshal(e.Event)
		if err != nil {
			ul.log.Error().Str("eventid", e.Id).Str("eventtype", e.Type).Msg("failed to umarshal event")
			continue
		}

		noti, err := conv.ConvertEvent(e.Id, einterface)
		if err != nil {
			ul.log.Error().Err(err).Str("eventid", e.Id).Str("eventtype", e.Type).Msg("failed to convert event")
			continue
		}

		resp.OCS.Data = append(resp.OCS.Data, noti)
	}

	glevs, err := ul.GetGlobalEvents()
	if err != nil {
		ul.log.Error().Err(err).Int("returned statuscode", http.StatusInternalServerError).Msg("get global events failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for t, data := range glevs {
		noti, err := conv.ConvertGlobalEvent(t, data)
		if err != nil {
			ul.log.Error().Err(err).Str("eventtype", t).Msg("failed to convert event")
			continue
		}

		resp.OCS.Data = append(resp.OCS.Data, noti)
	}

	resp.OCS.Meta.StatusCode = http.StatusOK
	b, _ := json.Marshal(resp)
	w.Write(b)
}

// HandleSSE is the GET handler for events
func (ul *UserlogService) HandleSSE(w http.ResponseWriter, r *http.Request) {
	u, ok := ctx.ContextGetUser(r.Context())
	if !ok {
		ul.log.Error().Msg("sse: no user in context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uid := u.GetId().GetOpaqueId()
	if uid == "" {
		ul.log.Error().Msg("sse: user in context is broken")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stream := ul.sse.CreateStream(uid)
	stream.AutoReplay = false

	// add stream to URL
	q := r.URL.Query()
	q.Set("stream", uid)
	r.URL.RawQuery = q.Encode()

	ul.sse.ServeHTTP(w, r)
}

// HandlePostEvent is the POST handler for events
func (ul *UserlogService) HandlePostEvent(w http.ResponseWriter, r *http.Request) {
	u, ok := ctx.ContextGetUser(r.Context())
	if !ok {
		ul.log.Error().Msg("post: no user in context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uid := u.GetId().GetOpaqueId()
	if uid == "" {
		ul.log.Error().Msg("post: user in context is broken")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req PostEventsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ul.log.Error().Err(err).Int("returned statuscode", http.StatusBadRequest).Msg("request body is malformed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ul.StoreGlobalEvent(req.Type, req.Data); err != nil {
		ul.log.Error().Err(err).Msg("post: error storing global event")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
		Data []OC10Notification `json:"data"`
	} `json:"ocs"`
}

// DeleteEventsRequest is the expected body for the delete request
type DeleteEventsRequest struct {
	IDs []string `json:"ids"`
}

// PostEventsRequest is the expected body for the post request
type PostEventsRequest struct {
	// the event type, e.g. "deprovision"
	Type string `json:"type"`
	// arbitray data for the event
	Data map[string]string `json:"data"`
}

// RequireAdmin middleware is used to require the user in context to be an admin / have account management permissions
func RequireAdmin(rm *roles.Manager, logger log.Logger) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			u, ok := revactx.ContextGetUser(r.Context())
			if !ok {
				errorcode.AccessDenied.Render(w, r, http.StatusUnauthorized, "Unauthorized")
				return
			}
			if u.Id == nil || u.Id.OpaqueId == "" {
				errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "user is missing an id")
				return
			}
			// get roles from context
			roleIDs, ok := roles.ReadRoleIDsFromContext(r.Context())
			if !ok {
				logger.Debug().Str("userid", u.Id.OpaqueId).Msg("No roles in context, contacting settings service")
				var err error
				roleIDs, err = rm.FindRoleIDsForUser(r.Context(), u.Id.OpaqueId)
				if err != nil {
					logger.Err(err).Str("userid", u.Id.OpaqueId).Msg("failed to get roles for user")
					errorcode.AccessDenied.Render(w, r, http.StatusUnauthorized, "Unauthorized")
					return
				}
				if len(roleIDs) == 0 {
					errorcode.AccessDenied.Render(w, r, http.StatusUnauthorized, "Unauthorized")
					return
				}
			}

			// check if permission is present in roles of the authenticated account
			if rm.FindPermissionByID(r.Context(), roleIDs, settings.AccountManagementPermissionID) != nil {
				next.ServeHTTP(w, r)
				return
			}

			errorcode.AccessDenied.Render(w, r, http.StatusForbidden, "Forbidden")
		}
	}
}
