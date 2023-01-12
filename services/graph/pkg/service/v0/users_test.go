package svc_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/status"
	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go-micro.dev/v4/client"

	libregraph "github.com/owncloud/libre-graph-api-go"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settings "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
)

type userList struct {
	Value []*libregraph.User
}

var _ = Describe("Users", func() {
	var (
		svc             service.Service
		ctx             context.Context
		cfg             *config.Config
		gatewayClient   *mocks.GatewayClient
		eventsPublisher mocks.Publisher
		roleService     *mocks.RoleService
		identityBackend *identitymocks.Backend

		rr *httptest.ResponseRecorder

		currentUser = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "user",
			},
		}
	)

	BeforeEach(func() {
		eventsPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

		identityBackend = &identitymocks.Backend{}
		roleService = &mocks.RoleService{}
		gatewayClient = &mocks.GatewayClient{}

		rr = httptest.NewRecorder()
		ctx = context.Background()

		cfg = defaults.FullDefaultConfig()
		cfg.Identity.LDAP.CACert = "" // skip the startup checks, we don't use LDAP at all in this tests
		cfg.TokenManager.JWTSecret = "loremipsum"
		cfg.Commons = &shared.Commons{}
		cfg.GRPCClientTLS = &shared.GRPCClientTLS{}
		cfg.Application.ID = "some-application-ID"

		_ = ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
		svc, _ = service.NewService(
			service.Config(cfg),
			service.WithGatewayClient(gatewayClient),
			service.EventsPublisher(&eventsPublisher),
			service.WithIdentityBackend(identityBackend),
			service.WithRoleService(roleService),
		)
	})

	Describe("GetMe", func() {
		It("handles missing user", func() {
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me", nil)
			svc.GetMe(rr, r)

			Expect(rr.Code).To(Equal(http.StatusInternalServerError))
		})

		It("gets the information", func() {
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetMe(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
		})

		It("expands the memberOf", func() {
			user := &libregraph.User{
				Id: libregraph.PtrString("user1"),
				MemberOf: []libregraph.Group{
					{DisplayName: libregraph.PtrString("somegroup")},
				},
			}
			identityBackend.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(user, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me?$expand=memberOf", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetMe(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			responseUser := &libregraph.User{}
			err = json.Unmarshal(data, &responseUser)
			Expect(err).ToNot(HaveOccurred())

			Expect(responseUser.GetId()).To(Equal("user1"))
			Expect(responseUser.GetMemberOf()).To(HaveLen(1))
			Expect(responseUser.GetMemberOf()[0].GetDisplayName()).To(Equal("somegroup"))

		})

		It("expands the appRoleAssignments", func() {
			assignments := []*settingsmsg.UserRoleAssignment{
				{
					Id:          "some-appRoleAssignment-ID",
					AccountUuid: "user",
					RoleId:      "some-appRole-ID",
				},
			}
			roleService.On("ListRoleAssignments", mock.Anything, mock.Anything, mock.Anything).Return(&settings.ListRoleAssignmentsResponse{Assignments: assignments}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me?$expand=appRoleAssignments", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetMe(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			responseUser := &libregraph.User{}
			err = json.Unmarshal(data, &responseUser)
			Expect(err).ToNot(HaveOccurred())

			Expect(responseUser.GetId()).To(Equal("user"))
			Expect(responseUser.GetAppRoleAssignments()).To(HaveLen(1))
			Expect(responseUser.GetAppRoleAssignments()[0].GetId()).To(Equal("some-appRoleAssignment-ID"))
			Expect(responseUser.GetAppRoleAssignments()[0].GetAppRoleId()).To(Equal("some-appRole-ID"))
			Expect(responseUser.GetAppRoleAssignments()[0].GetPrincipalId()).To(Equal("user"))
			Expect(responseUser.GetAppRoleAssignments()[0].GetResourceId()).To(Equal("some-application-ID"))
		})
	})

	Describe("GetUsers", func() {
		It("handles invalid requests", func() {
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/users?$invalid=true", nil)
			svc.GetUsers(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("lists the users", func() {
			user := &libregraph.User{}
			user.SetId("user1")
			users := []*libregraph.User{user}

			identityBackend.On("GetUsers", mock.Anything, mock.Anything, mock.Anything).Return(users, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/users", nil)
			svc.GetUsers(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := userList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(res.Value)).To(Equal(1))
			Expect(res.Value[0].GetId()).To(Equal("user1"))
		})

		It("sorts", func() {
			user := &libregraph.User{}
			user.SetId("user1")
			user.SetMail("z@example.com")
			user.SetDisplayName("9")
			user.SetOnPremisesSamAccountName("9")
			user2 := &libregraph.User{}
			user2.SetId("user2")
			user2.SetMail("a@example.com")
			user2.SetDisplayName("1")
			user2.SetOnPremisesSamAccountName("1")
			users := []*libregraph.User{user, user2}

			identityBackend.On("GetUsers", mock.Anything, mock.Anything, mock.Anything).Return(users, nil)

			getUsers := func(path string) []*libregraph.User {
				r := httptest.NewRequest(http.MethodGet, path, nil)
				rec := httptest.NewRecorder()
				svc.GetUsers(rec, r)

				Expect(rec.Code).To(Equal(http.StatusOK))
				data, err := io.ReadAll(rec.Body)
				Expect(err).ToNot(HaveOccurred())

				res := userList{}
				err = json.Unmarshal(data, &res)
				Expect(err).ToNot(HaveOccurred())
				return res.Value
			}

			unsorted := getUsers("/graph/v1.0/users")
			Expect(len(unsorted)).To(Equal(2))
			Expect(unsorted[0].GetId()).To(Equal("user1"))
			Expect(unsorted[1].GetId()).To(Equal("user2"))

			byMail := getUsers("/graph/v1.0/users?$orderby=mail")
			Expect(len(byMail)).To(Equal(2))
			Expect(byMail[0].GetId()).To(Equal("user2"))
			Expect(byMail[1].GetId()).To(Equal("user1"))
			byMail = getUsers("/graph/v1.0/users?$orderby=mail%20asc")
			Expect(len(byMail)).To(Equal(2))
			Expect(byMail[0].GetId()).To(Equal("user2"))
			Expect(byMail[1].GetId()).To(Equal("user1"))
			byMail = getUsers("/graph/v1.0/users?$orderby=mail%20desc")
			Expect(len(byMail)).To(Equal(2))
			Expect(byMail[0].GetId()).To(Equal("user1"))
			Expect(byMail[1].GetId()).To(Equal("user2"))

			byDisplayName := getUsers("/graph/v1.0/users?$orderby=displayName")
			Expect(len(byDisplayName)).To(Equal(2))
			Expect(byDisplayName[0].GetId()).To(Equal("user2"))
			Expect(byDisplayName[1].GetId()).To(Equal("user1"))
			byDisplayName = getUsers("/graph/v1.0/users?$orderby=displayName%20asc")
			Expect(len(byDisplayName)).To(Equal(2))
			Expect(byDisplayName[0].GetId()).To(Equal("user2"))
			Expect(byDisplayName[1].GetId()).To(Equal("user1"))
			byDisplayName = getUsers("/graph/v1.0/users?$orderby=displayName%20desc")
			Expect(len(byDisplayName)).To(Equal(2))
			Expect(byDisplayName[0].GetId()).To(Equal("user1"))
			Expect(byDisplayName[1].GetId()).To(Equal("user2"))

			byOnPremisesSamAccountName := getUsers("/graph/v1.0/users?$orderby=onPremisesSamAccountName")
			Expect(len(byOnPremisesSamAccountName)).To(Equal(2))
			Expect(byOnPremisesSamAccountName[0].GetId()).To(Equal("user2"))
			Expect(byOnPremisesSamAccountName[1].GetId()).To(Equal("user1"))
			byOnPremisesSamAccountName = getUsers("/graph/v1.0/users?$orderby=onPremisesSamAccountName%20asc")
			Expect(len(byOnPremisesSamAccountName)).To(Equal(2))
			Expect(byOnPremisesSamAccountName[0].GetId()).To(Equal("user2"))
			Expect(byOnPremisesSamAccountName[1].GetId()).To(Equal("user1"))
			byOnPremisesSamAccountName = getUsers("/graph/v1.0/users?$orderby=onPremisesSamAccountName%20desc")
			Expect(len(byOnPremisesSamAccountName)).To(Equal(2))
			Expect(byOnPremisesSamAccountName[0].GetId()).To(Equal("user1"))
			Expect(byOnPremisesSamAccountName[1].GetId()).To(Equal("user2"))

			// Handles invalid sort field
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/users?$orderby=invalid", nil)
			svc.GetUsers(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("expands the appRoleAssignments", func() {

			user := &libregraph.User{}
			user.SetId("user1")
			user.SetMail("z@example.com")
			user.SetDisplayName("9")
			user.SetOnPremisesSamAccountName("9")
			user2 := &libregraph.User{}
			user2.SetId("user2")
			user2.SetMail("a@example.com")
			user2.SetDisplayName("1")
			user2.SetOnPremisesSamAccountName("1")
			users := []*libregraph.User{user, user2}
			identityBackend.On("GetUsers", mock.Anything, mock.Anything, mock.Anything).Return(users, nil)

			roleService.On("ListRoleAssignments", mock.Anything, mock.Anything, mock.Anything).Return(func(ctx context.Context, in *settings.ListRoleAssignmentsRequest, opts ...client.CallOption) *settings.ListRoleAssignmentsResponse {
				return &settings.ListRoleAssignmentsResponse{Assignments: []*settingsmsg.UserRoleAssignment{
					{
						Id:          "some-appRoleAssignment-ID",
						AccountUuid: in.GetAccountUuid(),
						RoleId:      "some-appRole-ID",
					},
				}}
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/users?$expand=appRoleAssignments", nil)
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.GetUsers(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			res := userList{}
			err = json.Unmarshal(data, &res)
			Expect(err).ToNot(HaveOccurred())

			responseUsers := res.Value
			Expect(len(responseUsers)).To(Equal(2))
			Expect(responseUsers[0].GetId()).To(Equal("user1"))
			Expect(responseUsers[0].GetAppRoleAssignments()).To(HaveLen(1))
			Expect(responseUsers[0].GetAppRoleAssignments()[0].GetId()).To(Equal("some-appRoleAssignment-ID"))
			Expect(responseUsers[0].GetAppRoleAssignments()[0].GetAppRoleId()).To(Equal("some-appRole-ID"))
			Expect(responseUsers[0].GetAppRoleAssignments()[0].GetPrincipalId()).To(Equal("user1"))
			Expect(responseUsers[0].GetAppRoleAssignments()[0].GetResourceId()).To(Equal("some-application-ID"))

			Expect(responseUsers[1].GetId()).To(Equal("user2"))
			Expect(responseUsers[1].GetAppRoleAssignments()).To(HaveLen(1))
			Expect(responseUsers[1].GetAppRoleAssignments()[0].GetId()).To(Equal("some-appRoleAssignment-ID"))
			Expect(responseUsers[1].GetAppRoleAssignments()[0].GetAppRoleId()).To(Equal("some-appRole-ID"))
			Expect(responseUsers[1].GetAppRoleAssignments()[0].GetPrincipalId()).To(Equal("user2"))
			Expect(responseUsers[1].GetAppRoleAssignments()[0].GetResourceId()).To(Equal("some-application-ID"))

		})
	})

	Describe("GetUser", func() {
		It("handles missing userids", func() {
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/users", nil)
			svc.GetUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("gets the user", func() {
			user := &libregraph.User{}
			user.SetId("user1")

			identityBackend.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(user, nil)
			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/users", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", *user.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())
			responseUser := &libregraph.User{}
			err = json.Unmarshal(data, &responseUser)
			Expect(err).ToNot(HaveOccurred())
			Expect(responseUser.GetId()).To(Equal("user1"))
			Expect(len(responseUser.GetDrives())).To(Equal(0))
		})

		It("includes the personal space if requested", func() {
			user := &libregraph.User{}
			user.SetId("user1")

			identityBackend.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(user, nil)
			gatewayClient.On("GetQuota", mock.Anything, mock.Anything, mock.Anything).Return(&provider.GetQuotaResponse{
				Status:     status.NewOK(ctx),
				TotalBytes: 10,
			}, nil)
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{
					{
						Id:        &provider.StorageSpaceId{OpaqueId: "drive1"},
						Root:      &provider.ResourceId{SpaceId: "space", OpaqueId: "space"},
						SpaceType: "project",
					},
					{
						Id:        &provider.StorageSpaceId{OpaqueId: "personal"},
						Root:      &provider.ResourceId{SpaceId: "personal", OpaqueId: "personal"},
						SpaceType: "personal",
					},
				},
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/users?$expand=drive", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", *user.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())
			responseUser := &libregraph.User{}
			err = json.Unmarshal(data, &responseUser)
			Expect(err).ToNot(HaveOccurred())
			Expect(responseUser.GetId()).To(Equal("user1"))
			Expect(*responseUser.GetDrive().Id).To(Equal("personal"))
		})

		It("includes the drives if requested", func() {
			user := &libregraph.User{}
			user.SetId("user1")

			identityBackend.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(user, nil)
			gatewayClient.On("GetQuota", mock.Anything, mock.Anything, mock.Anything).Return(&provider.GetQuotaResponse{
				Status:     status.NewOK(ctx),
				TotalBytes: 10,
			}, nil)
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{
					{
						Id:   &provider.StorageSpaceId{OpaqueId: "drive1"},
						Root: &provider.ResourceId{SpaceId: "space", OpaqueId: "space"},
					},
				},
			}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/users?$expand=drives", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", *user.Id)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())
			responseUser := &libregraph.User{}
			err = json.Unmarshal(data, &responseUser)
			Expect(err).ToNot(HaveOccurred())
			Expect(responseUser.GetId()).To(Equal("user1"))
			Expect(len(responseUser.GetDrives())).To(Equal(1))
		})

		It("expands the appRoleAssignments", func() {
			user := &libregraph.User{}
			user.SetId("user1")

			identityBackend.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(user, nil)

			assignments := []*settingsmsg.UserRoleAssignment{
				{
					Id:          "some-appRoleAssignment-ID",
					AccountUuid: "user1",
					RoleId:      "some-appRole-ID",
				},
			}
			roleService.On("ListRoleAssignments", mock.Anything, mock.Anything, mock.Anything).Return(&settings.ListRoleAssignmentsResponse{Assignments: assignments}, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/users/user1?$expand=appRoleAssignments", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", user.GetId())
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.GetUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))

			data, err := io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			responseUser := &libregraph.User{}
			err = json.Unmarshal(data, &responseUser)
			Expect(err).ToNot(HaveOccurred())

			Expect(responseUser.GetId()).To(Equal("user1"))
			Expect(responseUser.GetAppRoleAssignments()).To(HaveLen(1))
			Expect(responseUser.GetAppRoleAssignments()[0].GetId()).To(Equal("some-appRoleAssignment-ID"))
			Expect(responseUser.GetAppRoleAssignments()[0].GetAppRoleId()).To(Equal("some-appRole-ID"))
			Expect(responseUser.GetAppRoleAssignments()[0].GetPrincipalId()).To(Equal("user1"))
			Expect(responseUser.GetAppRoleAssignments()[0].GetResourceId()).To(Equal("some-application-ID"))
		})
	})

	Describe("PostUser", func() {
		var (
			user *libregraph.User

			assertHandleBadAttributes = func(user *libregraph.User) {
				userJson, err := json.Marshal(user)
				Expect(err).ToNot(HaveOccurred())

				r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/users", bytes.NewBuffer(userJson))
				svc.PostUser(rr, r)

				Expect(rr.Code).To(Equal(http.StatusBadRequest))
			}
		)

		BeforeEach(func() {
			user = &libregraph.User{}
			user.SetDisplayName("Display Name")
			user.SetOnPremisesSamAccountName("user")
			user.SetMail("user@example.com")
		})

		It("handles invalid bodies", func() {
			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/users?$invalid=true", nil)
			svc.PostUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles missing display names", func() {
			user.DisplayName = nil
			assertHandleBadAttributes(user)

		})

		It("handles missing OnPremisesSamAccountName", func() {
			user.OnPremisesSamAccountName = nil
			assertHandleBadAttributes(user)

			user.SetOnPremisesSamAccountName("")
			assertHandleBadAttributes(user)
		})

		It("handles bad Mails", func() {
			user.Mail = nil
			assertHandleBadAttributes(user)

			user.SetMail("not-a-mail-address")
			assertHandleBadAttributes(user)
		})

		It("handles set Ids - they are read-only", func() {
			user.SetId("/users/user")
			assertHandleBadAttributes(user)
		})

		It("creates a user", func() {
			roleService.On("AssignRoleToUser", mock.Anything, mock.Anything).Return(&settings.AssignRoleToUserResponse{}, nil)
			identityBackend.On("CreateUser", mock.Anything, mock.Anything).Return(func(ctx context.Context, user libregraph.User) *libregraph.User {
				user.SetId("/users/user")
				return &user
			}, nil)
			userJson, err := json.Marshal(user)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/users", bytes.NewBuffer(userJson))
			r = r.WithContext(revactx.ContextSetUser(ctx, currentUser))
			svc.PostUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
		})
	})

	Describe("DeleteUser", func() {
		It("handles missing userids", func() {
			r := httptest.NewRequest(http.MethodDelete, "/graph/v1.0/users/{userid}", nil)
			svc.DeleteUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("prevents a user from deleting themselves", func() {
			lu := libregraph.User{}
			lu.SetId(currentUser.Id.OpaqueId)
			identityBackend.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(&lu, nil)

			r := httptest.NewRequest(http.MethodDelete, "/graph/v1.0/users/{userid}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", currentUser.Id.OpaqueId)
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusForbidden))
		})

		It("deletes a user from deleting themselves", func() {
			otheruser := &userv1beta1.User{
				Id: &userv1beta1.UserId{
					OpaqueId: "otheruser",
				},
			}

			lu := libregraph.User{}
			lu.SetId(otheruser.Id.OpaqueId)
			identityBackend.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(&lu, nil)
			identityBackend.On("DeleteUser", mock.Anything, mock.Anything).Return(nil)
			gatewayClient.On("DeleteStorageSpace", mock.Anything, mock.Anything).Return(&provider.DeleteStorageSpaceResponse{
				Status: status.NewOK(ctx),
			}, nil)
			gatewayClient.On("ListStorageSpaces", mock.Anything, mock.Anything, mock.Anything).Return(&provider.ListStorageSpacesResponse{
				Status: status.NewOK(ctx),
				StorageSpaces: []*provider.StorageSpace{
					{
						Opaque:    &typesv1beta1.Opaque{},
						Id:        &provider.StorageSpaceId{OpaqueId: "drive1"},
						Root:      &provider.ResourceId{SpaceId: "space", OpaqueId: "space"},
						SpaceType: "personal",
						Owner:     otheruser,
					},
				},
			}, nil)

			r := httptest.NewRequest(http.MethodDelete, "/graph/v1.0/users/{userid}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", lu.GetId())
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.DeleteUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusNoContent))
			gatewayClient.AssertNumberOfCalls(GinkgoT(), "DeleteStorageSpace", 2) // 2 calls for the home space. first trash, then purge
		})
	})

	Describe("PatchUser", func() {
		var (
			user *libregraph.User
		)

		BeforeEach(func() {
			user = &libregraph.User{}
			user.SetDisplayName("Display Name")
			user.SetOnPremisesSamAccountName("user")
			user.SetMail("user@example.com")
			user.SetId("/users/user")

			identityBackend.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(&user, nil)
		})

		It("handles missing userids", func() {
			r := httptest.NewRequest(http.MethodPatch, "/graph/v1.0/users/{userid}", nil)
			svc.PatchUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles invalid bodies", func() {
			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/users?$invalid=true", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", user.GetId())
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PatchUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("handles invalid email", func() {
			user.SetMail("invalid")
			data, err := json.Marshal(user)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/users?$invalid=true", bytes.NewBuffer(data))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", user.GetId())
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PatchUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusBadRequest))
		})

		It("updates attributes", func() {
			identityBackend.On("UpdateUser", mock.Anything, user.GetId(), mock.Anything).Return(user, nil)

			user.SetDisplayName("New Display Name")
			data, err := json.Marshal(user)
			Expect(err).ToNot(HaveOccurred())

			r := httptest.NewRequest(http.MethodPost, "/graph/v1.0/users?$invalid=true", bytes.NewBuffer(data))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userID", user.GetId())
			r = r.WithContext(context.WithValue(revactx.ContextSetUser(ctx, currentUser), chi.RouteCtxKey, rctx))
			svc.PatchUser(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
			data, err = io.ReadAll(rr.Body)
			Expect(err).ToNot(HaveOccurred())

			updatedUser := libregraph.User{}
			err = json.Unmarshal(data, &updatedUser)
			Expect(err).ToNot(HaveOccurred())
			Expect(updatedUser.GetDisplayName()).To(Equal("New Display Name"))
		})
	})
})
