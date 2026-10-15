package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend/mocks"
	userRoleMocks "github.com/owncloud/ocis/v2/services/proxy/pkg/userroles/mocks"
	"github.com/owncloud/reva/v2/pkg/auth/scope"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/token/manager/jwt"
)

var _ = Describe("Authenticating requests", Label("BasicAuthenticator"), func() {
	var authenticator Authenticator
	var roleAssigner *userRoleMocks.UserRoleAssigner
	ub := mocks.UserBackend{}
	ub.On("Authenticate", mock.Anything, "testuser", "testpassword").Return(
		&userv1beta1.User{
			Id: &userv1beta1.UserId{
				Idp:      "IdpId",
				OpaqueId: "OpaqueId",
			},
			Username: "testuser",
			Mail:     "testuser@example.com",
		},
		"",
		nil,
	)
	ub.On("Authenticate", mock.Anything, mock.Anything, mock.Anything).Return(nil, "", backend.ErrAccountNotFound)

	BeforeEach(func() {
		roleAssigner = &userRoleMocks.UserRoleAssigner{}
		roleAssigner.EXPECT().ApplyUserRole(mock.Anything, mock.Anything).RunAndReturn(
			func(ctx context.Context, user *userv1beta1.User) (*userv1beta1.User, error) { return user, nil },
		).Maybe()

		authenticator = BasicAuthenticator{
			Logger:           log.NewLogger(),
			UserProvider:     &ub,
			UserRoleAssigner: roleAssigner,
		}
	})

	When("the request contains correct data", func() {
		It("should successfully authenticate", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.SetBasicAuth("testuser", "testpassword")

			req2, err := authenticator.Authenticate(req)

			Expect(err).ToNot(HaveOccurred())
			Expect(req2).ToNot(BeNil())
		})
		It("puts the authenticated user into the request context", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.SetBasicAuth("testuser", "testpassword")

			req2, err := authenticator.Authenticate(req)
			Expect(err).ToNot(HaveOccurred())

			user, ok := revactx.ContextGetUser(req2.Context())
			Expect(ok).To(Equal(true))
			Expect(user.GetId().GetIdp()).To(Equal("IdpId"))
			Expect(user.GetId().GetOpaqueId()).To(Equal("OpaqueId"))
			Expect(user.GetUsername()).To(Equal("testuser"))
			Expect(user.GetMail()).To(Equal("testuser@example.com"))
		})
	})

	When("multi-instance mode is enabled", func() {
		It("rejects the request, since basic auth credentials carry no tenant-membership claims", func() {
			multiInstanceAuthenticator := BasicAuthenticator{
				Logger:               log.NewLogger(),
				UserProvider:         &ub,
				UserRoleAssigner:     roleAssigner,
				MultiInstanceEnabled: true,
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.SetBasicAuth("testuser", "testpassword")

			req2, err := multiInstanceAuthenticator.Authenticate(req)

			Expect(err).To(HaveOccurred())
			Expect(req2).To(BeNil())
		})
	})

	When("the account resolver is configured with a non-default PROXY_USER_OIDC_CLAIM", func() {
		It("still resolves the already-authenticated user, bypassing claim-based lookup", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/graph/v1.0/me", http.NoBody)
			req.SetBasicAuth("testuser", "testpassword")

			authenticatedReq, err := authenticator.Authenticate(req)
			Expect(err).ToNot(HaveOccurred())

			resolvedUser := &userv1beta1.User{
				Id:       &userv1beta1.UserId{Idp: "IdpId", OpaqueId: "OpaqueId"},
				Username: "testuser",
			}
			tokenManager, err := jwt.New(map[string]interface{}{"secret": "change-me", "expires": int64(60)})
			Expect(err).ToNot(HaveOccurred())
			s, err := scope.AddOwnerScope(nil)
			Expect(err).ToNot(HaveOccurred())
			token, err := tokenManager.MintToken(context.Background(), resolvedUser, s)
			Expect(err).ToNot(HaveOccurred())

			resolverBackend := mocks.UserBackend{}
			// services/proxy/pkg/middleware/account_resolver.go's verifyUser always looks
			// up the user by "username", regardless of the configured PROXY_USER_OIDC_CLAIM,
			// because the user arrives already authenticated via the request context.
			resolverBackend.On("GetUserByClaims", mock.Anything, "username", "testuser", mock.Anything).Return(resolvedUser, token, nil)

			roleAssigner := &userRoleMocks.UserRoleAssigner{}
			roleAssigner.EXPECT().ApplyUserRole(mock.Anything, mock.Anything).RunAndReturn(
				func(ctx context.Context, user *userv1beta1.User) (*userv1beta1.User, error) { return user, nil },
			).Maybe()

			resolver := AccountResolver(
				Logger(log.NewLogger()),
				UserProvider(&resolverBackend),
				UserRoleAssigner(roleAssigner),
				UserOIDCClaim("sub"),
				UserCS3Claim("userid"),
			)(mockHandler{})

			rw := httptest.NewRecorder()
			resolver.ServeHTTP(rw, authenticatedReq)

			Expect(rw.Code).To(Equal(http.StatusOK))
		})
	})
})
