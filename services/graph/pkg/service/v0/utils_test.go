package svc_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	rConversions "github.com/owncloud/reva/v2/pkg/conversions"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/utils"
	cs3mocks "github.com/owncloud/reva/v2/tests/cs3mocks/mocks"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/reva/v2/pkg/storagespace"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

var _ = Describe("Utils", func() {
	DescribeTable("GetDriveAndItemIDParam",
		func(driveID, itemID string, shouldPass bool) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", driveID)
			rctx.URLParams.Add("itemID", itemID)

			extractedDriveID, extractedItemID, err := service.GetDriveAndItemIDParam(
				httptest.NewRequest(http.MethodGet, "/", nil).
					WithContext(
						context.WithValue(context.Background(), chi.RouteCtxKey, rctx),
					),
				conversions.ToPointer(log.NopLogger()),
			)

			switch shouldPass {
			case true:
				Expect(err).To(BeNil())
				parsedItemID, _ := storagespace.ParseID(itemID)
				Expect(extractedItemID).To(BeComparableTo(&parsedItemID, protocmp.Transform()))

				parsedDriveID, _ := storagespace.ParseID(driveID)
				Expect(extractedDriveID).To(BeComparableTo(&parsedDriveID, protocmp.Transform()))
			default:
				Expect(err).ToNot(BeNil())
			}
		},
		Entry("fails: invalid driveID", "", "1$2!3", false),
		Entry("fails: invalid itemID", "1$2", "", false),
		Entry("fails: incompatible driveID and itemID", "1$2", "3$4!5", false),
		Entry("fails: no itemID opaqueId", "1$2", "3$4", false),
		Entry("pass: valid driveID and itemID", "1$2", "1$2!5", true),
	)

	DescribeTable("IsSpaceRoot",
		func(resourceID *provider.ResourceId, isRoot bool) {
			Expect(service.IsSpaceRoot(resourceID)).To(Equal(isRoot))
		},
		Entry("spaceId and opaqueID equal", &provider.ResourceId{
			StorageId: "1",
			OpaqueId:  "2",
			SpaceId:   "2",
		}, true),
		Entry("nil", nil, false),
		Entry("spaceID empty", &provider.ResourceId{
			StorageId: "1",
			OpaqueId:  "2",
		}, false),
		Entry("opaqueID empty", &provider.ResourceId{
			StorageId: "1",
			SpaceId:   "3",
		}, false),
		Entry("spaceID and opaqueID unequal", &provider.ResourceId{
			OpaqueId: "2",
			SpaceId:  "3",
		}, false),
	)

	DescribeTable("IsShareJail",
		func(resourceID *provider.ResourceId, isShareJail bool) {
			Expect(service.IsShareJail(resourceID)).To(Equal(isShareJail))
		},
		Entry("valid: share jail", &provider.ResourceId{
			StorageId: utils.ShareStorageProviderID,
			SpaceId:   utils.ShareStorageSpaceID,
		}, true),
		Entry("invalid: empty storageId", &provider.ResourceId{
			SpaceId: utils.ShareStorageSpaceID,
		}, false),
		Entry("invalid: empty spaceId", &provider.ResourceId{
			StorageId: utils.ShareStorageProviderID,
		}, false),
		Entry("invalid: empty storageId and spaceId", &provider.ResourceId{}, false),
		Entry("invalid: non share jail storageId", &provider.ResourceId{
			StorageId: "123",
			SpaceId:   utils.ShareStorageSpaceID,
		}, false),
		Entry("invalid: non share jail spaceId", &provider.ResourceId{
			StorageId: utils.ShareStorageProviderID,
			SpaceId:   "123",
		}, false),
		Entry("invalid: non share jail storageID and spaceId", &provider.ResourceId{
			StorageId: "123",
			SpaceId:   "123",
		}, false),
	)

	DescribeTable("_cs3ReceivedShareToLibreGraphPermissions",
		func(permissionSet *provider.ResourcePermissions, match func(*libregraph.Permission)) {
			permission, err := service.CS3ReceivedShareToLibreGraphPermissions(
				context.Background(),
				nil,
				identity.IdentityCache{},
				&collaboration.ReceivedShare{
					Share: &collaboration.Share{
						Permissions: &collaboration.SharePermissions{
							Permissions: permissionSet,
						},
					},
				}, &provider.ResourceInfo{
					Type: provider.ResourceType_RESOURCE_TYPE_FILE,
				},
				unifiedrole.GetRoles(unifiedrole.RoleFilterAll()),
			)
			Expect(err).ToNot(HaveOccurred())
			match(permission)
		},
		Entry(
			"permissions match a role",
			rConversions.NewViewerRole().CS3ResourcePermissions(),
			func(p *libregraph.Permission) {
				Expect(p.GetRoles()).To(HaveExactElements([]string{unifiedrole.UnifiedRoleViewerID}))
				Expect(p.GetLibreGraphPermissionsActions()).To(BeNil())
			},
		),
		Entry(
			"permissions do not match any role",
			&provider.ResourcePermissions{
				AddGrant: true,
			},
			func(p *libregraph.Permission) {
				Expect(p.GetRoles()).To(BeNil())
				Expect(p.GetLibreGraphPermissionsActions()).To(HaveExactElements([]string{unifiedrole.DriveItemPermissionsCreate}))
			},
		),
	)

	Describe("cs3ReceivedSharesToDriveItems", func() {
		It("sets spaceId in remoteItem", func() {
			pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
			gatewayClient := &cs3mocks.GatewayAPIClient{}
			gatewaySelector := pool.GetSelector[gateway.GatewayAPIClient](
				"GatewaySelector",
				"com.owncloud.api.gateway",
				func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
					return gatewayClient
				},
			)

			statResponse := &storageprovider.StatResponse{
				Status: &rpc.Status{Code: rpc.Code_CODE_OK},
				Info: &storageprovider.ResourceInfo{
					Id: &storageprovider.ResourceId{
						StorageId: "storage-id-123",
						SpaceId:   "space-id-456",
						OpaqueId:  "opaque-id-789",
					},
					Name:  "shared-file.txt",
					Type:  storageprovider.ResourceType_RESOURCE_TYPE_FILE,
					Etag:  "etag-abc",
					Size:  1024,
					Mtime: &typesv1beta1.Timestamp{Seconds: 1234567890},
					Space: &storageprovider.StorageSpace{
						SpaceType: "project",
						Root: &storageprovider.ResourceId{
							StorageId: "storage-id-123",
							SpaceId:   "space-id-456",
							OpaqueId:  "space-id-456",
						},
					},
				},
			}

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(statResponse, nil)

			getUserResponse := &userpb.GetUserResponse{
				Status: &rpc.Status{Code: rpc.Code_CODE_OK},
				User: &userpb.User{
					Id: &userpb.UserId{
						OpaqueId: "user-123",
					},
					DisplayName: "Test User",
					Mail:        "test@example.com",
				},
			}
			gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(getUserResponse, nil)

			identityCache := identity.NewIdentityCache(
				identity.IdentityCacheWithGatewaySelector(gatewaySelector),
			)

			receivedShares := []*collaboration.ReceivedShare{
				{
					Share: &collaboration.Share{
						Id: &collaboration.ShareId{
							OpaqueId: "share-123",
						},
						ResourceId: &storageprovider.ResourceId{
							StorageId: "storage-id-123",
							SpaceId:   "space-id-456",
							OpaqueId:  "opaque-id-789",
						},
						Permissions: &collaboration.SharePermissions{
							Permissions: rConversions.NewViewerRole().CS3ResourcePermissions(),
						},
						Creator: &userpb.UserId{
							OpaqueId: "user-123",
						},
						Ctime: &typesv1beta1.Timestamp{Seconds: 1234567890},
					},
					State: collaboration.ShareState_SHARE_STATE_ACCEPTED,
					MountPoint: &storageprovider.Reference{
						Path: "shared-file.txt",
					},
				},
			}

			driveItems, err := service.Cs3ReceivedSharesToDriveItems(
				context.Background(),
				conversions.ToPointer(log.NopLogger()),
				gatewayClient,
				identityCache,
				receivedShares,
				unifiedrole.GetRoles(unifiedrole.RoleFilterAll()),
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(driveItems).To(HaveLen(1))

			driveItem := driveItems[0]
			Expect(driveItem.RemoteItem).ToNot(BeNil())

			Expect(driveItem.RemoteItem.HasSpaceId()).To(BeTrue())
			expectedSpaceId := storagespace.FormatStorageID("storage-id-123", "space-id-456")
			Expect(driveItem.RemoteItem.GetSpaceId()).To(Equal(expectedSpaceId))
		})
	})
})
