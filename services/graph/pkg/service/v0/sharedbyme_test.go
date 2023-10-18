package svc_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	link "github.com/cs3org/go-cs3apis/cs3/sharing/link/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

var _ = Describe("sharedbyme", func() {
	var (
		svc             service.Service
		ctx             context.Context
		cfg             *config.Config
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		eventsPublisher mocks.Publisher
		identityBackend *identitymocks.Backend

		rr *httptest.ResponseRecorder

		newGroup *libregraph.Group
	)
	userShare := collaboration.Share{
		Id: &collaboration.ShareId{
			OpaqueId: "share-id",
		},
		ResourceId: &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "opaqueid",
		},
		Grantee: &provider.Grantee{
			Type: provider.GranteeType_GRANTEE_TYPE_USER,
			Id: &provider.Grantee_UserId{
				UserId: &userpb.UserId{
					OpaqueId: "user-id",
				},
			},
		},
	}
	groupShare := collaboration.Share{
		Id: &collaboration.ShareId{
			OpaqueId: "share-id",
		},
		ResourceId: &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "opaqueid",
		},
		Grantee: &provider.Grantee{
			Type: provider.GranteeType_GRANTEE_TYPE_GROUP,
			Id: &provider.Grantee_GroupId{
				GroupId: &grouppb.GroupId{
					OpaqueId: "group-id",
				},
			},
		},
	}
	userShareWithExpiration := collaboration.Share{
		Id: &collaboration.ShareId{
			OpaqueId: "expire-share-id",
		},
		ResourceId: &provider.ResourceId{
			StorageId: "storageid",
			SpaceId:   "spaceid",
			OpaqueId:  "expire-opaqueid",
		},
		Grantee: &provider.Grantee{
			Type: provider.GranteeType_GRANTEE_TYPE_USER,
			Id: &provider.Grantee_UserId{
				UserId: &userpb.UserId{
					OpaqueId: "user-id",
				},
			},
		},
		Expiration: utils.TimeToTS(time.Now()),
	}

	BeforeEach(func() {
		eventsPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		rr = httptest.NewRecorder()

		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		gatewayClient = &cs3mocks.GatewayAPIClient{}
		gatewayClient.On("ListPublicShares", mock.Anything, mock.Anything).Return(
			&link.ListPublicSharesResponse{
				Status: status.NewOK(ctx),
				Share:  []*link.PublicShare{},
			},
			nil,
		)
		// no stat for the image
		gatewayClient.On("Stat",
			mock.Anything,
			mock.MatchedBy(
				func(req *provider.StatRequest) bool {
					return req.Ref.ResourceId.OpaqueId == userShareWithExpiration.ResourceId.OpaqueId
				})).
			Return(&provider.StatResponse{
				Status: status.NewOK(ctx),
				Info: &provider.ResourceInfo{
					Id: userShareWithExpiration.ResourceId,
				},
			}, nil)
		gatewayClient.On("Stat",
			mock.Anything,
			mock.Anything).
			Return(&provider.StatResponse{
				Status: status.NewOK(ctx),
				Info: &provider.ResourceInfo{
					Id: userShare.ResourceId,
				},
			}, nil)
		gatewaySelector = pool.GetSelector[gateway.GatewayAPIClient](
			"GatewaySelector",
			"com.owncloud.api.gateway",
			func(cc *grpc.ClientConn) gateway.GatewayAPIClient {
				return gatewayClient
			},
		)

		identityBackend = &identitymocks.Backend{}
		newGroup = libregraph.NewGroup()
		newGroup.SetMembersodataBind([]string{"/users/user1"})
		newGroup.SetId("group1")

		rr = httptest.NewRecorder()
		ctx = context.Background()

		cfg = defaults.FullDefaultConfig()
		cfg.Identity.LDAP.CACert = "" // skip the startup checks, we don't use LDAP at all in this tests
		cfg.TokenManager.JWTSecret = "loremipsum"
		cfg.Commons = &shared.Commons{}
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}

		svc, _ = service.NewService(
			service.Config(cfg),
			service.WithGatewaySelector(gatewaySelector),
			service.EventsPublisher(&eventsPublisher),
			service.WithIdentityBackend(identityBackend),
		)
	})

	Describe("GetSharedByMe", func() {
		It("handles a failing ListShares", func() {
			gatewayClient.On("ListShares", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives/sharedByMe", nil)
			svc.GetSharedByMe(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("handles ListShares returning an error status", func() {
			gatewayClient.On("ListShares", mock.Anything, mock.Anything).Return(
				&collaboration.ListSharesResponse{Status: status.NewInternal(ctx, "error listing shares")},
				nil,
			)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives/sharedByMe", nil)
			svc.GetSharedByMe(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("succeeds, when no shares are returned", func() {
			gatewayClient.On("ListShares", mock.Anything, mock.Anything).Return(
				&collaboration.ListSharesResponse{
					Status: status.NewOK(ctx),
					Shares: []*collaboration.Share{},
				},
				nil,
			)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives/sharedByMe", nil)
			svc.GetSharedByMe(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := itemsList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(0))
		})

		It("returns a proper driveItem, when a single user share is returned", func() {
			gatewayClient.On("ListShares", mock.Anything, mock.Anything).Return(
				&collaboration.ListSharesResponse{
					Status: status.NewOK(ctx),
					Shares: []*collaboration.Share{
						&userShare,
					},
				},
				nil,
			)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives/sharedByMe", nil)
			svc.GetSharedByMe(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := itemsList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))

			di := res.Value[0]
			Expect(di.GetId()).To(Equal(storagespace.FormatResourceID(*userShare.GetResourceId())))
		})

		It("returns a proper driveItem, when a single group share is returned", func() {
			gatewayClient.On("ListShares", mock.Anything, mock.Anything).Return(
				&collaboration.ListSharesResponse{
					Status: status.NewOK(ctx),
					Shares: []*collaboration.Share{
						&groupShare,
					},
				},
				nil,
			)
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives/sharedByMe", nil)
			svc.GetSharedByMe(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := itemsList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))

			di := res.Value[0]
			Expect(di.GetId()).To(Equal(storagespace.FormatResourceID(*groupShare.GetResourceId())))
		})

		It("returns a single driveItem, when a mulitple shares for the same resource are returned", func() {
			gatewayClient.On("ListShares", mock.Anything, mock.Anything).Return(
				&collaboration.ListSharesResponse{
					Status: status.NewOK(ctx),
					Shares: []*collaboration.Share{
						&groupShare,
						&userShare,
					},
				},
				nil,
			)
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives/sharedByMe", nil)
			svc.GetSharedByMe(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := itemsList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))

			di := res.Value[0]
			Expect(di.GetId()).To(Equal(storagespace.FormatResourceID(*groupShare.GetResourceId())))
		})

		It("return a driveItem with the expiration date set, for expiring shares", func() {
			gatewayClient.On("ListShares", mock.Anything, mock.Anything).Return(
				&collaboration.ListSharesResponse{
					Status: status.NewOK(ctx),
					Shares: []*collaboration.Share{
						&userShareWithExpiration,
					},
				},
				nil,
			)
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drives/sharedByMe", nil)
			svc.GetSharedByMe(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := itemsList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))

			di := res.Value[0]
			Expect(di.GetId()).To(Equal(storagespace.FormatResourceID(*userShareWithExpiration.GetResourceId())))
		})
	})
})
