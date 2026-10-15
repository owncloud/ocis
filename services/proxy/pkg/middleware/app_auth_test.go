package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	backendmocks "github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend/mocks"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/userroles/mocks"
	"github.com/owncloud/reva/v2/pkg/auth/scope"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/reva/v2/pkg/token/manager/jwt"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

var _ = Describe("Authenticating requests", Label("AppAuthAuthenticator"), func() {
	var authenticator Authenticator
	var roleAssigner *mocks.UserRoleAssigner
	BeforeEach(func() {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")

		roleAssigner = &mocks.UserRoleAssigner{}
		roleAssigner.EXPECT().ApplyUserRole(mock.Anything, mock.Anything).RunAndReturn(
			func(ctx context.Context, user *userv1beta1.User) (*userv1beta1.User, error) { return user, nil },
		).Maybe()

		authenticator = AppAuthAuthenticator{
			Logger:           log.NewLogger(),
			UserRoleAssigner: roleAssigner,
			RevaGatewaySelector: pool.GetSelector[gateway.GatewayAPIClient](
				"GatewaySelector",
				"com.owncloud.api.gateway",
				func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
					return mockGatewayClient{
						AuthenticateFunc: func(authType, clientID, clientSecret string) *gateway.AuthenticateResponse {
							if authType != "appauth" {
								return &gateway.AuthenticateResponse{
									Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_NOT_FOUND},
								}
							}

							if clientID == "test-user" && clientSecret == "AppPassword" {
								return &gateway.AuthenticateResponse{
									Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_OK},
									Token:  "reva-token",
									User: &userv1beta1.User{
										Id:          &userv1beta1.UserId{Idp: "testIDP", OpaqueId: "abcd-1234", Type: userv1beta1.UserType_USER_TYPE_PRIMARY},
										Username:    "alice",
										Mail:        "alice@example.prv",
										DisplayName: "Alice Wong",
									},
								}
							}

							return &gateway.AuthenticateResponse{
								Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_NOT_FOUND},
							}
						},
					}
				},
			),
		}
	})

	When("the request contains correct data", func() {
		It("should successfully authenticate", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.SetBasicAuth("test-user", "AppPassword")

			req2, err := authenticator.Authenticate(req)

			Expect(err).ToNot(HaveOccurred())
			Expect(req2).ToNot(BeNil())
			Expect(req2.Header.Get("x-access-token")).To(Equal("reva-token"))

			user, ok := revactx.ContextGetUser(req2.Context())
			Expect(ok).To(Equal(true))
			Expect(user.GetId().GetIdp()).To(Equal("testIDP"))
			Expect(user.GetId().GetOpaqueId()).To(Equal("abcd-1234"))
			Expect(user.GetUsername()).To(Equal("alice"))
			Expect(user.GetMail()).To(Equal("alice@example.prv"))
		})
	})

	When("multi-instance mode is enabled", func() {
		It("rejects the request, since app passwords carry no tenant-membership claims", func() {
			multiInstanceAuthenticator := AppAuthAuthenticator{
				Logger:               log.NewLogger(),
				UserRoleAssigner:     roleAssigner,
				MultiInstanceEnabled: true,
				RevaGatewaySelector: pool.GetSelector[gateway.GatewayAPIClient](
					"GatewaySelector",
					"com.owncloud.api.gateway",
					func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
						return mockGatewayClient{
							AuthenticateFunc: func(authType, clientID, clientSecret string) *gateway.AuthenticateResponse {
								return &gateway.AuthenticateResponse{Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_OK}}
							},
						}
					},
				),
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.SetBasicAuth("test-user", "AppPassword")

			req2, err := multiInstanceAuthenticator.Authenticate(req)

			Expect(err).To(HaveOccurred())
			Expect(req2).To(BeNil())
		})
	})

	When("the account resolver is configured with a non-default PROXY_USER_OIDC_CLAIM", func() {
		It("still resolves the already-authenticated user, bypassing claim-based lookup", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/graph/v1.0/me", http.NoBody)
			req.SetBasicAuth("test-user", "AppPassword")

			authenticatedReq, err := authenticator.Authenticate(req)
			Expect(err).ToNot(HaveOccurred())

			resolvedUser := &userv1beta1.User{
				Id:       &userv1beta1.UserId{Idp: "testIDP", OpaqueId: "abcd-1234"},
				Username: "alice",
			}
			tokenManager, err := jwt.New(map[string]interface{}{"secret": "change-me", "expires": int64(60)})
			Expect(err).ToNot(HaveOccurred())
			s, err := scope.AddOwnerScope(nil)
			Expect(err).ToNot(HaveOccurred())
			token, err := tokenManager.MintToken(context.Background(), resolvedUser, s)
			Expect(err).ToNot(HaveOccurred())

			ub := backendmocks.UserBackend{}
			// services/proxy/pkg/middleware/account_resolver.go's verifyUser always looks
			// up the user by "username", regardless of the configured PROXY_USER_OIDC_CLAIM,
			// because the user arrives already authenticated via the request context.
			ub.On("GetUserByClaims", mock.Anything, "username", "alice", mock.Anything).Return(resolvedUser, token, nil)

			resolver := AccountResolver(
				Logger(log.NewLogger()),
				UserProvider(&ub),
				UserRoleAssigner(roleAssigner),
				UserOIDCClaim("sub"),
				UserCS3Claim("userid"),
			)(mockHandler{})

			rw := httptest.NewRecorder()
			resolver.ServeHTTP(rw, authenticatedReq)

			Expect(rw.Code).To(Equal(http.StatusOK))
		})
	})

	When("the request contains incorrect data", func() {
		It("should not successfully authenticate", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.SetBasicAuth("test-user", "WrongAppPassword")

			req2, err := authenticator.Authenticate(req)

			Expect(err).To(HaveOccurred())
			Expect(req2).To(BeNil())
		})
	})
})
