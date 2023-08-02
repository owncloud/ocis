package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/r3labs/sse/v2"
	micrometadata "go-micro.dev/v4/metadata"
	"go-micro.dev/v4/store"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
)

var tracer trace.Tracer

func init() {
	tracer = otel.Tracer("github.com/owncloud/ocis/services/userlog/pkg/service")
}

// UserlogService is the service responsible for user activities
type UserlogService struct {
	log              log.Logger
	m                *chi.Mux
	store            store.Store
	cfg              *config.Config
	historyClient    ehsvc.EventHistoryService
	gatewaySelector  pool.Selectable[gateway.GatewayAPIClient]
	valueClient      settingssvc.ValueService
	sse              *sse.Server
	registeredEvents map[string]events.Unmarshaller
}

// NewUserlogService returns an EventHistory service
func NewUserlogService(opts ...Option) (*UserlogService, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.Consumer == nil || o.Store == nil {
		return nil, fmt.Errorf("need non nil consumer (%v) and store (%v) to work properly", o.Consumer, o.Store)
	}

	ch, err := events.Consume(o.Consumer, "userlog", o.RegisteredEvents...)
	if err != nil {
		return nil, err
	}

	ul := &UserlogService{
		log:              o.Logger,
		m:                o.Mux,
		store:            o.Store,
		cfg:              o.Config,
		historyClient:    o.HistoryClient,
		gatewaySelector:  o.GatewaySelector,
		valueClient:      o.ValueClient,
		registeredEvents: make(map[string]events.Unmarshaller),
	}

	if !ul.cfg.DisableSSE {
		ul.sse = sse.New()
	}

	for _, e := range o.RegisteredEvents {
		typ := reflect.TypeOf(e)
		ul.registeredEvents[typ.String()] = e
	}

	m := roles.NewManager(
		// TODO: caching?
		roles.Logger(o.Logger),
		roles.RoleService(o.RoleClient),
	)

	ul.m.Route("/ocs/v2.php/apps/notifications/api/v1/notifications", func(r chi.Router) {
		r.Get("/", ul.HandleGetEvents)
		r.Delete("/", ul.HandleDeleteEvents)
		r.Post("/global", RequireAdminOrSecret(&m, o.Config.GlobalNotificationsSecret)(ul.HandlePostGlobalEvent))
		r.Delete("/global", RequireAdminOrSecret(&m, o.Config.GlobalNotificationsSecret)(ul.HandleDeleteGlobalEvent))

		if !ul.cfg.DisableSSE {
			r.Get("/sse", ul.HandleSSE)
		}
	})

	go ul.MemorizeEvents(ch)

	return ul, nil
}

