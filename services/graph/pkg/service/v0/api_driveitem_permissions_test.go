package svc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	collaboration "github.com/cs3org/go-cs3apis/cs3/sharing/collaboration/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	svc "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
)

var _ = Describe("DriveItemPermissionsService", func() {
	var (
		driveItemPermissionsService svc.DriveItemPermissionsService
		gatewayClient               *cs3mocks.GatewayAPIClient
		gatewaySelector             *mocks.Selectable[gateway.GatewayAPIClient]
		currentUser                 = &userpb.User{
			Id: &userpb.UserId{
				OpaqueId: "user",
			},
		}
	)

	BeforeEach(func() {
		logger := log.NewLogger()
		gatewayClient = cs3mocks.NewGatewayAPIClient(GinkgoT())

		gatewaySelector = mocks.NewSelectable[gateway.GatewayAPIClient](GinkgoT())
		gatewaySelector.On("Next").Return(gatewayClient, nil)

		cache := identity.NewIdentityCache(identity.IdentityCacheWithGatewaySelector(gatewaySelector))

		service, err := svc.NewDriveItemPermissionsService(logger, gatewaySelector, cache, false)
		Expect(err).ToNot(HaveOccurred())
		driveItemPermissionsService = service
	})

	Describe("Invite", func() {
		var (
			createShareResponse *collaboration.CreateShareResponse
			driveItemInvite     libregraph.DriveItemInvite
			driveItemId         provider.ResourceId
			statResponse        *provider.StatResponse
			getUserResponse     *userpb.GetUserResponse
			getGroupResponse    *grouppb.GetGroupResponse
		)

		BeforeEach(func() {
			driveItemId = provider.ResourceId{
				StorageId: "1",
				SpaceId:   "2",
				OpaqueId:  "3",
			}
			ctx := revactx.ContextSetUser(context.Background(), currentUser)

			statResponse = &provider.StatResponse{
				Status: status.NewOK(ctx),
			}
			gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(statResponse, nil)

			getUserResponse = &userpb.GetUserResponse{
				Status: status.NewOK(ctx),
				User: &userpb.User{
					Id:          &userpb.UserId{OpaqueId: "1"},
					DisplayName: "Cem Kaner",
				},
			}

			getGroupResponse = &grouppb.GetGroupResponse{
				Status: status.NewOK(ctx),
				Group: &grouppb.Group{
					Id:          &grouppb.GroupId{OpaqueId: "2"},
					GroupName:   "Florida Institute of Technology",
					DisplayName: "Florida Institute of Technology",
				},
			}

			createShareResponse = &collaboration.CreateShareResponse{
				Status: status.NewOK(ctx),
			}
		})

		It("creates user shares as expected (happy path)", func() {
			gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(getUserResponse, nil)
			gatewayClient.On("CreateShare", mock.Anything, mock.Anything).Return(createShareResponse, nil)
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				{ObjectId: libregraph.PtrString("1"), LibreGraphRecipientType: libregraph.PtrString("user")},
			}
			driveItemInvite.ExpirationDateTime = libregraph.PtrTime(time.Now().Add(time.Hour))
			createShareResponse.Share = &collaboration.Share{
				Id:         &collaboration.ShareId{OpaqueId: "123"},
				Expiration: utils.TimeToTS(*driveItemInvite.ExpirationDateTime),
			}

			permission, err := driveItemPermissionsService.Invite(context.Background(), driveItemId, driveItemInvite)
			Expect(err).ToNot(HaveOccurred())
			Expect(permission.GetId()).To(Equal("123"))
			Expect(permission.GetExpirationDateTime().Equal(*driveItemInvite.ExpirationDateTime)).To(BeTrue())
			Expect(permission.GrantedToV2.User.GetDisplayName()).To(Equal(getUserResponse.User.DisplayName))
			Expect(permission.GrantedToV2.User.GetId()).To(Equal("1"))
		})

		It("creates group shares as expected (happy path)", func() {
			gatewayClient.On("GetGroup", mock.Anything, mock.Anything).Return(getGroupResponse, nil)
			gatewayClient.On("CreateShare", mock.Anything, mock.Anything).Return(createShareResponse, nil)
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				{ObjectId: libregraph.PtrString("2"), LibreGraphRecipientType: libregraph.PtrString("group")},
			}
			driveItemInvite.ExpirationDateTime = libregraph.PtrTime(time.Now().Add(time.Hour))
			createShareResponse.Share = &collaboration.Share{
				Id:         &collaboration.ShareId{OpaqueId: "123"},
				Expiration: utils.TimeToTS(*driveItemInvite.ExpirationDateTime),
			}

			permission, err := driveItemPermissionsService.Invite(context.Background(), driveItemId, driveItemInvite)
			Expect(err).ToNot(HaveOccurred())
			Expect(permission.GetId()).To(Equal("123"))
			Expect(permission.GetExpirationDateTime().Equal(*driveItemInvite.ExpirationDateTime)).To(BeTrue())
			Expect(permission.GrantedToV2.Group.GetDisplayName()).To(Equal(getGroupResponse.Group.DisplayName))
			Expect(permission.GrantedToV2.Group.GetId()).To(Equal("2"))
		})

		It("with roles (happy path)", func() {
			gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(getUserResponse, nil)
			gatewayClient.On("CreateShare", mock.Anything, mock.Anything).Return(createShareResponse, nil)
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				{ObjectId: libregraph.PtrString("1"), LibreGraphRecipientType: libregraph.PtrString("user")},
			}
			driveItemInvite.Roles = []string{unifiedrole.NewViewerUnifiedRole(true).GetId()}

			permission, err := driveItemPermissionsService.Invite(context.Background(), driveItemId, driveItemInvite)
			Expect(err).ToNot(HaveOccurred())

			Expect(permission.GetRoles()).To(HaveLen(1))
			Expect(permission.GetRoles()[0]).To(Equal(unifiedrole.NewViewerUnifiedRole(true).GetId()))
		})

		It("fails with wrong role", func() {
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				{ObjectId: libregraph.PtrString("1"), LibreGraphRecipientType: libregraph.PtrString("user")},
			}
			driveItemInvite.Roles = []string{unifiedrole.NewManagerUnifiedRole().GetId()}
			permission, err := driveItemPermissionsService.Invite(context.Background(), driveItemId, driveItemInvite)

			Expect(err).To(MatchError(errorcode.New(errorcode.InvalidRequest, "role not applicable to this resource")))
			Expect(permission).To(BeZero())
		})

		It("with actions (happy path)", func() {
			gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(getUserResponse, nil)
			gatewayClient.On("CreateShare", mock.Anything, mock.Anything).Return(createShareResponse, nil)
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				{ObjectId: libregraph.PtrString("1"), LibreGraphRecipientType: libregraph.PtrString("user")},
			}
			driveItemInvite.Roles = nil
			driveItemInvite.LibreGraphPermissionsActions = []string{unifiedrole.DriveItemContentRead}

			permission, err := driveItemPermissionsService.Invite(context.Background(), driveItemId, driveItemInvite)
			Expect(err).ToNot(HaveOccurred())

			Expect(permission).NotTo(BeZero())
			Expect(permission.GetRoles()).To(HaveLen(0))
			Expect(permission.GetLibreGraphPermissionsActions()).To(HaveLen(1))
			Expect(permission.GetLibreGraphPermissionsActions()[0]).To(Equal(unifiedrole.DriveItemContentRead))
		})
		It("fails with a missing driveritem", func() {
			statResponse.Status = status.NewNotFound(context.Background(), "not found")
			permission, err := driveItemPermissionsService.Invite(context.Background(), driveItemId, driveItemInvite)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(errorcode.New(errorcode.ItemNotFound, "not found")))
			Expect(permission).To(BeZero())
		})
	})
	Describe("SpaceRootInvite", func() {
		var (
			listSpacesResponse  *provider.ListStorageSpacesResponse
			createShareResponse *collaboration.CreateShareResponse
			driveItemInvite     libregraph.DriveItemInvite
			driveId             provider.ResourceId
			statResponse        *provider.StatResponse
			getUserResponse     *userpb.GetUserResponse
		)

		BeforeEach(func() {
			driveId = provider.ResourceId{
				StorageId: "1",
				SpaceId:   "2",
			}
			ctx := revactx.ContextSetUser(context.Background(), currentUser)

			statResponse = &provider.StatResponse{
				Status: status.NewOK(ctx),
			}

			listSpacesResponse = &provider.ListStorageSpacesResponse{
				Status: status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{
					{
						Id: &provider.StorageSpaceId{
							OpaqueId: "2",
						},
					},
				},
			}

			getUserResponse = &userpb.GetUserResponse{
				Status: status.NewOK(ctx),
				User: &userpb.User{
					Id:          &userpb.UserId{OpaqueId: "1"},
					DisplayName: "Cem Kaner",
				},
			}

			createShareResponse = &collaboration.CreateShareResponse{
				Status: status.NewOK(ctx),
			}
		})

		It("adds a user to a space as expected (happy path)", func() {
			listSpacesResponse.StorageSpaces[0].SpaceType = "project"
			listSpacesResponse.StorageSpaces[0].Root = &provider.ResourceId{
				StorageId: "1",
				SpaceId:   "2",
				OpaqueId:  "3",
			}

			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(listSpacesResponse, nil)
			gatewayClient.On("GetUser", mock.Anything, mock.Anything).Return(getUserResponse, nil)
			gatewayClient.On("Stat", mock.Anything, mock.Anything).Return(statResponse, nil)
			gatewayClient.On("CreateShare", mock.Anything, mock.Anything).Return(createShareResponse, nil)
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				{ObjectId: libregraph.PtrString("1"), LibreGraphRecipientType: libregraph.PtrString("user")},
			}
			driveItemInvite.ExpirationDateTime = libregraph.PtrTime(time.Now().Add(time.Hour))
			createShareResponse.Share = &collaboration.Share{
				Id:         &collaboration.ShareId{OpaqueId: "123"},
				Expiration: utils.TimeToTS(*driveItemInvite.ExpirationDateTime),
			}

			permission, err := driveItemPermissionsService.SpaceRootInvite(context.Background(), driveId, driveItemInvite)
			Expect(err).ToNot(HaveOccurred())
			Expect(permission.GetId()).To(Equal("123"))
			Expect(permission.GetExpirationDateTime().Equal(*driveItemInvite.ExpirationDateTime)).To(BeTrue())
			Expect(permission.GrantedToV2.User.GetDisplayName()).To(Equal(getUserResponse.User.DisplayName))
			Expect(permission.GrantedToV2.User.GetId()).To(Equal("1"))
		})
		It("rejects to add a user to a personal space", func() {
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything).Return(listSpacesResponse, nil)
			driveItemInvite.Recipients = []libregraph.DriveRecipient{
				{ObjectId: libregraph.PtrString("1"), LibreGraphRecipientType: libregraph.PtrString("user")},
			}
			driveItemInvite.ExpirationDateTime = libregraph.PtrTime(time.Now().Add(time.Hour))
			createShareResponse.Share = &collaboration.Share{
				Id:         &collaboration.ShareId{OpaqueId: "123"},
				Expiration: utils.TimeToTS(*driveItemInvite.ExpirationDateTime),
			}

			permission, err := driveItemPermissionsService.SpaceRootInvite(context.Background(), driveId, driveItemInvite)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(errorcode.New(errorcode.InvalidRequest, "unsupported space type")))
			Expect(permission).To(BeZero())
		})
	})
})

