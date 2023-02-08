package service

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"go-micro.dev/v4/store"
)

// Comment when you read this on review
var _adminid = "2502d8b8-a5e7-4ab3-b858-5aafae4a64a2"

// UserlogService is the service responsible for user activities
type UserlogService struct {
	log              log.Logger
	ch               <-chan events.Event
	m                *chi.Mux
	store            store.Store
	cfg              *config.Config
	historyClient    ehsvc.EventHistoryService
	gwClient         gateway.GatewayAPIClient
	registeredEvents map[string]events.Unmarshaller
}

// NewUserlogService returns an EventHistory service
func NewUserlogService(cfg *config.Config, consumer events.Consumer, store store.Store, gwClient gateway.GatewayAPIClient, registeredEvents []events.Unmarshaller, log log.Logger) (*UserlogService, error) {
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

	ul := &UserlogService{
		log:              log,
		ch:               ch,
		store:            store,
		cfg:              cfg,
		historyClient:    c,
		gwClient:         gwClient,
		registeredEvents: make(map[string]events.Unmarshaller),
	}

	ul.BuildMux()

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
		var (
			spaceID    string       // we need the spaceid to inform other space members
			spaceOwner *user.UserId // we need a space owner to query space members
		)
		switch e := event.Event.(type) {
		case events.UploadReady:
			spaceID = e.FileRef.GetResourceId().GetSpaceId()
			spaceOwner = e.SpaceOwner
		default:
			ul.log.Error().Interface("event", e).Msg("unhandled event")
			continue
		}

		// for each event type we need to:
		// I) find users eligible to receive the event
		users, err := ul.findEligibleUsers(spaceOwner, spaceID, event.Type)
		if err != nil {
			ul.log.Error().Err(err).Str("spaceID", spaceID).Msg("failed to find eligible users")
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
		return nil, err
	}

	if len(rec) == 0 {
		// no events available
		return nil, nil
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

	return resp.Events, nil

}

// DeleteEvents will delete the specified events
func (ul *UserlogService) DeleteEvents(userid string, evids []string) error {
	return ul.removeEventsFromUser(userid, evids...)
}

func (ul *UserlogService) addEventsToUser(userid string, eventids ...string) error {
	return ul.alterUserEventList(userid, func(ids []string) []string {
		return append(ids, eventids...)
	})
}

func (ul *UserlogService) removeEventsFromUser(userid string, eventids ...string) error {
	toDelete := make(map[string]struct{})
	for _, e := range eventids {
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

	b, err := json.Marshal(ids)
	if err != nil {
		return err
	}

	return ul.store.Write(&store.Record{
		Key:   userid,
		Value: b,
	})

}

func (ul *UserlogService) findEligibleUsers(spaceOwner *user.UserId, spaceID string, evtype string) ([]string, error) {
	ctx, _, err := utils.Impersonate(spaceOwner, ul.gwClient, ul.cfg.MachineAuthAPIKey)
	if err != nil {
		return nil, err
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
		if users, err = ul.gatherSpaceMembers(space, evtype); err != nil {
			return nil, err
		}
	default:
		// TODO: shares? other space types?
		return nil, fmt.Errorf("unsupported space type: %s", space.SpaceType)
	}

	return users, nil
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

func (ul *UserlogService) gatherSpaceMembers(space *storageprovider.StorageSpace, evtype string) ([]string, error) {
	var permissionsMap map[string]*storageprovider.ResourcePermissions
	if err := json.Unmarshal(space.Opaque.Map["grants"].GetValue(), &permissionsMap); err != nil {
		return nil, err
	}

	groupsMap := make(map[string]struct{})
	if opaqueGroups, ok := space.Opaque.Map["groups"]; ok {
		_ = json.Unmarshal(opaqueGroups.GetValue(), &groupsMap)
	}

	// we use a map to avoid duplicates
	usermap := make(map[string]struct{})
	for id, perm := range permissionsMap {
		if _, isGroup := groupsMap[id]; isGroup {
			usrs, err := ul.resolveGroup(id)
			if err != nil {
				ul.log.Error().Err(err).Str("groupID", id).Msg("failed to resolve group")
				continue
			}

			for _, u := range usrs {
				usermap[u] = struct{}{}
			}
			continue
		}

		// TODO: needed permission is depended on the event
		if ul.evaluatePermission(perm, evtype) {
			usermap[id] = struct{}{}
		}
	}

	var users []string
	for id := range usermap {
		users = append(users, id)
	}

	return users, nil
}

// resolves the users of a group
func (ul *UserlogService) resolveGroup(groupID string) ([]string, error) {
	// TODO: Please implement me!
	return []string{}, nil
}

func (ul *UserlogService) evaluatePermission(perms *storageprovider.ResourcePermissions, evtype string) bool {
	switch evtype {
	case "events.UploadReady":
		return perms.Stat
	default:
		return false
	}
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
