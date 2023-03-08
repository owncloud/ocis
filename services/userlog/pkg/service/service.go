package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"go-micro.dev/v4/store"
)

// UserlogService is the service responsible for user activities
type UserlogService struct {
	log              log.Logger
	m                *chi.Mux
	store            store.Store
	cfg              *config.Config
	historyClient    ehsvc.EventHistoryService
	gwClient         gateway.GatewayAPIClient
	registeredEvents map[string]events.Unmarshaller
}

// NewUserlogService returns an EventHistory service
func NewUserlogService(opts ...Option) (*UserlogService, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.Consumer == nil || o.Store == nil {
		return nil, fmt.Errorf("Need non nil consumer (%v) and store (%v) to work properly", o.Consumer, o.Store)
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
		gwClient:         o.GatewayClient,
		registeredEvents: make(map[string]events.Unmarshaller),
	}

	for _, e := range o.RegisteredEvents {
		typ := reflect.TypeOf(e)
		ul.registeredEvents[typ.String()] = e
	}

	ul.m.Route("/", func(r chi.Router) {
		r.Get("/*", ul.HandleGetEvents)
		r.Delete("/*", ul.HandleDeleteEvents)
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
			users []string
			err   error
		)

		switch e := event.Event.(type) {
		default:
			err = errors.New("unhandled event")
		// space related // TODO: how to find spaceadmins?
		case events.SpaceDisabled:
			users, err = ul.findSpaceMembers(ul.impersonate(e.Executant), e.ID.GetOpaqueId(), viewer)
		case events.SpaceDeleted:
			for u, _ := range e.FinalMembers {
				users = append(users, u)
			}
		case events.SpaceShared:
			users, err = ul.resolveID(ul.impersonate(e.Executant), e.GranteeUserID, e.GranteeGroupID)
		case events.SpaceUnshared:
			users, err = ul.resolveID(ul.impersonate(e.Executant), e.GranteeUserID, e.GranteeGroupID)
		case events.SpaceMembershipExpired:
			users, err = ul.resolveID(ul.impersonate(e.SpaceOwner), e.GranteeUserID, e.GranteeGroupID)

		// share related
		case events.ShareCreated:
			users, err = ul.resolveID(ul.impersonate(e.Executant), e.GranteeUserID, e.GranteeGroupID)
		case events.ShareRemoved:
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

		// III) store the eventID for each user
		for _, id := range users {
			if err := ul.addEventsToUser(id, event.ID); err != nil {
				ul.log.Error().Err(err).Str("userID", id).Str("eventid", event.ID).Msg("failed to store event for user")
				continue
			}
		}
	}
}

// GetEvents allows to retrieve events from the eventhistory by userid
func (ul *UserlogService) GetEvents(ctx context.Context, userid string) ([]*ehmsg.Event, error) {
	rec, err := ul.store.Read(userid)
	if err != nil && err != store.ErrNotFound {
		ul.log.Fatal().Err(err).Str("userid", userid).Msg("failed to read record from database")
		return nil, err
	}

	if len(rec) == 0 {
		// no events available
		return []*ehmsg.Event{}, nil
	}

	var eventIDs []string
	if err := json.Unmarshal(rec[0].Value, &eventIDs); err != nil {
		ul.log.Fatal().Err(err).Str("userid", userid).Msg("failed to umarshal record from database")
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

func (ul *UserlogService) addEventsToUser(userid string, eventids ...string) error {
	return ul.alterUserEventList(userid, func(ids []string) []string {
		return append(ids, eventids...)
	})
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

// we need the spaceid to inform other space members
// we need an owner to query space members
// we need to check the user has the required role to see the event
func (ul *UserlogService) findSpaceMembers(ctx context.Context, spaceID string, requiredRole permissionChecker) ([]string, error) {
	if ctx == nil {
		return nil, errors.New("need authenticated context to find space members")
	}

	space, err := ul.getSpace(ctx, spaceID)
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

	return ul.resolveGroup(ctx, groupid.GetOpaqueId())
}

// resolves the users of a group
func (ul *UserlogService) resolveGroup(ctx context.Context, groupID string) ([]string, error) {
	grp, err := ul.getGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var userIDs []string
	for _, m := range grp.GetMembers() {
		userIDs = append(userIDs, m.GetOpaqueId())
	}

	return userIDs, nil
}

func (ul *UserlogService) impersonate(u *user.UserId) context.Context {
	if u == nil {
		ul.log.Debug().Msg("cannot impersonate nil user")
		return nil
	}

	ctx, _, err := utils.Impersonate(u, ul.gwClient, ul.cfg.MachineAuthAPIKey)
	if err != nil {
		ul.log.Error().Err(err).Str("userid", u.GetOpaqueId()).Msg("failed to impersonate user")
		return nil
	}
	return ctx
}

func (ul *UserlogService) getSpace(ctx context.Context, spaceID string) (*storageprovider.StorageSpace, error) {
	res, err := ul.gwClient.ListStorageSpaces(ctx, listStorageSpaceRequest(spaceID))
	if err != nil {
		return nil, err
	}

	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("Unexpected status code while getting space: %v", res.GetStatus().GetCode())
	}

	if len(res.StorageSpaces) == 0 {
		return nil, fmt.Errorf("error getting storage space %s: no space returned", spaceID)
	}

	return res.StorageSpaces[0], nil
}

func (ul *UserlogService) getUser(ctx context.Context, userid *user.UserId) (*user.User, error) {
	getUserResponse, err := ul.gwClient.GetUser(context.Background(), &user.GetUserRequest{
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

func (ul *UserlogService) getGroup(ctx context.Context, groupid string) (*group.Group, error) {
	r, err := ul.gwClient.GetGroup(ctx, &group.GetGroupRequest{GroupId: &group.GroupId{OpaqueId: groupid}})
	if err != nil {
		return nil, err
	}

	if r.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("unexpected status code from gateway client: %d", r.GetStatus().GetCode())
	}

	return r.GetGroup(), nil
}

func (ul *UserlogService) getResource(ctx context.Context, resourceid *storageprovider.ResourceId) (*storageprovider.ResourceInfo, error) {
	res, err := ul.gwClient.Stat(ctx, &storageprovider.StatRequest{Ref: &storageprovider.Reference{ResourceId: resourceid}})
	if err != nil {
		return nil, err
	}

	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("Unexpected status code while getting space: %v", res.GetStatus().GetCode())
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
