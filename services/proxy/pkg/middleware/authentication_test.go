package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	oidcmocks "github.com/owncloud/ocis/v2/ocis-pkg/oidc/mocks"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/router"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend/mocks"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/stretchr/testify/mock"
	"go-micro.dev/v4/store"
	"google.golang.org/grpc"
)

var _ = Describe("authentication helpers", func() {
	DescribeTable("isPublicPath should recognize public paths",
		func(input string, expected bool) {
			isPublic := isPublicPath(input)
			Expect(isPublic).To(Equal(expected))
		},
		Entry("public files path", "/remote.php/dav/public-files/", true),
		Entry("public files path without remote.php", "/remote.php/dav/public-files/", true),
		Entry("token info path", "/ocs/v1.php/apps/files_sharing/api/v1/tokeninfo/unprotected", true),
		Entry("token info path", "/ocs/v2.php/apps/files_sharing/api/v1/tokeninfo/unprotected", true),
		Entry("capabilities", "/ocs/v1.php/cloud/capabilities", true),
		// The signing-key endpoint must NOT be a public path; otherwise a public-share
		// guest could fetch the share owner's signing key and forge signed URLs.
		Entry("signing-key v1 must not be public", "/ocs/v1.php/cloud/user/signing-key", false),
		Entry("signing-key v2 must not be public", "/ocs/v2.php/cloud/user/signing-key", false),
	)
})

