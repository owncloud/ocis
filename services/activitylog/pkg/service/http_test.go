package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/store"
	"github.com/owncloud/reva/v2/pkg/utils"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go-micro.dev/v4/client"
	"google.golang.org/grpc"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ehmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/eventhistory/v0"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	ehmocks "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0/mocks"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/activitylog/pkg/config"
)

func newTestActivitylogService(t *testing.T, cfg *config.Config) (*ActivitylogService, *ehmocks.EventHistoryService, *provider.ResourceId) {
	t.Helper()

	rid := resourceID("item1")
	parentID := resourceID("parent1")
	selectorKey := "activitylog-test-" + t.Name()
	pool.RemoveSelector(selectorKey + "com.owncloud.api.gateway")

	gwc := cs3mocks.NewGatewayAPIClient(t)
	gwc.On("Stat", mock.Anything, mock.Anything).Return(&provider.StatResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_OK},
		Info: &provider.ResourceInfo{
			Id:       rid,
			ParentId: parentID,
			Name:     "file.txt",
			PermissionSet: &provider.ResourcePermissions{
				ListGrants: true,
			},
			Space: &provider.StorageSpace{
				Root: resourceID("spaceid"),
			},
		},
	}, nil).Maybe()
	gwc.On("GetUser", mock.Anything, mock.Anything).Return(&user.GetUserResponse{
		Status: &rpc.Status{Code: rpc.Code_CODE_OK},
		User: &user.User{
			Id:          &user.UserId{OpaqueId: "userid"},
			DisplayName: "Test User",
		},
	}, nil).Maybe()

	t.Cleanup(func() {
		pool.RemoveSelector(selectorKey + "com.owncloud.api.gateway")
	})

	gws := pool.GetSelector[gateway.GatewayAPIClient](
		selectorKey,
		"com.owncloud.api.gateway",
		func(_ grpc.ClientConnInterface) gateway.GatewayAPIClient {
			return gwc
		},
	)

	eh := ehmocks.NewEventHistoryService(t)
	if cfg == nil {
		cfg = &config.Config{}
	}
	if cfg.DefaultLanguage == "" {
		cfg.DefaultLanguage = "en"
	}

	svc := &ActivitylogService{
		cfg:       cfg,
		log:       log.NewLogger(),
		store:     store.Create(),
		gws:       gws,
		evHistory: eh,
		valService: settingssvc.MockValueService{
			GetValueByUniqueIdentifiersFunc: func(_ context.Context, _ *settingssvc.GetValueByUniqueIdentifiersRequest, _ ...client.CallOption) (*settingssvc.GetValueResponse, error) {
				return &settingssvc.GetValueResponse{}, nil
			},
		},
		registeredEvents: map[string]events.Unmarshaller{
			reflect.TypeOf(events.FileTouched{}).String(): events.FileTouched{},
		},
	}
	return svc, eh, rid
}

func makeFileTouchedEvent(t *testing.T, id string, rid *provider.ResourceId, ts time.Time) *ehmsg.Event {
	t.Helper()
	b, err := json.Marshal(events.FileTouched{
		Executant: &user.UserId{OpaqueId: "userid"},
		Ref:       &provider.Reference{ResourceId: rid},
		Timestamp: utils.TimeToTS(ts),
	})
	require.NoError(t, err)
	return &ehmsg.Event{
		Id:    id,
		Type:  reflect.TypeOf(events.FileTouched{}).String(),
		Event: b,
	}
}

func performGetActivities(t *testing.T, svc *ActivitylogService, rid *provider.ResourceId) GetActivitiesResponse {
	t.Helper()
	q := url.Values{}
	q.Set("kql", "itemid:"+storagespace.FormatResourceID(rid))
	req := httptest.NewRequest(http.MethodGet, "/graph/v1beta1/extensions/org.libregraph/activities?"+q.Encode(), nil)
	ctx := revactx.ContextSetUser(req.Context(), &user.User{Id: &user.UserId{OpaqueId: "userid"}})
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()

	svc.HandleGetItemActivities(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var out GetActivitiesResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &out))
	return out
}