// MemorizeEvents stores eventIDs a user wants to receive
func (ul *UserlogService) MemorizeEvents(ch <-chan events.Event) {
	for event := range ch {
		// for each event we need to:
		// I) find users eligible to receive the event
		var (
			users     []string
			executant *user.UserId
			err       error
		)

		switch e := event.Event.(type) {
		default:
			err = errors.New("unhandled event")
		// file related
		case events.PostprocessingStepFinished:
			switch e.FinishedStep {
			case events.PPStepAntivirus:
				result := e.Result.(events.VirusscanResult)
				if !result.Infected {
					continue
				}

				// TODO: should space mangers also be informed?
				users = append(users, e.ExecutingUser.GetId().GetOpaqueId())
			case events.PPStepPolicies:
				if e.Outcome == events.PPOutcomeContinue {
					continue
				}
				users = append(users, e.ExecutingUser.GetId().GetOpaqueId())
			default:
				continue

			}
		// space related // TODO: how to find spaceadmins?
		case events.SpaceDisabled:
			executant = e.Executant
			users, err = ul.findSpaceMembers(ul.impersonate(e.Executant), e.ID.GetOpaqueId(), viewer)
		case events.SpaceDeleted:
			executant = e.Executant
			for u := range e.FinalMembers {
				users = append(users, u)
			}
		case events.SpaceShared:
			executant = e.Executant
			users, err = ul.resolveID(ul.impersonate(e.Executant), e.GranteeUserID, e.GranteeGroupID)
		case events.SpaceUnshared:
			executant = e.Executant
			users, err = ul.resolveID(ul.impersonate(e.Executant), e.GranteeUserID, e.GranteeGroupID)
		case events.SpaceMembershipExpired:
			users, err = ul.resolveID(ul.impersonate(e.SpaceOwner), e.GranteeUserID, e.GranteeGroupID)

		// share related
		case events.ShareCreated:
			executant = e.Executant
			users, err = ul.resolveID(ul.impersonate(e.Executant), e.GranteeUserID, e.GranteeGroupID)
		case events.ShareRemoved:
			executant = e.Executant
			users, err = ul.resolveID(ul.impersonate(e.Executant), e.GranteeUserID, e.GranteeGroupID)
		case events.ShareExpired:
			users, err = ul.resolveID(ul.impersonate(e.ShareOwner), e.GranteeUserID, e.GranteeGroupID)
		}

		if err != nil {
			// TODO: Find out why this errors on ci pipeline
			ul.log.Debug().Err(err).Interface("event", event).Msg("error gathering members for event")
			continue
		}

		// II) filter users who want to receive the event
		// This step is postponed for later.
		// For now each user should get all events she is eligible to receive
		// ...except notifications for their own actions
		users = removeExecutant(users, executant)

		// III) store the eventID for each user
		for _, id := range users {
			if err := ul.addEventToUser(id, event); err != nil {
				ul.log.Error().Err(err).Str("userID", id).Str("eventid", event.ID).Msg("failed to store event for user")
				continue
			}
		}
	}
}

// GetEvents allows retrieving events from the eventhistory by userid
func (ul *UserlogService) GetEvents(ctx context.Context, userid string) ([]*ehmsg.Event, error) {
	ctx, span := tracer.Start(ctx, "GetEvents")
	defer span.End()
	rec, err := ul.store.Read(userid)
	if err != nil && err != store.ErrNotFound {
		ul.log.Error().Err(err).Str("userid", userid).Msg("failed to read record from store")
		return nil, err
	}

	if len(rec) == 0 {
		// no events available
		return []*ehmsg.Event{}, nil
	}

	var eventIDs []string
	if err := json.Unmarshal(rec[0].Value, &eventIDs); err != nil {
		ul.log.Error().Err(err).Str("userid", userid).Msg("failed to umarshal record from store")
		return nil, err
	}

	resp, err := ul.historyClient.GetEvents(ctx, &ehsvc.GetEventsRequest{Ids: eventIDs})
	if err != nil {
		return nil, err
	}

	// remove expired events from list asynchronously
	go func() {
		if err := ul.removeExpiredEvents(userid, eventIDs, resp.Events); err != nil {
			ul.log.Error().Err(err).Str("userid", userid).Msg("could not remove expired events from user")
		}

	}()

	return resp.Events, nil

}

// DeleteEvents will delete the specified events
func (ul *UserlogService) DeleteEvents(userid string, evids []string) error {
	toDelete := make(map[string]struct{})
	for _, e := range evids {
		toDelete[e] = struct{}{}
	}

	return ul.alterUserEventList(userid, func(ids []string) []string {
		var newids []string
		for _, id := range ids {
			if _, delete := toDelete[id]; delete {
				continue
			}

			newids = append(newids, id)
		}
		return newids
	})
}

