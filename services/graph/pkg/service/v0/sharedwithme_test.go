package svc_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaborationv1beta1 "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
)

var _ = Describe("SharedWithMe", func() {
	var (
		svc             service.Service
		cfg             *config.Config
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		identityBackend *identitymocks.Backend
		ctx             context.Context
		tape            *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		identityBackend = &identitymocks.Backend{}

		tape = httptest.NewRecorder()
		ctx = context.Background()

		cfg = defaults.FullDefaultConfig()
		cfg.Identity.LDAP.CACert = "" // skip the startup checks, we don't use LDAP at all in this tests
		cfg.TokenManager.JWTSecret = "loremipsum"
		cfg.Commons = &shared.Commons{}
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}

		svc, _ = service.NewService(
			service.Config(cfg),
			service.WithGatewaySelector(gatewaySelector),
			service.WithIdentityBackend(identityBackend),
		)
	})

	Describe("ListSharedWithMe", func() {
		var (
			listReceivedSharesResponse *collaborationv1beta1.ListReceivedSharesResponse
			statResponse               *providerv1beta1.StatResponse
			getUserResponse            *userv1beta1.GetUserResponse
		)

		toResourceID := func(in string) *providerv1beta1.ResourceId {
			out, err := storagespace.ParseID(in)
			Expect(err).To(BeNil())

			return &out
		}

		BeforeEach(func() {

			getUserResponse = &userv1beta1.GetUserResponse{
				Status: status.NewOK(ctx),
				User: &userv1beta1.User{
					Id: &userv1beta1.UserId{
						OpaqueId: "2699b42d-c6ca-4ce1-90de-89dedfb3022c",
					},
					DisplayName: "John Romero",
				},
			}

			gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(getUserResponse, nil)

			listReceivedSharesResponse = &collaborationv1beta1.ListReceivedSharesResponse{
				Status: status.NewOK(ctx),
				Shares: []*collaborationv1beta1.ReceivedShare{
					{Share: &collaborationv1beta1.Share{ResourceId: toResourceID("1$2!3")}},
				},
			}

			gatewayClient.On("ListReceivedShares", mock.Anything, mock.Anything).Return(listReceivedSharesResponse, nil)

			statResponse = &providerv1beta1.StatResponse{
				Status: status.NewOK(ctx),
				Info: &providerv1beta1.ResourceInfo{
					Type: providerv1beta1.ResourceType_RESOURCE_TYPE_CONTAINER,
				},
			}

			gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(func(_ context.Context, r *providerv1beta1.StatRequest, _ ...grpc.CallOption) (*providerv1beta1.StatResponse, error) {
				for _, share := range listReceivedSharesResponse.Shares {
					if share.Share.ResourceId != r.Ref.ResourceId {
						continue
					}

					if statResponse.Info.Id == nil {
						statResponse.Info.Id = share.Share.ResourceId
					}

					return statResponse, nil
				}

				return nil, nil
			})
		})

		It("fails if no received shares were found", func() {
			listReceivedSharesResponse.Status = status.NewNotFound(ctx, "msg")

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			Expect(tape.Code, errorcode.ItemNotFound)
		})

		It("ignores hidden received shares by default", func() {
			listReceivedSharesResponse.Shares = append(listReceivedSharesResponse.Shares, &collaborationv1beta1.ReceivedShare{
				Hidden: true,
			})

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value")

			Expect(len(listReceivedSharesResponse.Shares)).To(Equal(2))
			Expect(jsonData.Get("#").Num).To(Equal(float64(1)))
		})

		It("includes hidden shares if explicitly stated", func() {
			listReceivedSharesResponse.Shares = append(listReceivedSharesResponse.Shares, &collaborationv1beta1.ReceivedShare{
				Hidden: true,
				Share: &collaborationv1beta1.Share{
					ResourceId: toResourceID("7$8!9"),
				},
			})

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe?show-hidden=true", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value")

			Expect(len(listReceivedSharesResponse.Shares)).To(Equal(2))
			Expect(jsonData.Get("#").Num).To(Equal(float64(2)))
		})

		It("populates the driveItem properties", func() {

			share := listReceivedSharesResponse.Shares[0].Share
			share.Id = &collaborationv1beta1.ShareId{OpaqueId: "1:2:3"}
			share.Ctime = &typesv1beta1.Timestamp{Seconds: 4000}
			share.Mtime = &typesv1beta1.Timestamp{Seconds: 40000}

			etag := "5ffb8e4bec7026050af7fde9482b289a"

			resourceInfo := statResponse.Info
			resourceInfo.Name = "some folder"
			resourceInfo.Etag = "\"" + etag + "\""

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0")

			Expect(jsonData.Get("createdDateTime").String()).To(Equal(utils.TSToTime(share.Ctime).Format(time.RFC3339Nano)))
			Expect(jsonData.Get("eTag").String()).To(Equal(etag))
			Expect(jsonData.Get("id").String()).To(Equal(share.Id.OpaqueId))
			Expect(jsonData.Get("lastModifiedDateTime").String()).To(Equal(utils.TSToTime(share.Mtime).Format(time.RFC3339Nano)))
			Expect(jsonData.Get("name").String()).To(Equal(resourceInfo.Name))
		})

		It("populates the driveItem parentReference properties", func() {
			share := listReceivedSharesResponse.Shares[0].Share
			share.Id = &collaborationv1beta1.ShareId{OpaqueId: "1:2:3"}

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.parentReference")

			Expect(jsonData.Get("driveId").String()).To(Equal(share.Id.OpaqueId))
		})

		It("populates the driveItem remoteItem properties", func() {

			share := listReceivedSharesResponse.Shares[0].Share

			resourceInfo := statResponse.Info
			resourceInfo.Name = "some folder"
			resourceInfo.Mtime = &typesv1beta1.Timestamp{Seconds: 40000}
			resourceInfo.Size = 500

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.remoteItem")

			Expect(jsonData.Get("id").String()).To(Equal(storagespace.FormatResourceID(*share.ResourceId)))
			Expect(jsonData.Get("lastModifiedDateTime").String()).To(Equal(utils.TSToTime(resourceInfo.Mtime).Format(time.RFC3339Nano)))
			Expect(jsonData.Get("name").String()).To(Equal(resourceInfo.Name))
			Expect(jsonData.Get("size").Num).To(Equal(float64(resourceInfo.Size)))
		})

		It("populates the driveItem.remoteItem.createdBy properties", func() {
			driveOwner := getUserResponse.User

			resourceInfo := statResponse.Info
			resourceInfo.Owner = driveOwner.Id

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.remoteItem.createdBy")

			Expect(jsonData.Get("user.displayName").String()).To(Equal(driveOwner.DisplayName))
			Expect(jsonData.Get("user.id").String()).To(Equal(driveOwner.Id.OpaqueId))
		})

		It("populates the driveItem.remoteItem.fileSystemInfo properties", func() {

			share := listReceivedSharesResponse.Shares[0].Share
			share.Ctime = &typesv1beta1.Timestamp{Seconds: 400}

			resourceInfo := statResponse.Info
			resourceInfo.Mtime = &typesv1beta1.Timestamp{Seconds: 4000}

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.remoteItem.fileSystemInfo")

			Expect(jsonData.Get("createdDateTime").String()).To(Equal(utils.TSToTime(share.Ctime).Format(time.RFC3339Nano)))
			Expect(jsonData.Get("lastModifiedDateTime").String()).To(Equal(utils.TSToTime(resourceInfo.Mtime).Format(time.RFC3339Nano)))
		})

		It("populates the driveItem.remoteItem.folder properties", func() {

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.remoteItem")

			Expect(jsonData.Get("file").Exists()).To(BeFalse())
			Expect(jsonData.Get("folder.childCount").Num).To(Equal(float64(0)))
		})

		It("populates the driveItem.remoteItem.file properties", func() {
			resourceInfo := statResponse.Info
			resourceInfo.Type = providerv1beta1.ResourceType_RESOURCE_TYPE_FILE
			resourceInfo.MimeType = "application/pdf"

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.remoteItem")

			Expect(jsonData.Get("folder").Exists()).To(BeFalse())
			Expect(jsonData.Get("file.mimeType").String()).To(Equal(resourceInfo.MimeType))
		})

		It("populates the driveItem.remoteItem.folder properties", func() {
			resourceInfo := statResponse.Info
			resourceInfo.Type = providerv1beta1.ResourceType_RESOURCE_TYPE_CONTAINER

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.remoteItem")

			Expect(jsonData.Get("folder").Exists()).To(BeTrue())
		})

		It("populates the driveItem.remoteItem.shared properties", func() {
			share := listReceivedSharesResponse.Shares[0].Share
			share.Ctime = &typesv1beta1.Timestamp{Seconds: 4000}

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.remoteItem.shared")

			Expect(jsonData.Get("sharedDateTime").String()).To(Equal(utils.TSToTime(share.Ctime).Format(time.RFC3339Nano)))
		})

		It("populates the driveItem.remoteItem.shared.owner properties", func() {
			shareCreator := getUserResponse.User

			share := listReceivedSharesResponse.Shares[0].Share
			share.Creator = shareCreator.Id

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.remoteItem.shared.owner")

			Expect(jsonData.Get("user.displayName").String()).To(Equal(shareCreator.DisplayName))
			Expect(jsonData.Get("user.id").String()).To(Equal(shareCreator.Id.OpaqueId))
		})

		It("populates the driveItem.remoteItem.shared.sharedBy properties", func() {
			shareCreator := getUserResponse.User

			share := listReceivedSharesResponse.Shares[0].Share
			share.Creator = shareCreator.Id

			svc.ListSharedWithMe(
				tape,
				httptest.NewRequest(http.MethodGet, "/graph/v1beta1/me/drive/sharedWithMe", nil),
			)

			jsonData := gjson.Get(tape.Body.String(), "value.0.remoteItem.shared.sharedBy")

			Expect(jsonData.Get("user.displayName").String()).To(Equal(shareCreator.DisplayName))
			Expect(jsonData.Get("user.id").String()).To(Equal(shareCreator.Id.OpaqueId))
		})
	})
})
