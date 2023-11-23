package svc_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
	"google.golang.org/grpc"

	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
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

	Describe("Invite", func() {
		var (
			itemID              string
			driveItemInvite     *libregraph.DriveItemInvite
			listSharesMock      *mock.Call
			listSharesResponse  *collaboration.ListSharesResponse
			statMock            *mock.Call
			statResponse        *provider.StatResponse
			getUserResponse     *userpb.GetUserResponse
			getUserMock         *mock.Call
			getGroupResponse    *grouppb.GetGroupResponse
			getGroupMock        *mock.Call
			createShareResponse *collaboration.CreateShareResponse
			createShareMock     *mock.Call
		)

		BeforeEach(func() {
			itemID = "f0042750-23c5-441c-9f2c-ff7c53e5bd2a$cd621428-dfbe-44c1-9393-65bf0dd440a6!1177add3-b4eb-434e-a2e8-1859b31b17bf"
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "f0042750-23c5-441c-9f2c-ff7c53e5bd2a$cd621428-dfbe-44c1-9393-65bf0dd440a6!cd621428-dfbe-44c1-9393-65bf0dd440a6")
			rctx.URLParams.Add("itemID", itemID)

			ctx = context.WithValue(ctx, chi.RouteCtxKey, rctx)
			ctx = revactx.ContextSetUser(ctx, currentUser)

			driveItemInvite = &libregraph.DriveItemInvite{
				Recipients: []libregraph.DriveRecipient{
					{ObjectId: libregraph.PtrString("1")},
				},
				Roles: []string{unifiedrole.NewViewerUnifiedRole(true).GetId()},
			}

			statMock = gatewayClient.On("Stat", mock.Anything, mock.Anything)
			statResponse = &provider.StatResponse{
				Status: status.NewOK(ctx),
			}
			statMock.Return(statResponse, nil)

			getUserMock = gatewayClient.On("GetUser", mock.Anything, mock.Anything)
			getUserResponse = &userpb.GetUserResponse{
				Status: status.NewOK(ctx),
				User: &userpb.User{
					Id:          &userpb.UserId{OpaqueId: "1"},
					DisplayName: "Cem Kaner",
				},
			}
			getUserMock.Return(getUserResponse, nil)

			getGroupMock = gatewayClient.On("GetGroup", mock.Anything, mock.Anything)
			getGroupResponse = &grouppb.GetGroupResponse{
				Status: status.NewOK(ctx),
				Group: &grouppb.Group{
					Id:        &grouppb.GroupId{OpaqueId: "2"},
					GroupName: "Florida Institute of Technology",
				},
			}
			getGroupMock.Return(getGroupResponse, nil)

			listSharesMock = gatewayClient.On("ListShares", mock.Anything, mock.Anything)
			listSharesResponse = &collaboration.ListSharesResponse{
				Status: status.NewOK(ctx),
			}
			listSharesMock.Return(listSharesResponse, nil)

			createShareMock = gatewayClient.On("CreateShare", mock.Anything, mock.Anything)
			createShareResponse = &collaboration.CreateShareResponse{
				Status: status.NewOK(ctx),
			}
			createShareMock.Return(createShareResponse, nil)
		})

		toJSONReader := func(v any) *strings.Reader {
			driveItemInviteBytes, err := json.Marshal(v)
			Expect(err).ToNot(HaveOccurred())

			return strings.NewReader(string(driveItemInviteBytes))
		}

		It("creates user and group shares as expected (happy path)", func() {
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				{ObjectId: libregraph.PtrString("1")},
				{ObjectId: libregraph.PtrString("2"), LibreGraphRecipientType: libregraph.PtrString("group")},
			}
			driveItemInvite.ExpirationDateTime = libregraph.PtrTime(time.Now().Add(time.Hour))
			createShareResponse.Share = &collaboration.Share{
				Id:         &collaboration.ShareId{OpaqueId: "123"},
				Expiration: utils.TimeToTS(*driveItemInvite.ExpirationDateTime),
			}

			svc.Invite(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
					WithContext(ctx),
			)

			jsonData := gjson.Get(rr.Body.String(), "value")

			Expect(rr.Code).To(Equal(http.StatusCreated))
			Expect(jsonData.Get("#").Num).To(Equal(float64(2)))

			Expect(jsonData.Get("0.id").Str).To(Equal("123"))
			Expect(jsonData.Get("1.id").Str).To(Equal("123"))

			Expect(jsonData.Get("0.expirationDateTime").Str).To(Equal(driveItemInvite.ExpirationDateTime.Format(time.RFC3339Nano)))
			Expect(jsonData.Get("1.expirationDateTime").Str).To(Equal(driveItemInvite.ExpirationDateTime.Format(time.RFC3339Nano)))

			Expect(jsonData.Get("#.grantedToV2.user.displayName").Array()[0].Str).To(Equal(getUserResponse.User.DisplayName))
			Expect(jsonData.Get("#.grantedToV2.user.id").Array()[0].Str).To(Equal("1"))

			Expect(jsonData.Get("#.grantedToV2.group.displayName").Array()[0].Str).To(Equal(getGroupResponse.Group.GroupName))
			Expect(jsonData.Get("#.grantedToV2.group.id").Array()[0].Str).To(Equal("2"))
		})

		It("with roles (happy path)", func() {
			svc.Invite(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
					WithContext(ctx),
			)

			jsonData := gjson.Get(rr.Body.String(), "value")

			Expect(rr.Code).To(Equal(http.StatusCreated))

			Expect(jsonData.Get(`0.@libre\.graph\.permissions\.actions`).Exists()).To(BeFalse())
			Expect(jsonData.Get("0.roles.#").Num).To(Equal(float64(1)))
			Expect(jsonData.Get("0.roles.0").String()).To(Equal(unifiedrole.NewViewerUnifiedRole(true).GetId()))
		})

		It("with actions (happy path)", func() {
			driveItemInvite.Roles = nil
			driveItemInvite.LibreGraphPermissionsActions = []string{unifiedrole.DriveItemContentRead}
			svc.Invite(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
					WithContext(ctx),
			)

			jsonData := gjson.Get(rr.Body.String(), "value")

			Expect(rr.Code).To(Equal(http.StatusCreated))

			Expect(jsonData.Get("0.roles").Exists()).To(BeFalse())
			Expect(jsonData.Get(`0.@libre\.graph\.permissions\.actions.#`).Num).To(Equal(float64(1)))
			Expect(jsonData.Get(`0.@libre\.graph\.permissions\.actions.0`).String()).To(Equal(unifiedrole.DriveItemContentRead))
		})

		It("validates the driveID", func() {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "")

			ctx = context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

			svc.Invite(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
					WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("validates the itemID", func() {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "f0042750-23c5-441c-9f2c-ff7c53e5bd2a$cd621428-dfbe-44c1-9393-65bf0dd440a6!cd621428-dfbe-44c1-9393-65bf0dd440a6")
			rctx.URLParams.Add("itemID", "")

			ctx = context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

			svc.Invite(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
					WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("checks if the itemID and driveID is compatible to each other", func() {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("driveID", "1$2!3")
			rctx.URLParams.Add("itemID", "4$5!6")

			ctx = context.WithValue(context.Background(), chi.RouteCtxKey, rctx)

			svc.Invite(
				rr,
				httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
					WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("fails if the request body is empty", func() {
			svc.Invite(
				rr,
				httptest.NewRequest(http.MethodPost, "/", nil).
					WithContext(ctx),
			)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		DescribeTable("request validations",
			func(body func() *strings.Reader, code int) {
				svc.Invite(
					rr,
					httptest.NewRequest(http.MethodPost, "/", body()).
						WithContext(ctx),
				)

				Expect(rr.Code).To(Equal(code))
			},
			Entry("fails on unknown fields", func() *strings.Reader {
				return strings.NewReader(`{"unknown":"field"}`)
			}, http.StatusBadRequest),
		)

		DescribeTable("Stat",
			func(prep func(), code int) {
				prep()
				svc.Invite(
					rr,
					httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
						WithContext(ctx),
				)

				Expect(rr.Code).To(Equal(code))
				statMock.Parent.AssertNumberOfCalls(GinkgoT(), "Stat", 1)
			},
			Entry("fails if not ok", func() {
				statResponse.Status = status.NewNotFound(context.Background(), "")
			}, http.StatusInternalServerError),
			Entry("fails if errors", func() {
				statMock.Return(nil, errors.New("error"))
			}, http.StatusInternalServerError),
		)

		DescribeTable("GetGroup",
			func(prep func(), code int) {
				driveItemInvite.Recipients = []libregraph.DriveRecipient{
					{ObjectId: libregraph.PtrString("1"), LibreGraphRecipientType: libregraph.PtrString("group")},
				}

				prep()

				svc.Invite(
					rr,
					httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
						WithContext(ctx),
				)

				Expect(rr.Code).To(Equal(code))
				getGroupMock.Parent.AssertNumberOfCalls(GinkgoT(), "GetGroup", 1)
			},
			Entry("fails if not ok", func() {
				getGroupResponse.Status = status.NewNotFound(context.Background(), "")
			}, http.StatusInternalServerError),
			Entry("fails if errors", func() {
				getGroupMock.Return(nil, errors.New("error"))
			}, http.StatusInternalServerError),
		)

		DescribeTable("GetUser",
			func(prep func(), code int) {
				prep()

				svc.Invite(
					rr,
					httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
						WithContext(ctx),
				)

				Expect(rr.Code).To(Equal(code))
				getUserMock.Parent.AssertNumberOfCalls(GinkgoT(), "GetUser", 1)
			},
			Entry("fails if not ok", func() {
				getUserResponse.Status = status.NewNotFound(context.Background(), "")
			}, http.StatusInternalServerError),
			Entry("fails if errors", func() {
				getUserMock.Return(nil, errors.New("error"))
			}, http.StatusInternalServerError),
		)

		DescribeTable("CreateShare",
			func(prep func(), code int) {
				prep()

				svc.Invite(
					rr,
					httptest.NewRequest(http.MethodPost, "/", toJSONReader(driveItemInvite)).
						WithContext(ctx),
				)

				Expect(rr.Code).To(Equal(code))
				createShareMock.Parent.AssertNumberOfCalls(GinkgoT(), "CreateShare", 1)
			},
			Entry("fails if not ok", func() {
				createShareResponse.Status = status.NewNotFound(context.Background(), "")
			}, http.StatusInternalServerError),
			Entry("fails if errors", func() {
				createShareMock.Return(nil, errors.New("error"))
			}, http.StatusInternalServerError),
		)
	})

	Describe("GetRootDriveChildren", func() {
		It("handles ListStorageSpaces not found", func() {
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

		Context("it succeeds", func() {
			var (
				r     *http.Request
				mtime = time.Now()
			)

			BeforeEach(func() {
				r = httptest.NewRequest(http.MethodGet, "/graph/v1.0/drives/storageid$spaceid/items/storageid$spaceid!nodeid/children", nil)
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("driveID", "storageid$spaceid")
				rctx.URLParams.Add("driveItemID", "storageid$spaceid!nodeid")
				r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			})

			assertItemsList := func(length int) itemsList {
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
				Expect(res.Value[0].GetId()).To(Equal("storageid$spaceid!opaqueid"))

				return res
			}

			It("returns a generic file", func() {
				gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
					Status: status.NewOK(ctx),
					Infos: []*provider.ResourceInfo{
						{
							Type:              provider.ResourceType_RESOURCE_TYPE_FILE,
							Id:                &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
							Etag:              "etag",
							Mtime:             utils.TimeToTS(mtime),
							ArbitraryMetadata: nil,
						},
					},
				}, nil)

				res := assertItemsList(1)
				Expect(res.Value[0].Audio).To(BeNil())
			})

			It("returns the audio facet if metadata is available", func() {
				gatewayClient.On("ListContainer", mock.Anything, mock.Anything).Return(&provider.ListContainerResponse{
					Status: status.NewOK(ctx),
					Infos: []*provider.ResourceInfo{
						{
							Type:     provider.ResourceType_RESOURCE_TYPE_FILE,
							Id:       &provider.ResourceId{StorageId: "storageid", SpaceId: "spaceid", OpaqueId: "opaqueid"},
							Etag:     "etag",
							Mtime:    utils.TimeToTS(mtime),
							MimeType: "audio/mpeg",
							ArbitraryMetadata: &provider.ArbitraryMetadata{
								Metadata: map[string]string{
									"libre.graph.audio.album":             "Some Album",
									"libre.graph.audio.albumArtist":       "Some AlbumArtist",
									"libre.graph.audio.artist":            "Some Artist",
									"libre.graph.audio.bitrate":           "192",
									"libre.graph.audio.composers":         "Some Composers",
									"libre.graph.audio.copyright":         "Some Copyright",
									"libre.graph.audio.disc":              "2",
									"libre.graph.audio.discCount":         "5",
									"libre.graph.audio.duration":          "225000",
									"libre.graph.audio.genre":             "Some Genre",
									"libre.graph.audio.hasDrm":            "false",
									"libre.graph.audio.isVariableBitrate": "true",
									"libre.graph.audio.title":             "Some Title",
									"libre.graph.audio.track":             "6",
									"libre.graph.audio.trackCount":        "9",
									"libre.graph.audio.year":              "1994",
								},
							},
						},
					},
				}, nil)

				res := assertItemsList(1)
				audio := res.Value[0].Audio

				Expect(audio).ToNot(BeNil())
				Expect(audio.Album).To(Equal(libregraph.PtrString("Some Album")))
				Expect(audio.AlbumArtist).To(Equal(libregraph.PtrString("Some AlbumArtist")))
				Expect(audio.Artist).To(Equal(libregraph.PtrString("Some Artist")))
				Expect(audio.Bitrate).To(Equal(libregraph.PtrInt64(192)))
				Expect(audio.Composers).To(Equal(libregraph.PtrString("Some Composers")))
				Expect(audio.Copyright).To(Equal(libregraph.PtrString("Some Copyright")))
				Expect(audio.Disc).To(Equal(libregraph.PtrInt32(2)))
				Expect(audio.DiscCount).To(Equal(libregraph.PtrInt32(5)))
				Expect(audio.Duration).To(Equal(libregraph.PtrInt64(225000)))
				Expect(audio.Genre).To(Equal(libregraph.PtrString("Some Genre")))
				Expect(audio.HasDrm).To(Equal(libregraph.PtrBool(false)))
				Expect(audio.IsVariableBitrate).To(Equal(libregraph.PtrBool(true)))
				Expect(audio.Title).To(Equal(libregraph.PtrString("Some Title")))
				Expect(audio.Track).To(Equal(libregraph.PtrInt32(6)))
				Expect(audio.TrackCount).To(Equal(libregraph.PtrInt32(9)))
				Expect(audio.Year).To(Equal(libregraph.PtrInt32(1994)))
			})
		})
	})
})