// StoreGlobalEvent will store a global event that will be returned with each `GetEvents` request
func (ul *UserlogService) StoreGlobalEvent(ctx context.Context, typ string, data map[string]string) error {
	ctx, span := tracer.Start(ctx, "StoreGlobalEvent")
	defer span.End()
	switch typ {
	default:
		return fmt.Errorf("unknown event type: %s", typ)
	case "deprovision":
		dps, ok := data["deprovision_date"]
		if !ok {
			return errors.New("need 'deprovision_date' in request body")
		}

		format := data["deprovision_date_format"]
		if format == "" {
			format = time.RFC3339
		}

		date, err := time.Parse(format, dps)
		if err != nil {
			fmt.Println("", format, "\n", dps)
			return fmt.Errorf("cannot parse time to format. time: '%s' format: '%s'", dps, format)
		}

		ev := DeprovisionData{
			DeprovisionDate:   date,
			DeprovisionFormat: format,
		}

		b, err := json.Marshal(ev)
		if err != nil {
			return err
		}

		return ul.alterGlobalEvents(ctx, func(evs map[string]json.RawMessage) error {
			evs[typ] = b
			return nil
		})
	}

}

// GetGlobalEvents will return all global events
func (ul *UserlogService) GetGlobalEvents(ctx context.Context) (map[string]json.RawMessage, error) {
	_, span := tracer.Start(ctx, "GetGlobalEvents")
	defer span.End()
	out := make(map[string]json.RawMessage)

	recs, err := ul.store.Read(_globalEventsKey)
	if err != nil && err != store.ErrNotFound {
		return out, err
	}

	if len(recs) > 0 {
		if err := json.Unmarshal(recs[0].Value, &out); err != nil {
			return out, err
		}
	}

	return out, nil
}

// DeleteGlobalEvents will delete the specified event
func (ul *UserlogService) DeleteGlobalEvents(ctx context.Context, evnames []string) error {
	_, span := tracer.Start(ctx, "DeleteGlobalEvents")
	defer span.End()
	return ul.alterGlobalEvents(ctx, func(evs map[string]json.RawMessage) error {
		for _, name := range evnames {
			delete(evs, name)
		}
		return nil
	})
}

func (ul *UserlogService) addEventToUser(userid string, event events.Event) error {
	if !ul.cfg.DisableSSE {
		if err := ul.sendSSE(userid, event); err != nil {
			ul.log.Error().Err(err).Str("userid", userid).Str("eventid", event.ID).Msg("cannot create sse event")
		}
	}
	return ul.alterUserEventList(userid, func(ids []string) []string {
		return append(ids, event.ID)
	})
}

func (ul *UserlogService) sendSSE(userid string, event events.Event) error {
	ev, err := ul.getConverter(ul.getUserLocale(userid)).ConvertEvent(event.ID, event.Event)
	if err != nil {
		return err
	}

	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	ul.sse.Publish(userid, &sse.Event{Data: b})
	return nil
}

func (ul *UserlogService) removeExpiredEvents(userid string, all []string, received []*ehmsg.Event) error {
	exists := make(map[string]struct{}, len(received))
	for _, e := range received {
		exists[e.Id] = struct{}{}
	}

	var toDelete []string
	for _, eid := range all {
		if _, ok := exists[eid]; !ok {
			toDelete = append(toDelete, eid)
		}
	}

	if len(toDelete) == 0 {
		return nil
	}

	return ul.DeleteEvents(userid, toDelete)
}

func (ul *UserlogService) alterUserEventList(userid string, alter func([]string) []string) error {
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

	ids = alter(ids)

	// store reacts unforseeable when trying to store nil values
	if len(ids) == 0 {
		return ul.store.Delete(userid)
	}

	b, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	return ul.store.Write(&store.Record{
		Key:   userid,
		Value: b,
	})
}

func (ul *UserlogService) alterGlobalEvents(ctx context.Context, alter func(map[string]json.RawMessage) error) error {
	_, span := tracer.Start(ctx, "alterGlobalEvents")
	defer span.End()
	evs, err := ul.GetGlobalEvents(ctx)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	if err := alter(evs); err != nil {
		return err
	}

	val, err := json.Marshal(evs)
	if err != nil {
		return err
	}

	return ul.store.Write(&store.Record{
		Key:   "global-events",
		Value: val,
	})
}

