package service

import (
	"embed"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"google.golang.org/grpc/metadata"

	"github.com/owncloud/ocis/v2/ocis-pkg/ast"
	"github.com/owncloud/ocis/v2/ocis-pkg/kql"
	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
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
	ctx := r.Context()
	ctx = metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, r.Header.Get("X-Access-Token"))

	activeUser, ok := revactx.ContextGetUser(ctx)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rid, limit, rawActivityAccepted, activityAccepted, err := s.getFilters(r.URL.Query().Get("kql"))
	if err != nil {
		s.log.Info().Str("query", r.URL.Query().Get("kql")).Err(err).Msg("error getting filters")
		_, _ = w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	raw, err := s.Activities(rid)
	if err != nil {
		s.log.Error().Err(err).Msg("error getting activities")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ids := make([]string, 0, len(raw))
	toDelete := make(map[string]struct{}, len(raw))
	for _, a := range raw {
		if !rawActivityAccepted(a) {
			continue
		}
		ids = append(ids, a.EventID)
		toDelete[a.EventID] = struct{}{}
	}

	evRes, err := s.evHistory.GetEvents(r.Context(), &ehsvc.GetEventsRequest{Ids: ids})
	if err != nil {
		s.log.Error().Err(err).Msg("error getting events")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var resp GetActivitiesResponse
	for _, e := range evRes.GetEvents() {
		delete(toDelete, e.GetId())

		if limit != 0 && len(resp.Activities) >= limit {
			continue
		}

		if !activityAccepted(e) {
			continue
		}

		var (
			message string
			ts      time.Time
			vars    map[string]interface{}
		)

		switch ev := s.unwrapEvent(e).(type) {
		case nil:
			// error already logged in unwrapEvent
			continue
		case events.UploadReady:
			message = MessageResourceCreated
			ts = utils.TSToTime(ev.Timestamp)
			vars, err = s.GetVars(ctx, WithResource(ev.FileRef, true), WithUser(ev.ExecutingUser.GetId(), ev.ExecutingUser.GetDisplayName()))
		case events.FileTouched:
			message = MessageResourceCreated
			ts = utils.TSToTime(ev.Timestamp)
			vars, err = s.GetVars(ctx, WithResource(ev.Ref, true), WithUser(ev.Executant, ""))
		case events.ContainerCreated:
			message = MessageResourceCreated
			ts = utils.TSToTime(ev.Timestamp)
			vars, err = s.GetVars(ctx, WithResource(ev.Ref, true), WithUser(ev.Executant, ""))
		case events.ItemTrashed:
			message = MessageResourceTrashed
			ts = utils.TSToTime(ev.Timestamp)
			vars, err = s.GetVars(ctx, WithTrashedResource(ev.Ref, ev.ID), WithUser(ev.Executant, ""), WithSpace(toSpace(ev.Ref)))
		case events.ItemMoved:
			switch isRename(ev.OldReference, ev.Ref) {
			case true:
				message = MessageResourceRenamed
				vars, err = s.GetVars(ctx, WithResource(ev.Ref, false), WithOldResource(ev.OldReference), WithUser(ev.Executant, ""))
			case false:
				message = MessageResourceMoved
				vars, err = s.GetVars(ctx, WithResource(ev.Ref, true), WithUser(ev.Executant, ""))
			}
			ts = utils.TSToTime(ev.Timestamp)
		case events.ShareCreated:
			message = MessageShareCreated
			ts = utils.TSToTime(ev.CTime)
			vars, err = s.GetVars(ctx, WithResource(toRef(ev.ItemID), false), WithUser(ev.Executant, ""), WithSharee(ev.GranteeUserID, ev.GranteeGroupID))
		case events.ShareRemoved:
			message = MessageShareDeleted
			ts = ev.Timestamp
			vars, err = s.GetVars(ctx, WithResource(toRef(ev.ItemID), false), WithUser(ev.Executant, ""), WithSharee(ev.GranteeUserID, ev.GranteeGroupID))
		case events.LinkCreated:
			message = MessageLinkCreated
			ts = utils.TSToTime(ev.CTime)
			vars, err = s.GetVars(ctx, WithResource(toRef(ev.ItemID), false), WithUser(ev.Executant, ""))
		case events.LinkRemoved:
			message = MessageLinkDeleted
			ts = utils.TSToTime(ev.Timestamp)
			vars, err = s.GetVars(ctx, WithResource(toRef(ev.ItemID), false), WithUser(ev.Executant, ""))
		case events.SpaceShared:
			message = MessageSpaceShared
			ts = ev.Timestamp
			vars, err = s.GetVars(ctx, WithSpace(ev.ID), WithUser(ev.Executant, ""), WithSharee(ev.GranteeUserID, ev.GranteeGroupID))
		case events.SpaceUnshared:
			message = MessageSpaceUnshared
			ts = ev.Timestamp
			vars, err = s.GetVars(ctx, WithSpace(ev.ID), WithUser(ev.Executant, ""), WithSharee(ev.GranteeUserID, ev.GranteeGroupID))
		}

		if err != nil {
			s.log.Error().Err(err).Msg("error getting response data")
			continue
		}

		// FIXME: configurable default locale?
		loc := l10n.MustGetUserLocale(r.Context(), activeUser.GetId().GetOpaqueId(), r.Header.Get(l10n.HeaderAcceptLanguage), s.valService)
		t := l10n.NewTranslatorFromCommonConfig("en", _domain, "", _localeFS, _localeSubPath)

		resp.Activities = append(resp.Activities, NewActivity(t.Translate(message, loc), ts, e.GetId(), vars))
	}

	// delete activities in separate go routine
	if len(toDelete) > 0 {
		go func() {
			err := s.RemoveActivities(rid, toDelete)
			if err != nil {
				s.log.Error().Err(err).Msg("error removing activities")
			}
		}()
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

func (s *ActivitylogService) getFilters(query string) (*provider.ResourceId, int, func(RawActivity) bool, func(*ehmsg.Event) bool, error) {
	qast, err := kql.Builder{}.Build(query)
	if err != nil {
		return nil, 0, nil, nil, err
	}

	prefilters := make([]func(RawActivity) bool, 0)
	postfilters := make([]func(*ehmsg.Event) bool, 0)

	var (
		itemID string
		limit  int
	)

	for _, n := range qast.Nodes {
		switch v := n.(type) {
		case *ast.StringNode:
			switch strings.ToLower(v.Key) {
			case "itemid":
				itemID = v.Value
			case "depth":
				depth, err := strconv.Atoi(v.Value)
				if err != nil {
					return nil, limit, nil, nil, err
				}

				prefilters = append(prefilters, func(a RawActivity) bool {
					return a.Depth <= depth
				})
			case "limit":
				l, err := strconv.Atoi(v.Value)
				if err != nil {
					return nil, limit, nil, nil, err
				}

				limit = l
			}
		case *ast.DateTimeNode:
			switch v.Operator.Value {
			case "<", "<=":
				prefilters = append(prefilters, func(a RawActivity) bool {
					return a.Timestamp.Before(v.Value)
				})
			case ">", ">=":
				prefilters = append(prefilters, func(a RawActivity) bool {
					return a.Timestamp.After(v.Value)
				})
			}
		case *ast.OperatorNode:
			if v.Value != "AND" {
				return nil, limit, nil, nil, errors.New("only AND operator is supported")
			}
		}
	}

	rid, err := storagespace.ParseID(itemID)
	if err != nil {
		return nil, limit, nil, nil, err
	}
	pref := func(a RawActivity) bool {
		for _, f := range prefilters {
			if !f(a) {
				return false
			}
		}
		return true
	}
	postf := func(e *ehmsg.Event) bool {
		for _, f := range postfilters {
			if !f(e) {
				return false
			}
		}
		return true
	}
	return &rid, limit, pref, postf, nil
}

// returns true if this is just a rename
func isRename(o, n *provider.Reference) bool {
	// if resourceids are different we assume it is a move
	if !utils.ResourceIDEqual(o.GetResourceId(), n.GetResourceId()) {
		return false
	}
	return filepath.Base(o.GetPath()) != filepath.Base(n.GetPath())
}
