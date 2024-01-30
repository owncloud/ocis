package middleware

import (
	"net/http"
	"net/http/httptest"
	"time"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/golang-jwt/jwt/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	oidcmocks "github.com/owncloud/ocis/v2/ocis-pkg/oidc/mocks"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	ubmocks "github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend/mocks"
	"github.com/stretchr/testify/mock"
	"go-micro.dev/v4/store"
)

var _ = Describe("Authenticating requests", Label("OIDCAuthenticator"), func() {
	var authenticator Authenticator
	ub := ubmocks.UserBackend{}
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
	/*
		// to test with skipUserInfo:  true, we need to also use an interface so we can mock the UserInfo.Claim call
		oc.On("UserInfo", mock.Anything, mock.Anything).Return(
			&oidc.UserInfo{
				Subject:       "my-sub",
				EmailVerified: true,
				Email:         "test@example.org",
			},
			nil,
		)
	*/

	BeforeEach(func() {
		authenticator = &OIDCAuthenticator{
			OIDCIss:       "http://idp.example.com",
			Logger:        log.NewLogger(),
			oidcClient:    &oc,
			userInfoCache: store.NewMemoryStore(),
			skipUserInfo:  true,
		}
	})

	When("the request contains correct data", func() {
		It("should successfully authenticate", func() {
			req := httptest.NewRequest(http.MethodGet, "http://example.com/example/path", http.NoBody)
			req.Header.Set(_headerAuthorization, "Bearer jwt.token.sig")

			req2, valid := authenticator.Authenticate(req)

			Expect(valid).To(Equal(true))
			Expect(req2).ToNot(BeNil())
		})
	})
})