// we need the spaceid to inform other space members
// we need an owner to query space members
// we need to check the user has the required role to see the event
func (ul *UserlogService) findSpaceMembers(ctx context.Context, spaceID string, requiredRole permissionChecker) ([]string, error) {
	if ctx == nil {
		return nil, errors.New("need authenticated context to find space members")
	}

	space, err := getSpace(ctx, spaceID, ul.gatewaySelector)
	if err != nil {
		return nil, err
	}

	var users []string
	switch space.SpaceType {
	case "personal":
		users = []string{space.GetOwner().GetId().GetOpaqueId()}
	case "project":
		if users, err = ul.gatherSpaceMembers(ctx, space, requiredRole); err != nil {
			return nil, err
		}
	default:
		// TODO: shares? other space types?
		return nil, fmt.Errorf("unsupported space type: %s", space.SpaceType)
	}

	return users, nil
}

func (ul *UserlogService) gatherSpaceMembers(ctx context.Context, space *storageprovider.StorageSpace, hasRequiredRole permissionChecker) ([]string, error) {
	var permissionsMap map[string]*storageprovider.ResourcePermissions
	if err := utils.ReadJSONFromOpaque(space.GetOpaque(), "grants", &permissionsMap); err != nil {
		return nil, err
	}

	groupsMap := make(map[string]struct{})
	if opaqueGroups, ok := space.Opaque.Map["groups"]; ok {
		_ = json.Unmarshal(opaqueGroups.GetValue(), &groupsMap)
	}

	// we use a map to avoid duplicates
	usermap := make(map[string]struct{})
	for id, perm := range permissionsMap {
		if !hasRequiredRole(perm) {
			// not allowed to receive event
			continue
		}

		if _, isGroup := groupsMap[id]; !isGroup {
			usermap[id] = struct{}{}
			continue
		}

		usrs, err := ul.resolveGroup(ctx, id)
		if err != nil {
			ul.log.Error().Err(err).Str("groupID", id).Msg("failed to resolve group")
			continue
		}

		for _, u := range usrs {
			usermap[u] = struct{}{}
		}
	}

	var users []string
	for id := range usermap {
		users = append(users, id)
	}

	return users, nil
}

func (ul *UserlogService) resolveID(ctx context.Context, userid *user.UserId, groupid *group.GroupId) ([]string, error) {
	if userid != nil {
		return []string{userid.GetOpaqueId()}, nil
	}

	if ctx == nil {
		return nil, errors.New("need ctx to resolve group id")
	}

	return ul.resolveGroup(ctx, groupid.GetOpaqueId())
}

// resolves the users of a group
func (ul *UserlogService) resolveGroup(ctx context.Context, groupID string) ([]string, error) {
	grp, err := getGroup(ctx, groupID, ul.gatewaySelector)
	if err != nil {
		return nil, err
	}

	var userIDs []string
	for _, m := range grp.GetMembers() {
		userIDs = append(userIDs, m.GetOpaqueId())
	}

	return userIDs, nil
}

func (ul *UserlogService) impersonate(uid *user.UserId) context.Context {
	if uid == nil {
		ul.log.Error().Msg("cannot impersonate nil user")
		return nil
	}

	u, err := getUser(context.Background(), uid, ul.gatewaySelector)
	if err != nil {
		ul.log.Error().Err(err).Msg("cannot get user")
		return nil
	}

	ctx, err := authenticate(u, ul.gatewaySelector, ul.cfg.MachineAuthAPIKey)
	if err != nil {
		ul.log.Error().Err(err).Str("userid", u.GetId().GetOpaqueId()).Msg("failed to impersonate user")
		return nil
	}
	return ctx
}

func (ul *UserlogService) getUserLocale(userid string) string {
	resp, err := ul.valueClient.GetValueByUniqueIdentifiers(
		micrometadata.Set(context.Background(), middleware.AccountID, userid),
		&settingssvc.GetValueByUniqueIdentifiersRequest{
			AccountUuid: userid,
			SettingId:   defaults.SettingUUIDProfileLanguage,
		},
	)
	if err != nil {
		ul.log.Error().Err(err).Str("userid", userid).Msg("cannot get users locale")
		return ""
	}
	val := resp.GetValue().GetValue().GetListValue().GetValues()
	if len(val) == 0 {
		return ""
	}
	return val[0].GetStringValue()
}