var _ = Describe("Authenticating requests", Label("Authentication"), func() {
	var (
		authenticators []Authenticator
	)
	oc := oidcmocks.OIDCClient{}
	oc.On("VerifyAccessToken", mock.Anything, mock.Anything).Return(
		oidc.RegClaimsWithSID{
			SessionID: "a-session-id",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Unix(1147483647, 0)),
			},
		}, jwt.MapClaims{
			"sid": "a-session-id",
			"exp": 1147483647,
		},
		nil,
	)

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
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")

		logger := log.NewLogger()
		authenticators = []Authenticator{
			BasicAuthenticator{
				Logger:       logger,
				UserProvider: &ub,
			},
			&OIDCAuthenticator{
				OIDCIss:       "http://idp.example.com",
				Logger:        logger,
				oidcClient:    &oc,
				userInfoCache: store.NewMemoryStore(),
				skipUserInfo:  true,
			},
			PublicShareAuthenticator{
				Logger: logger,
				RevaGatewaySelector: pool.GetSelector[gateway.GatewayAPIClient](
					"GatewaySelector",
					"com.owncloud.api.gateway",
					func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
						return mockGatewayClient{
							AuthenticateFunc: func(authType, clientID, clientSecret string) *gateway.AuthenticateResponse {
								if authType != "publicshares" {
									return &gateway.AuthenticateResponse{
										Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_NOT_FOUND},
									}
								}

								if clientID == "sharetoken" && (clientSecret == "password|examples3cr3t" || clientSecret == "signature|examplesignature|exampleexpiration") {
									return &gateway.AuthenticateResponse{
										Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_OK},
										Token:  "exampletoken",
									}
								}

								if clientID == "sharetoken" && clientSecret == "password|" {
									return &gateway.AuthenticateResponse{
										Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_OK},
										Token:  "otherexampletoken",
									}
								}

								return &gateway.AuthenticateResponse{
									Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_NOT_FOUND},
								}
							},
						}
					},
				),
			},
		}
	})

	When("the public request must contains correct data", func() {
		It("ensures the context oidc data when the Bearer authentication is successful", func() {
			req := httptest.NewRequest("PROPFIND", "http://example.com/remote.php/dav/public-files/", http.NoBody)
			req = req.WithContext(router.SetRoutingInfo(context.Background(), router.RoutingInfo{}))
			req.Header.Set(_headerAuthorization, "Bearer jwt.token.sig")

			handler := Authentication(authenticators,
				EnableBasicAuth(true),
			)
			testHandler := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(oidc.FromContext(r.Context())).To(Equal(map[string]interface{}{
					"sid": "a-session-id",
					"exp": int64(1147483647),
				}))
			}))
			rr := httptest.NewRecorder()
			testHandler.ServeHTTP(rr, req)
			Expect(rr).To(HaveHTTPStatus(http.StatusOK))
		})
		It("ensures the context oidc data when user the Basic authentication is successful", func() {
			req := httptest.NewRequest("PROPFIND", "http://example.com/remote.php/dav/public-files/", http.NoBody)
			req = req.WithContext(router.SetRoutingInfo(context.Background(), router.RoutingInfo{}))
			req.SetBasicAuth("testuser", "testpassword")

			handler := Authentication(authenticators,
				EnableBasicAuth(true),
			)
			testHandler := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(oidc.FromContext(r.Context())).To(Equal(map[string]interface{}{
					"email":              "testuser@example.com",
					"ownclouduuid":       "OpaqueId",
					"iss":                "IdpId",
					"preferred_username": "testuser",
				}))
			}))
			rr := httptest.NewRecorder()
			testHandler.ServeHTTP(rr, req)
			Expect(rr).To(HaveHTTPStatus(http.StatusOK))
		})
		It("ensures the x-access-token header when public-token URL parameter is set", func() {
			req := httptest.NewRequest("PROPFIND", "http://example.com/dav/public-files/?public-token=sharetoken", http.NoBody)
			req = req.WithContext(router.SetRoutingInfo(context.Background(), router.RoutingInfo{}))

			handler := Authentication(authenticators,
				EnableBasicAuth(true),
			)
			testHandler := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header.Get(_headerRevaAccessToken)).To(Equal("otherexampletoken"))
			}))
			rr := httptest.NewRecorder()
			testHandler.ServeHTTP(rr, req)
			Expect(rr).To(HaveHTTPStatus(http.StatusOK))
		})
		It("ensures the x-access-token header when public-token URL parameter and BasicAuth are set", func() {
			req := httptest.NewRequest("PROPFIND", "http://example.com/dav/public-files/?public-token=sharetoken", http.NoBody)
			req.SetBasicAuth("public", "examples3cr3t")
			req = req.WithContext(router.SetRoutingInfo(context.Background(), router.RoutingInfo{}))

			handler := Authentication(authenticators,
				EnableBasicAuth(true),
			)
			testHandler := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header.Get(_headerRevaAccessToken)).To(Equal("exampletoken"))
			}))
			rr := httptest.NewRecorder()
			testHandler.ServeHTTP(rr, req)
			Expect(rr).To(HaveHTTPStatus(http.StatusOK))
		})
		It("ensures the x-access-token header when public-token BasicAuth is set", func() {
			req := httptest.NewRequest("GET", "http://example.com/archiver", http.NoBody)
			req.Header.Set("public-token", "sharetoken")
			req = req.WithContext(router.SetRoutingInfo(context.Background(), router.RoutingInfo{}))

			handler := Authentication(authenticators,
				EnableBasicAuth(true),
			)
			testHandler := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Expect(r.Header.Get(_headerRevaAccessToken)).To(Equal("otherexampletoken"))
			}))
			rr := httptest.NewRecorder()
			testHandler.ServeHTTP(rr, req)
			Expect(rr).To(HaveHTTPStatus(http.StatusOK))
		})
	})

	When("the userinfo call to the IdP times out", func() {
		It("returns a retryable status, not 401, for a transient backend failure", func() {
			logger := log.NewLogger()

			timingOutClient := oidcmocks.OIDCClient{}
			timingOutClient.On("VerifyAccessToken", mock.Anything, mock.Anything).Return(
				oidc.RegClaimsWithSID{
					SessionID: "a-session-id",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Unix(1147483647, 0)),
					},
				}, jwt.MapClaims{
					"sid": "a-session-id",
					"exp": 1147483647,
				},
				nil,
			)
			// what the oidc client returns for a userinfo timeout
			timingOutClient.On("UserInfo", mock.Anything, mock.Anything).Return(
				(*oidc.UserInfo)(nil),
				errors.Join(oidc.ErrTemporarilyUnavailable,
					fmt.Errorf("Get \"http://idp.example.com/userinfo\": %w", context.DeadlineExceeded)),
			)

			authenticators := []Authenticator{
				&OIDCAuthenticator{
					OIDCIss:       "http://idp.example.com",
					Logger:        logger,
					oidcClient:    &timingOutClient,
					userInfoCache: store.NewMemoryStore(),
					skipUserInfo:  false,
				},
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com/graph/v1.0/me/drives", http.NoBody)
			req = req.WithContext(router.SetRoutingInfo(context.Background(), router.RoutingInfo{}))
			req.Header.Set(_headerAuthorization, "Bearer jwt.token.sig")

			handler := Authentication(authenticators,
				EnableBasicAuth(false),
			)
			testHandler := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Fail("next handler must not be reached when userinfo times out")
			}))
			rr := httptest.NewRecorder()
			testHandler.ServeHTTP(rr, req)

			Expect(rr).To(HaveHTTPStatus(http.StatusServiceUnavailable))
		})

		It("still returns 401 when userinfo fails without a transient signal", func() {
			logger := log.NewLogger()

			failingClient := oidcmocks.OIDCClient{}
			failingClient.On("VerifyAccessToken", mock.Anything, mock.Anything).Return(
				oidc.RegClaimsWithSID{
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Unix(1147483647, 0)),
					},
				}, jwt.MapClaims{"exp": 1147483647}, nil,
			)
			failingClient.On("UserInfo", mock.Anything, mock.Anything).Return(
				(*oidc.UserInfo)(nil),
				fmt.Errorf("401 Unauthorized: token revoked"),
			)

			authenticators := []Authenticator{
				&OIDCAuthenticator{
					OIDCIss:       "http://idp.example.com",
					Logger:        logger,
					oidcClient:    &failingClient,
					userInfoCache: store.NewMemoryStore(),
					skipUserInfo:  false,
				},
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com/graph/v1.0/me/drives", http.NoBody)
			req = req.WithContext(router.SetRoutingInfo(context.Background(), router.RoutingInfo{}))
			req.Header.Set(_headerAuthorization, "Bearer jwt.token.sig")

			handler := Authentication(authenticators, EnableBasicAuth(false))
			testHandler := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Fail("next handler must not be reached for a failed authentication")
			}))
			rr := httptest.NewRecorder()
			testHandler.ServeHTTP(rr, req)

			Expect(rr).To(HaveHTTPStatus(http.StatusUnauthorized))
		})

		It("returns 401 for an invalid/malformed token", func() {
			logger := log.NewLogger()

			invalidTokenClient := oidcmocks.OIDCClient{}
			invalidTokenClient.On("VerifyAccessToken", mock.Anything, mock.Anything).Return(
				oidc.RegClaimsWithSID{}, jwt.MapClaims{}, jwt.ErrTokenMalformed,
			)

			authenticators := []Authenticator{
				&OIDCAuthenticator{
					OIDCIss:       "http://idp.example.com",
					Logger:        logger,
					oidcClient:    &invalidTokenClient,
					userInfoCache: store.NewMemoryStore(),
					skipUserInfo:  true,
				},
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com/graph/v1.0/me/drives", http.NoBody)
			req = req.WithContext(router.SetRoutingInfo(context.Background(), router.RoutingInfo{}))
			req.Header.Set(_headerAuthorization, "Bearer invalid.token.sig")

			handler := Authentication(authenticators, EnableBasicAuth(false))
			testHandler := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				Fail("next handler must not be reached for an invalid token")
			}))
			rr := httptest.NewRecorder()
			testHandler.ServeHTTP(rr, req)

			Expect(rr).To(HaveHTTPStatus(http.StatusUnauthorized))
		})
	})
})