func activityIDs(t *testing.T, svc *ActivitylogService, rid *provider.ResourceId) []string {
	t.Helper()
	acts, err := svc.Activities(rid)
	require.NoError(t, err)
	ids := make([]string, 0, len(acts))
	for _, a := range acts {
		ids = append(ids, a.EventID)
	}
	return ids
}

func TestHandleGetItemActivities_PreservesPendingEvents(t *testing.T) {
	svc, eh, rid := newTestActivitylogService(t, nil)
	now := time.Now()
	require.NoError(t, svc.storeActivity(storagespace.FormatResourceID(rid), "pending-event", 0, now))

	eh.On("GetEvents", mock.Anything, &ehsvc.GetEventsRequest{Ids: []string{"pending-event"}}).
		Return(&ehsvc.GetEventsResponse{}, nil).Once()

	first := performGetActivities(t, svc, rid)
	require.Empty(t, first.Activities)

	require.Never(t, func() bool {
		return len(activityIDs(t, svc, rid)) == 0
	}, 100*time.Millisecond, 10*time.Millisecond, "newly stored activity must not be deleted while eventhistory catches up")

	eh.On("GetEvents", mock.Anything, &ehsvc.GetEventsRequest{Ids: []string{"pending-event"}}).
		Return(&ehsvc.GetEventsResponse{
			Events: []*ehmsg.Event{makeFileTouchedEvent(t, "pending-event", rid, now)},
		}, nil).Once()

	second := performGetActivities(t, svc, rid)
	require.Len(t, second.Activities, 1)
	require.Equal(t, "pending-event", second.Activities[0].GetId())
}

func TestHandleGetItemActivities_DeletesOnlyExpiredMissingEvents(t *testing.T) {
	svc, eh, rid := newTestActivitylogService(t, nil)
	now := time.Now()
	resKey := storagespace.FormatResourceID(rid)

	require.NoError(t, svc.storeActivity(resKey, "recent-missing", 0, now.Add(-5*time.Minute)))
	require.NoError(t, svc.storeActivity(resKey, "zero-ts-missing", 0, time.Time{}))
	require.NoError(t, svc.storeActivity(resKey, "expired-missing", 0, now.Add(-400*time.Hour)))
	require.NoError(t, svc.storeActivity(resKey, "expired-present", 0, now.Add(-400*time.Hour)))

	eh.On("GetEvents", mock.Anything, mock.Anything).
		Return(&ehsvc.GetEventsResponse{
			Events: []*ehmsg.Event{makeFileTouchedEvent(t, "expired-present", rid, now.Add(-400*time.Hour))},
		}, nil).Once()

	resp := performGetActivities(t, svc, rid)
	require.Len(t, resp.Activities, 1)
	require.Equal(t, "expired-present", resp.Activities[0].GetId())

	require.Eventually(t, func() bool {
		ids := activityIDs(t, svc, rid)
		return len(ids) == 3
	}, time.Second, 10*time.Millisecond)

	require.ElementsMatch(t, []string{"recent-missing", "zero-ts-missing", "expired-present"}, activityIDs(t, svc, rid))
}

func TestHandleGetItemActivities_RespectsConfiguredStoreTTL(t *testing.T) {
	svc, eh, rid := newTestActivitylogService(t, &config.Config{
		Store: config.Store{
			TTL: 2 * time.Hour,
		},
	})
	now := time.Now()
	resKey := storagespace.FormatResourceID(rid)

	require.NoError(t, svc.storeActivity(resKey, "within-ttl", 0, now.Add(-1*time.Hour)))
	require.NoError(t, svc.storeActivity(resKey, "past-ttl", 0, now.Add(-3*time.Hour)))

	eh.On("GetEvents", mock.Anything, mock.Anything).
		Return(&ehsvc.GetEventsResponse{}, nil).Once()

	resp := performGetActivities(t, svc, rid)
	require.Empty(t, resp.Activities)

	require.Eventually(t, func() bool {
		ids := activityIDs(t, svc, rid)
		return len(ids) == 1
	}, time.Second, 10*time.Millisecond)

	require.Equal(t, []string{"within-ttl"}, activityIDs(t, svc, rid))
}
