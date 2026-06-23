package middleware

import (
	"net/http"
	"net/http/httptest"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

var _ = Describe("Authenticating requests", Label("AppAuthAuthenticator"), func() {
	var authenticator Authenticator
	BeforeEach(func() {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		authenticator = AppAuthAuthenticator{
			Logger: log.NewLogger(),
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

			req2, valid := authenticator.Authenticate(req)

			Expect(valid).To(Equal(true))
			Expect(req2).ToNot(BeNil())
			Expect(req2.Header.Get("x-access-token")).To(Equal("reva-token"))
		})
	})

	When("the request contains incorrect data", func() {
		It("should not successfully authenticate", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.SetBasicAuth("test-user", "WrongAppPassword")

			req2, valid := authenticator.Authenticate(req)

			Expect(valid).To(Equal(false))
			Expect(req2).To(BeNil())
		})
	})
})
