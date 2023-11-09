package svc_test

import (
	libregraph "github.com/owncloud/libre-graph-api-go"
	"net/http/httptest"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	cs3mocks "github.com/cs3org/reva/v2/tests/cs3mocks/mocks"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config/defaults"
	identitymocks "github.com/owncloud/ocis/v2/services/graph/pkg/identity/mocks"
	service "github.com/owncloud/ocis/v2/services/graph/pkg/service/v0"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

var _ = Describe("Language", func() {
	var (
		svc service.Service
		//ctx             context.Context
		cfg             *config.Config
		gatewayClient   *cs3mocks.GatewayAPIClient
		gatewaySelector pool.Selectable[gateway.GatewayAPIClient]
		eventsPublisher mocks.Publisher
		identityBackend *identitymocks.Backend

		rr *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		eventsPublisher.On("Publish", mock.Anything, mock.Anything, mock.Anything).Return(nil)

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

		rr = httptest.NewRecorder()
		//ctx = context.Background()

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

	It("should return the language of the current user", func() {
		user := libregraph.NewUser()
		user.SetId("disallowed")
		user.SetDisplayName("foobar")
		user.SetPreferredLanguage("en-EN")

		r := httptest.NewRequest("GET", "/graph/v1.0/me/language", nil)
		svc.(*service.Graph).GetOwnLanguage(rr, r)
		Expect(rr.Code).To(Equal(200))
		Expect(rr.Body.String()).To(Equal("en-EN"))

	})

	It("should set the language of the current user", func() {
		r := httptest.NewRequest("PUT", "/graph/v1.0/me/language/en-EN", nil)
		svc.(*service.Graph).SetOwnLanguage(rr, r)
		Expect(rr.Code).To(Equal(204))
		svc.(*service.Graph).GetOwnLanguage(rr, r)
		Expect(rr.Code).To(Equal(200))
		Expect(rr.Body.String()).To(Equal("en-EN"))
	})
})
