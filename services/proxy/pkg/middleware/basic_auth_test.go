package middleware

import (
	"net/http"
	"net/http/httptest"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend/mocks"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("Authenticating requests", Label("BasicAuthenticator"), func() {
	var authenticator Authenticator
	BeforeEach(func() {
		var user userv1beta1.User
		var err error
		backendMock := &mocks.UserBackend{}
		backendMock.On("Authenticate", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(
			func(args mock.Arguments) {
				if args.String(1) == "testuser" && args.String(2) == "testpassword" {
					user = userv1beta1.User{
						Id: &userv1beta1.UserId{
							Idp:      "IdpId",
							OpaqueId: "OpaqueId",
						},
						Username: "testuser",
						Mail:     "testuser@example.com",
					}
				} else {
					err = backend.ErrAccountNotFound
				}
			},
		).Return(&user, "", err)

		authenticator = BasicAuthenticator{
			Logger:       log.NewLogger(),
			UserProvider: backendMock,
		}
	})

	When("the request contains correct data", func() {
		It("should successfully authenticate", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.SetBasicAuth("testuser", "testpassword")

			req2, valid := authenticator.Authenticate(req)

			Expect(valid).To(Equal(true))
			Expect(req2).ToNot(BeNil())
		})
		It("adds claims to the request context", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.SetBasicAuth("testuser", "testpassword")

			req2, valid := authenticator.Authenticate(req)
			Expect(valid).To(Equal(true))

			claims := oidc.FromContext(req2.Context())
			Expect(claims).ToNot(BeNil())
			Expect(claims[oidc.Iss]).To(Equal("IdpId"))
			Expect(claims[oidc.PreferredUsername]).To(Equal("testuser"))
			Expect(claims[oidc.Email]).To(Equal("testuser@example.com"))
			Expect(claims[oidc.OwncloudUUID]).To(Equal("OpaqueId"))
		})
	})
})
