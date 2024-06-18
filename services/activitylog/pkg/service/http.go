package service

import (
	"embed"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
	"github.com/owncloud/ocis/v2/services/search/pkg/query/kql"
)

var (
	//go:embed l10n/locale
	_localeFS embed.FS

	// subfolder where the translation files are stored
	_localeSubPath = "l10n/locale"

	// domain of the activitylog service (transifex)
	_domain = "activitylog"
)

// ServeHTTP implements the http.Handler interface.
func (s *ActivitylogService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// HandleGetItemActivities handles the request to get the activities of an item.
func (s *ActivitylogService) HandleGetItemActivities(w http.ResponseWriter, r *http.Request) {
	activeUser, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	qraw := r.URL.Query().Get("kql")
	if qraw == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	qBuilder := kql.Builder{}
	qast, err := qBuilder.Build(qraw)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var itemID string

	for _, n := range qast.Nodes {
		v, ok := n.(*ast.StringNode)
		if !ok {
			continue
		}

		if strings.ToLower(v.Key) != "itemid" {
			continue
		}

		itemID = v.Value
	}

	rid, err := storagespace.ParseID(itemID)
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

	evRes, err := s.evHistory.GetEvents(r.Context(), &ehsvc.GetEventsRequest{Ids: ids})
	if err != nil {
		s.log.Error().Err(err).Msg("error getting events")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var resp GetActivitiesResponse
	for _, e := range evRes.GetEvents() {
		// TODO: compare returned events with initial list and remove missing ones

		// FIXME: Should all users get all events? If not we can filter here

		var (
			message string
			res     Resource
			act     Actor
			ts      libregraph.ActivityTimes
		)

		switch ev := s.unwrapEvent(e).(type) {
		case nil:
			// error already logged in unwrapEvent
			continue
		case events.UploadReady:
			message = MessageResourceCreated
			res, act, ts, err = s.ResponseData(ev.FileRef, ev.ExecutingUser.GetId(), ev.ExecutingUser.GetDisplayName(), utils.TSToTime(ev.Timestamp))
		case events.FileTouched:
			message = MessageResourceCreated
			res, act, ts, err = s.ResponseData(ev.Ref, ev.Executant, "", utils.TSToTime(ev.Timestamp))
		case events.ContainerCreated:
			message = MessageResourceCreated
			res, act, ts, err = s.ResponseData(ev.Ref, ev.Executant, "", utils.TSToTime(ev.Timestamp))
		case events.ItemTrashed:
			message = MessageResourceTrashed
			res, act, ts, err = s.ResponseData(ev.Ref, ev.Executant, "", utils.TSToTime(ev.Timestamp))
		case events.ItemPurged:
			message = MessageResourcePurged
			res, act, ts, err = s.ResponseData(ev.Ref, ev.Executant, "", utils.TSToTime(ev.Timestamp))
		case events.ItemMoved:
			message = MessageResourceMoved
			res, act, ts, err = s.ResponseData(ev.Ref, ev.Executant, "", utils.TSToTime(ev.Timestamp))
		case events.ShareCreated:
			message = MessageShareCreated
			res, act, ts, err = s.ResponseData(toRef(ev.ItemID), ev.Executant, "", utils.TSToTime(ev.CTime))
		case events.ShareUpdated:
			message = MessageShareUpdated
			res, act, ts, err = s.ResponseData(toRef(ev.ItemID), ev.Executant, "", utils.TSToTime(ev.MTime))
		case events.ShareRemoved:
			message = MessageShareDeleted
			res, act, ts, err = s.ResponseData(toRef(ev.ItemID), ev.Executant, "", ev.Timestamp)
		case events.LinkCreated:
			message = MessageLinkCreated
			res, act, ts, err = s.ResponseData(toRef(ev.ItemID), ev.Executant, "", utils.TSToTime(ev.CTime))
		case events.LinkUpdated:
			message = MessageLinkUpdated
			res, act, ts, err = s.ResponseData(toRef(ev.ItemID), ev.Executant, "", utils.TSToTime(ev.CTime))
		case events.LinkRemoved:
			message = MessageLinkDeleted
			res, act, ts, err = s.ResponseData(toRef(ev.ItemID), ev.Executant, "", utils.TSToTime(ev.Timestamp))
		case events.SpaceShared:
			message = MessageSpaceShared
			res, act, ts, err = s.ResponseData(sToRef(ev.ID), ev.Executant, "", ev.Timestamp)
		case events.SpaceShareUpdated:
			message = MessageSpaceShareUpdated
			res, act, ts, err = s.ResponseData(sToRef(ev.ID), ev.Executant, "", ev.Timestamp)
		case events.SpaceUnshared:
			message = MessageSpaceUnshared
			res, act, ts, err = s.ResponseData(sToRef(ev.ID), ev.Executant, "", ev.Timestamp)
		}

		if err != nil {
			s.log.Error().Err(err).Msg("error getting response data")
			continue
		}

		// todo: configurable default locale?
		loc := l10n.MustGetUserLocale(r.Context(), activeUser.GetId().GetOpaqueId(), r.Header.Get(l10n.HeaderAcceptLanguage), s.valService)
		t := l10n.NewTranslatorFromCommonConfig("en", _domain, "", _localeFS, _localeSubPath)

		resp.Activities = append(resp.Activities, NewActivity(t.Translate(message, loc), res, act, ts, e.GetId()))
	}

	b, err := json.Marshal(resp)
	if err != nil {
		s.log.Error().Err(err).Msg("error marshalling activities")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		s.log.Error().Err(err).Msg("error writing response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
