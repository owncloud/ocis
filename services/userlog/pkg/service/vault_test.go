package service_test

import (
	"context"
	"encoding/json"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	group "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"github.com/owncloud/reva/v2/pkg/store"
	"github.com/owncloud/reva/v2/pkg/utils"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"go-micro.dev/v4/client"
	microstore "go-micro.dev/v4/store"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	"github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0/mocks"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/service"
)

var _ = Describe("UserlogService vault events", func() {
	const (
		member   = "member"
		outsider = "outsider"
	)

	var (
		bus testBus
		sto microstore.Store

		gatewayClient *cs3mocks.GatewayAPIClient
		checkedUsers  chan string

		vaultResourceID = &provider.ResourceId{
			StorageId: utils.VaultStorageProviderID,
			SpaceId:   "spaceid",
			OpaqueId:  "itemid",
		}
		regularResourceID = &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "itemid",
		}
		vaultSpaceID = &provider.StorageSpaceId{
			OpaqueId: storagespace.FormatResourceID(&provider.ResourceId{
				StorageId: utils.VaultStorageProviderID,
				SpaceId:   "vaultspaceid",
				OpaqueId:  "vaultspaceid",
			}),
		}
		vaultGroup = &group.GroupId{OpaqueId: "vaultgroup"}
	)

	BeforeEach(func() {
		sto = store.Create()
		bus = testBus(make(chan events.Event))
		checkedUsers = make(chan string, 10)

		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector := pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		gatewayClient.On("Authenticate", mock.Anything, mock.Anything).Return(&gateway.AuthenticateResponse{Status: &rpc.Status{Code: rpc.Code_CODE_OK}}, nil)
		gatewayClient.On("GetGroup", mock.Anything, mock.Anything).Return(&group.GetGroupResponse{
			Status: &rpc.Status{Code: rpc.Code_CODE_OK},
			Group:  &group.Group{Members: []*user.UserId{{OpaqueId: member}, {OpaqueId: outsider}}},
		}, nil)
		gatewayClient.On("CheckPermission", mock.Anything, mock.Anything).Return(func(_ context.Context, r *permissions.CheckPermissionRequest, _ ...grpc.CallOption) (*permissions.CheckPermissionResponse, error) {
			id := r.GetSubjectRef().GetUserId().GetOpaqueId()
			checkedUsers <- id
			code := rpc.Code_CODE_PERMISSION_DENIED
			if id == member {
				code = rpc.Code_CODE_OK
			}
			return &permissions.CheckPermissionResponse{Status: &rpc.Status{Code: code}}, nil
		})

		vc := &settingssvc.MockValueService{}
		vc.GetValueByUniqueIdentifiersFunc = func(ctx context.Context, req *settingssvc.GetValueByUniqueIdentifiersRequest, opts ...client.CallOption) (*settingssvc.GetValueResponse, error) {
			return &settingssvc.GetValueResponse{
				Value: &settingsmsg.ValueWithIdentifier{
					Value: &settingsmsg.Value{
						Value: &settingsmsg.Value_CollectionValue{
							CollectionValue: &settingsmsg.CollectionValue{
								Values: []*settingsmsg.CollectionOption{
									{
										Key:    "in-app",
										Option: &settingsmsg.CollectionOption_BoolValue{BoolValue: true},
									},
								},
							},
						},
					},
				},
			}, nil
		}

		_, err := service.NewUserlogService(
			service.Config(&config.Config{MaxConcurrency: 5, DisableSSE: true}),
			service.Stream(bus),
			service.Store(sto),
			service.Logger(log.NewLogger()),
			service.Mux(chi.NewMux()),
			service.GatewaySelector(gatewaySelector),
			service.HistoryClient(&mocks.EventHistoryService{}),
			service.ValueClient(vc),
			service.RegisteredEvents([]events.Unmarshaller{
				events.ShareCreated{},
				events.SpaceShared{},
			}),
			service.TraceProvider(trace.NewNoopTracerProvider()),
		)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		close(bus)
	})

	storedEvents := func(userID string) []string {
		recs, err := sto.Read(userID)
		if err != nil {
			return nil
		}
		var ids []string
		Expect(json.Unmarshal(recs[0].Value, &ids)).To(Succeed())
		return ids
	}

	DescribeTable("stores the event only for group members with the vault mode permission",
		func(ev interface{}) {
			id := bus.publish(ev)

			Eventually(func() []string { return storedEvents(member) }, time.Second).Should(ConsistOf(id))
			Consistently(func() []string { return storedEvents(outsider) }, 200*time.Millisecond).Should(BeEmpty())
		},
		Entry("Share Created", events.ShareCreated{
			Executant:      &user.UserId{OpaqueId: "sharer"},
			GranteeGroupID: vaultGroup,
			ItemID:         vaultResourceID,
		}),
		Entry("Added to Space", events.SpaceShared{
			Executant:      &user.UserId{OpaqueId: "sharer"},
			GranteeGroupID: vaultGroup,
			ID:             vaultSpaceID,
		}),
	)

	It("does not check the vault mode permission for a regular resource", func() {
		id := bus.publish(events.ShareCreated{
			Executant:      &user.UserId{OpaqueId: "sharer"},
			GranteeGroupID: vaultGroup,
			ItemID:         regularResourceID,
		})

		Eventually(func() []string { return storedEvents(member) }, time.Second).Should(ConsistOf(id))
		Eventually(func() []string { return storedEvents(outsider) }, time.Second).Should(ConsistOf(id))
		Expect(checkedUsers).To(BeEmpty())
	})
})
