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
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/go-chi/chi/v5"
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

type itemsList struct {
	Value []*libregraph.DriveItem
}

var _ = Describe("Driveitems", func() {
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

		currentUser = &userpb.User{
			Id: &userpb.UserId{
				OpaqueId: "user",
			},
		}
	)

	BeforeEach(func() {
		eventsPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		rr = httptest.NewRecorder()

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

	Describe("GetRootDriveChildren", func() {
		It("handles ListStorageSpaces not found ", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListStorageSpaces error", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewInternal(ctx, "internal error"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("handles ListContainer not found", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListContainer permission denied", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewPermissionDenied(ctx, errors.New("denied"), "denied"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})

		It("handles ListContainer error", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewInternal(ctx, "internal"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("succeeds", func() {
			mtime := time.Now()
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status:        status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{{Owner: currentUser, Root: &provider.ResourceId{}}},
			}, nil)
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewOK(ctx),
				Infos: []*provider.ResourceInfo{
					{
						Type:  provider.ResourceType_RESOURCE_TYPE_FILE,
						Id:    &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
						Etag:  "etag",
						Mtime: utils.TimeToTS(mtime),
					},
				},
			}, nil)
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me/drive/root/children", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetRootDriveChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := itemsList{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))
			Expect(res.Value[0].GetLastModifiedDateTime().Equal(mtime)).To(BeTrue())
			Expect(res.Value[0].GetETag()).To(Equal("etag"))
			Expect(res.Value[0].GetId()).To(Equal("storageid$spaceid!opaqueid"))
		})
	})

	Describe("GetDriveItemChildren", func() {
		It("handles ListContainer not found", func() {
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewNotFound(ctx, "not found"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListContainer permission denied as not found", func() {
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewPermissionDenied(ctx, errors.New("denied"), "denied"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusNotFound))
		})

		It("handles ListContainer error", func() {
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewInternal(ctx, "internal"),
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("succeeds", func() {
			mtime := time.Now()
			gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
				Status: status.NewOK(ctx),
				Infos: []*provider.ResourceInfo{
					{
						Type:  provider.ResourceType_RESOURCE_TYPE_FILE,
						Id:    &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
						Etag:  "etag",
						Mtime: utils.TimeToTS(mtime),
					},
				},
			}, nil)
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "storageid$spaceid")
			rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetDriveItemChildren(rr, r)
			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := itemsList{}

			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))
			Expect(res.Value[0].GetLastModifiedDateTime().Equal(mtime)).To(BeTrue())
			Expect(res.Value[0].GetETag()).To(Equal("etag"))
			Expect(res.Value[0].GetId()).To(Equal("storageid$spaceid!opaqueid"))
		})
	})
})
