package svc_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	libregraph "github.com/owncloud/libre-graph-api-go"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
)

var _ = Describe("Users", func() {
	var (
		svc             service.Service
		ctx             context.Context
		cfg             *config.Config
		gatewayClient   *mocks.GatewayClient
		eventsPublisher mocks.Publisher
		identityBackend *identitymocks.Backend

		rr *httptest.ResponseRecorder

		newGroup    *libregraph.Group
		currentUser = &userv1beta1.User{
			Id: &userv1beta1.UserId{
				OpaqueId: "user",
			},
		}
	)

	BeforeEach(func() {
		eventsPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

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

		_ = ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
		svc = service.NewService(
			service.Config(cfg),
			service.WithGatewayClient(gatewayClient),
			service.EventsPublisher(&eventsPublisher),
			service.WithIdentityBackend(identityBackend),
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
			r = r.WithContext(ctxpkg.ContextSetUser(ctx, currentUser))
			svc.GetMe(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
		})

		It("expands the user", func() {
			user := &libregraph.User{}
			identityBackend.On("GetUser", mock.Anything, mock.Anything, mock.Anything).Return(user, nil)

			r := httptest.NewRequest(http.MethodGet, "/graph/v1.0/me?$expand=memberOf", nil)
			r = r.WithContext(ctxpkg.ContextSetUser(ctx, currentUser))
			svc.GetMe(rr, r)

			Expect(rr.Code).To(Equal(http.StatusOK))
		})
	})
})
