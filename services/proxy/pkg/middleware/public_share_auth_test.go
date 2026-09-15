package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/reva/v2/pkg/auth/manager/publicshares"
	"github.com/owncloud/reva/v2/pkg/rgrpc/todo/pool"
	"google.golang.org/grpc"
)

var _ = Describe("Authenticating requests", Label("PublicShareAuthenticator"), func() {
	var authenticator Authenticator
	BeforeEach(func() {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		authenticator = PublicShareAuthenticator{
			Logger: log.NewLogger(),
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
		}
	})
	When("the request contains correct data", func() {
		Context("using password authentication", func() {
			It("should successfully authenticate", func() {
				req := httptest.NewRequest(http.MethodGet, "http://example.com/dav/public-files/?public-token=sharetoken", http.NoBody)
				req.SetBasicAuth("public", "examples3cr3t")

				req2, valid := authenticator.Authenticate(req)

				Expect(valid).To(Equal(true))
				Expect(req2).ToNot(BeNil())

				h := req2.Header
				Expect(h.Get(_headerRevaAccessToken)).To(Equal("exampletoken"))
			})
		})
		Context("using signature authentication", func() {
			It("should successfully authenticate", func() {
				req := httptest.NewRequest(http.MethodGet, "http://example.com/dav/public-files/?public-token=sharetoken&signature=examplesignature&expiration=exampleexpiration", http.NoBody)

				req2, valid := authenticator.Authenticate(req)

				Expect(valid).To(Equal(true))
				Expect(req2).ToNot(BeNil())

				h := req2.Header
				Expect(h.Get(_headerRevaAccessToken)).To(Equal("exampletoken"))
			})
		})
	})
	When("the reguest is for the archiver", func() {
		Context("using a public-token", func() {
			It("should successfully authenticate", func() {
				req := httptest.NewRequest(http.MethodGet, "http://example.com/archiver?public-token=sharetoken", http.NoBody)
				req2, valid := authenticator.Authenticate(req)

				Expect(valid).To(Equal(true))
				Expect(req2).ToNot(BeNil())

				h := req2.Header
				Expect(h.Get(_headerRevaAccessToken)).To(Equal("otherexampletoken"))
			})
		})
		Context("not using a public-token", func() {
			It("should fail to authenticate", func() {
				req := httptest.NewRequest(http.MethodGet, "http://example.com/archiver", http.NoBody)
				req2, valid := authenticator.Authenticate(req)

				Expect(valid).To(Equal(false))
				Expect(req2).To(BeNil())
			})
		})
	})
})

var _ = Describe("Brute force protection skip marker", Label("PublicShareAuthenticator"), func() {
	// authenticatorCapturingCtx builds a PublicShareAuthenticator whose gateway
	// mock records the context passed to Authenticate, so the test can assert
	// whether the brute-force-protection skip marker was set for the request.
	authenticatorCapturingCtx := func(captured *context.Context) Authenticator {
		pool.RemoveSelector("GatewaySelector" + "com.owncloud.api.gateway")
		return PublicShareAuthenticator{
			Logger: log.NewLogger(),
			RevaGatewaySelector: pool.GetSelector[gateway.GatewayAPIClient](
				"GatewaySelector",
				"com.owncloud.api.gateway",
				func(cc grpc.ClientConnInterface) gateway.GatewayAPIClient {
					return mockGatewayClient{
						OnAuthenticate: func(ctx context.Context, in *gatewayv1beta1.AuthenticateRequest) {
							*captured = ctx
						},
						AuthenticateFunc: func(authType, clientID, clientSecret string) *gateway.AuthenticateResponse {
							// Simulate a wrong password: the gateway returns a non-OK
							// status with a nil error, matching real behaviour.
							return &gateway.AuthenticateResponse{
								Status: &rpcv1beta1.Status{Code: rpcv1beta1.Code_CODE_PERMISSION_DENIED},
							}
						},
					}
				},
			),
		}
	}

	// A failed password attempt must be counted (skip marker NOT set) whenever the
	// proxy is the sole authenticator. It may only be skipped for requests that
	// are re-authenticated downstream by ocdav (the public-files DAV paths), where
	// counting here would double count.
	DescribeTable("skips protection for a password attempt only when re-authenticated downstream",
		func(url string, expectSkip bool) {
			var capturedCtx context.Context
			authenticator := authenticatorCapturingCtx(&capturedCtx)

			req := httptest.NewRequest(http.MethodGet, url, http.NoBody)
			req.SetBasicAuth("public", "wrong-password")

			_, _ = authenticator.Authenticate(req)

			Expect(capturedCtx).ToNot(BeNil(), "gateway Authenticate should have been called")
			Expect(publicshares.CheckSkipAttempt(capturedCtx, "sharetoken")).To(Equal(expectSkip))
		},
		Entry("public-files DAV path is re-authenticated by ocdav", "http://example.com/dav/public-files/?public-token=sharetoken", true),
		Entry("public-files DAV path with remote.php is re-authenticated by ocdav", "http://example.com/remote.php/dav/public-files/?public-token=sharetoken", true),
		Entry("archiver is authenticated only by the proxy", "http://example.com/archiver?public-token=sharetoken", false),
		Entry("app open is authenticated only by the proxy", "http://example.com/app/open?public-token=sharetoken", false),
		Entry("app new is authenticated only by the proxy", "http://example.com/app/new?public-token=sharetoken", false),
		Entry("ocm is authenticated only by the proxy", "http://example.com/ocm/?public-token=sharetoken", false),
	)

	// Signature based auth is a different, non-guessable credential and must not
	// be counted towards the brute force protection, otherwise invalid or expired
	// signed URLs could trip the throttle and lock out a token.
	DescribeTable("never counts signature based auth",
		func(url string) {
			var capturedCtx context.Context
			authenticator := authenticatorCapturingCtx(&capturedCtx)

			req := httptest.NewRequest(http.MethodGet, url, http.NoBody)

			_, _ = authenticator.Authenticate(req)

			Expect(capturedCtx).ToNot(BeNil(), "gateway Authenticate should have been called")
			Expect(publicshares.CheckSkipAttempt(capturedCtx, "sharetoken")).To(BeTrue())
		},
		Entry("archiver with signature", "http://example.com/archiver?public-token=sharetoken&signature=sig&expiration=exp"),
		Entry("app open with signature", "http://example.com/app/open?public-token=sharetoken&signature=sig&expiration=exp"),
	)
})

type mockGatewayClient struct {
	gatewayv1beta1.GatewayAPIClient
	AuthenticateFunc func(authType, clientID, clientSecret string) *gatewayv1beta1.AuthenticateResponse
	// OnAuthenticate, when set, is invoked with the context the caller passed to
	// Authenticate. It lets tests inspect context-propagated metadata such as the
	// brute-force-protection skip marker.
	OnAuthenticate func(ctx context.Context, in *gatewayv1beta1.AuthenticateRequest)
}

func (c mockGatewayClient) Authenticate(ctx context.Context, in *gatewayv1beta1.AuthenticateRequest, opts ...grpc.CallOption) (*gatewayv1beta1.AuthenticateResponse, error) {
	if c.OnAuthenticate != nil {
		c.OnAuthenticate(ctx, in)
	}
	response := c.AuthenticateFunc(in.GetType(), in.GetClientId(), in.GetClientSecret())
	return response, nil
}