func (ul *UserlogService) getConverter(locale string) *Converter {
	return NewConverter(locale, ul.gatewaySelector, ul.cfg.MachineAuthAPIKey, ul.cfg.Service.Name, ul.cfg.TranslationPath)
}

func authenticate(usr *user.User, gatewaySelector pool.Selectable[gateway.GatewayAPIClient], machineAuthAPIKey string) (context.Context, error) {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	ctx := revactx.ContextSetUser(context.Background(), usr)
	authRes, err := gatewayClient.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + usr.GetId().GetOpaqueId(),
		ClientSecret: machineAuthAPIKey,
	})
	if err != nil {
		return nil, err
	}
	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error impersonating user: %s", authRes.Status.Message)
	}

	return metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, authRes.Token), nil
}

func getSpace(ctx context.Context, spaceID string, gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) (*storageprovider.StorageSpace, error) {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	res, err := gatewayClient.ListStorageSpaces(ctx, listStorageSpaceRequest(spaceID))
	if err != nil {
		return nil, err
	}

	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error while getting space: (%v) %s", res.GetStatus().GetCode(), res.GetStatus().GetMessage())
	}

	if len(res.StorageSpaces) == 0 {
		return nil, fmt.Errorf("error getting storage space %s: no space returned", spaceID)
	}

	return res.StorageSpaces[0], nil
}

func getUser(ctx context.Context, userid *user.UserId, gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) (*user.User, error) {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	getUserResponse, err := gatewayClient.GetUser(context.Background(), &user.GetUserRequest{
		UserId: userid,
	})
	if err != nil {
		return nil, err
	}

	if getUserResponse.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("error getting user: %s", getUserResponse.Status.Message)
	}

	return getUserResponse.GetUser(), nil
}

func getGroup(ctx context.Context, groupid string, gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) (*group.Group, error) {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	r, err := gatewayClient.GetGroup(ctx, &group.GetGroupRequest{GroupId: &group.GroupId{OpaqueId: groupid}})
	if err != nil {
		return nil, err
	}

	if r.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("unexpected status code from gateway client: %d", r.GetStatus().GetCode())
	}

	return r.GetGroup(), nil
}

func getResource(ctx context.Context, resourceid *storageprovider.ResourceId, gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) (*storageprovider.ResourceInfo, error) {
	gatewayClient, err := gatewaySelector.Next()
	if err != nil {
		return nil, err
	}

	res, err := gatewayClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: resourceid}})
	if err != nil {
		return nil, err
	}

	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("unexpected status code while getting space: %v", res.GetStatus().GetCode())
	}

	return res.GetInfo(), nil
}

func listStorageSpaceRequest(spaceID string) *storageprovider.ListStorageSpacesRequest {
	return &storageprovider.ListStorageSpacesRequest{
		Filters: []*storageprovider.ListStorageSpacesRequest_Filter{
			{
				Type: storageprovider.ListStorageSpacesRequest_Filter_TYPE_ID,
				Term: &storageprovider.ListStorageSpacesRequest_Filter_Id{
					Id: &storageprovider.StorageSpaceId{
						OpaqueId: spaceID,
					},
				},
			},
		},
	}
}

func removeExecutant(users []string, executant *user.UserId) []string {
	var usrs []string
	for _, u := range users {
		if u != executant.GetOpaqueId() {
			usrs = append(usrs, u)
		}
	}
	return usrs
}

type permissionChecker func(*storageprovider.ResourcePermissions) bool

func viewer(perms *storageprovider.ResourcePermissions) bool {
	return perms.Stat
}

func editor(perms *storageprovider.ResourcePermissions) bool {
	return perms.InitiateFileUpload
}

func manager(perms *storageprovider.ResourcePermissions) bool {
	return perms.DenyGrant
}