var _ = Describe("DriveItemPermissionsApi", func() {
	var (
		mockProvider *mocks.DriveItemPermissionsProvider
		httpAPI      svc.DriveItemPermissionsApi
		rCTX         *chi.Context
		invite       libregraph.DriveItemInvite
	)

	BeforeEach(func() {
		logger := log.NewLogger()

		mockProvider = mocks.NewDriveItemPermissionsProvider(GinkgoT())
		api, err := svc.NewDriveItemPermissionsApi(mockProvider, logger)
		Expect(err).ToNot(HaveOccurred())

		httpAPI = api

		rCTX = chi.NewRouteContext()
		rCTX.URLParams.Add("driveID", "1$2")

		invite = libregraph.DriveItemInvite{
			Recipients: []libregraph.DriveRecipient{
				{
					ObjectId:                libregraph.PtrString("1"),
					LibreGraphRecipientType: libregraph.PtrString("user")},
			},
			Roles: []string{unifiedrole.NewViewerUnifiedRole(true).GetId()},
		}
	})

	checkDriveIDAndItemIDValidation := func(handler http.HandlerFunc) {
		rCTX.URLParams.Add("itemID", "3$4!5")

		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", nil).
			WithContext(
				context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
			)

		handler(responseRecorder, request)

		Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))

		jsonData := gjson.Get(responseRecorder.Body.String(), "error")
		Expect(jsonData.Get("message").String()).To(Equal("invalid driveID or itemID"))
	}

	Describe("Invite", func() {
		It("validates the driveID and itemID url param", func() {
			checkDriveIDAndItemIDValidation(httpAPI.Invite)
		})

		It("return an error when the Invite provider errors", func() {
			rCTX.URLParams.Add("itemID", "1$2!3")
			responseRecorder := httptest.NewRecorder()
			inviteJson, err := json.Marshal(invite)
			Expect(err).ToNot(HaveOccurred())

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inviteJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)

			onInvite := mockProvider.On("Invite", mock.Anything, mock.Anything, mock.Anything)

			onInvite.Return(func(ctx context.Context, resourceID storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error) {
				return libregraph.Permission{}, errors.New("any")
			}).Once()

			httpAPI.Invite(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
		})

		It("call the Invite provider with the correct arguments", func() {
			rCTX.URLParams.Add("itemID", "1$2!3")
			responseRecorder := httptest.NewRecorder()
			inviteJson, err := json.Marshal(invite)
			Expect(err).ToNot(HaveOccurred())

			onInvite := mockProvider.On("Invite", mock.Anything, mock.Anything, mock.Anything)
			onInvite.Return(func(ctx context.Context, resourceID storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error) {
				Expect(storagespace.FormatResourceID(resourceID)).To(Equal("1$2!3"))
				return libregraph.Permission{}, nil
			}).Once()

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inviteJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)
			httpAPI.Invite(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
	})
	Describe("SpaceRootInvite", func() {
		It("call the Invite provider with the correct arguments", func() {
			responseRecorder := httptest.NewRecorder()
			inviteJson, err := json.Marshal(invite)
			Expect(err).ToNot(HaveOccurred())

			onInvite := mockProvider.On("SpaceRootInvite", mock.Anything, mock.Anything, mock.Anything)
			onInvite.Return(func(ctx context.Context, driveID storageprovider.ResourceId, invite libregraph.DriveItemInvite) (libregraph.Permission, error) {
				Expect(storagespace.FormatResourceID(driveID)).To(Equal("1$2"))
				return libregraph.Permission{}, nil
			}).Once()

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inviteJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)
			httpAPI.SpaceRootInvite(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})
		It("call the Invite provider with the correct arguments", func() {
			rCTX.URLParams.Add("driveID", "")
			responseRecorder := httptest.NewRecorder()
			inviteJson, err := json.Marshal(invite)
			Expect(err).ToNot(HaveOccurred())

			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(inviteJson)).
				WithContext(
					context.WithValue(context.Background(), chi.RouteCtxKey, rCTX),
				)
			httpAPI.SpaceRootInvite(responseRecorder, request)

			Expect(responseRecorder.Code).To(Equal(http.StatusUnprocessableEntity))
		})
	})
})
