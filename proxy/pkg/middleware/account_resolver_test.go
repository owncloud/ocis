package middleware

import (
	"context"
	"github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/pkg/token"
	"github.com/owncloud/ocis/ocis-pkg/log"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"github.com/owncloud/ocis/proxy/pkg/config"
	"github.com/owncloud/ocis/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/proxy/pkg/user/backend/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTokenIsAddedWithMailClaim(t *testing.T) {
	sut := newMockAccountResolver(&userv1beta1.User{
		Id:   &userv1beta1.UserId{Idp: "https://idx.example.com", OpaqueId: "123"},
		Mail: "foo@example.com",
	}, nil)

	req, rw := mockRequest(&oidc.StandardClaims{
		Iss:   "https://idx.example.com",
		Email: "foo@example.com",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(token.TokenHeader)
	assert.NotEmpty(t, token)
	assert.Contains(t, token, "eyJ")
}

func TestTokenIsAddedWithUsernameClaim(t *testing.T) {
	sut := newMockAccountResolver(&userv1beta1.User{
		Id:   &userv1beta1.UserId{Idp: "https://idx.example.com", OpaqueId: "123"},
		Mail: "foo@example.com",
	}, nil)

	req, rw := mockRequest(&oidc.StandardClaims{
		Iss:               "https://idx.example.com",
		PreferredUsername: "foo",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(token.TokenHeader)
	assert.NotEmpty(t, token)

	assert.Contains(t, token, "eyJ")
}

func TestNSkipOnNoClaims(t *testing.T) {
	sut := newMockAccountResolver(nil, backend.ErrAccountDisabled)
	req, rw := mockRequest(nil)

	sut.ServeHTTP(rw, req)

	token := req.Header.Get("x-access-token")
	assert.Empty(t, token)
	assert.Equal(t, http.StatusOK, rw.Code)
}

func TestUnauthorizedOnUserNotFound(t *testing.T) {
	sut := newMockAccountResolver(nil, backend.ErrAccountNotFound)
	req, rw := mockRequest(&oidc.StandardClaims{
		Iss:               "https://idx.example.com",
		PreferredUsername: "foo",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(token.TokenHeader)
	assert.Empty(t, token)
	assert.Equal(t, http.StatusUnauthorized, rw.Code)
}

func TestUnauthorizedOnUserDisabled(t *testing.T) {
	sut := newMockAccountResolver(nil, backend.ErrAccountDisabled)
	req, rw := mockRequest(&oidc.StandardClaims{
		Iss:               "https://idx.example.com",
		PreferredUsername: "foo",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(token.TokenHeader)
	assert.Empty(t, token)
	assert.Equal(t, http.StatusUnauthorized, rw.Code)
}

func TestInternalServerErrorOnMissingMailAndUsername(t *testing.T) {
	sut := newMockAccountResolver(nil, backend.ErrAccountDisabled)
	req, rw := mockRequest(&oidc.StandardClaims{
		Iss: "https://idx.example.com",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(token.TokenHeader)
	assert.Empty(t, token)
	assert.Equal(t, http.StatusInternalServerError, rw.Code)
}

func newMockAccountResolver(userBackendResult *userv1beta1.User, userBackendErr error) http.Handler {
	mock := &test.UserBackendMock{
		GetUserByClaimsFunc: func(ctx context.Context, claim string, value string, withRoles bool) (*userv1beta1.User, error) {
			return userBackendResult, userBackendErr
		},
	}

	return AccountResolver(
		Logger(log.NewLogger()),
		UserProvider(mock),
		TokenManagerConfig(config.TokenManager{JWTSecret: "secret"}),
		AutoprovisionAccounts(false),
	)(mockHandler{})
}

func mockRequest(claims *oidc.StandardClaims) (*http.Request, *httptest.ResponseRecorder) {
	if claims == nil {
		return httptest.NewRequest("GET", "http://example.com/foo", nil), httptest.NewRecorder()
	}

	ctx := oidc.NewContext(context.Background(), claims)
	req := httptest.NewRequest("GET", "http://example.com/foo", nil).WithContext(ctx)
	rw := httptest.NewRecorder()

	return req, rw
}

type mockHandler struct{}

func (m mockHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {}
