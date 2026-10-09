package middleware

import (
	"net/http"
	"net/http/httptest"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend/mocks"
)

var _ = Describe("Authenticating requests", Label("BasicAuthenticator"), func() {
	var authenticator Authenticator
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
		authenticator = BasicAuthenticator{
			Logger:       log.NewLogger(),
			UserProvider: &ub,
		}
	})

	// FIXME: test currently commented out due to missing roleManager / roleService mock
	// It must be restored and adjusted once the all the mock services are up to date
	//When("the request contains correct data", func() {
	//	It("should successfully authenticate", func() {
	//		req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
	//		req.SetBasicAuth("testuser", "testpassword")

	//		req2, err := authenticator.Authenticate(req)

	//		Expect(err).ToNot(HaveOccurred())
	//		Expect(req2).ToNot(BeNil())
	//	})
	//	It("adds claims to the request context", func() {
	//		req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
	//		req.SetBasicAuth("testuser", "testpassword")

	//		req2, err := authenticator.Authenticate(req)
	//		Expect(err).ToNot(HaveOccurred())

	//		claims := oidc.FromContext(req2.Context())
	//		Expect(claims).ToNot(BeNil())
	//		Expect(claims[oidc.Iss]).To(Equal("IdpId"))
	//		Expect(claims[oidc.PreferredUsername]).To(Equal("testuser"))
	//		Expect(claims[oidc.Email]).To(Equal("testuser@example.com"))
	//		Expect(claims[oidc.OwncloudUUID]).To(Equal("OpaqueId"))
	//	})
	//})
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
