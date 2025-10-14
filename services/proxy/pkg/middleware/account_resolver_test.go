package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend/mocks"
	userRoleMocks "github.com/owncloud/ocis/v2/services/proxy/pkg/userroles/mocks"
	"github.com/owncloud/reva/v2/pkg/auth/scope"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/token/manager/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTokenIsAddedWithMailClaim(t *testing.T) {
	sut := newMockAccountResolver(&userv1beta1.User{
		Id:   &userv1beta1.UserId{Idp: "https://idx.example.com", OpaqueId: "123"},
		Mail: "foo@example.com",
	}, nil, oidc.Email, "mail")

	req, rw := mockRequest(map[string]interface{}{
		oidc.Iss:   "https://idx.example.com",
		oidc.Email: "foo@example.com",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(revactx.TokenHeader)
	assert.NotEmpty(t, token)
	assert.Contains(t, token, "eyJ")
}

func TestTokenIsAddedWithUsernameClaim(t *testing.T) {
	sut := newMockAccountResolver(&userv1beta1.User{
		Id:   &userv1beta1.UserId{Idp: "https://idx.example.com", OpaqueId: "123"},
		Mail: "foo@example.com",
	}, nil, oidc.PreferredUsername, "username")

	req, rw := mockRequest(map[string]interface{}{
		oidc.Iss:               "https://idx.example.com",
		oidc.PreferredUsername: "foo",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(revactx.TokenHeader)
	assert.NotEmpty(t, token)

	assert.Contains(t, token, "eyJ")
}

func TestTokenIsAddedWithDotUsernamePathClaim(t *testing.T) {
	sut := newMockAccountResolver(&userv1beta1.User{
		Id:   &userv1beta1.UserId{Idp: "https://idx.example.com", OpaqueId: "123"},
		Mail: "foo@example.com",
	}, nil, "li.un", "username")

	// This is how lico adds the username to the access token
	req, rw := mockRequest(map[string]interface{}{
		oidc.Iss: "https://idx.example.com",
		"li": map[string]interface{}{
			"un": "foo",
		},
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(revactx.TokenHeader)
	assert.NotEmpty(t, token)

	assert.Contains(t, token, "eyJ")
}

func TestTokenIsAddedWithDotEscapedUsernameClaim(t *testing.T) {
	sut := newMockAccountResolver(&userv1beta1.User{
		Id:   &userv1beta1.UserId{Idp: "https://idx.example.com", OpaqueId: "123"},
		Mail: "foo@example.com",
	}, nil, "li\\.un", "username")

	// This tests the . escaping of the readUserIDClaim
	req, rw := mockRequest(map[string]interface{}{
		oidc.Iss: "https://idx.example.com",
		"li.un":  "foo",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(revactx.TokenHeader)
	assert.NotEmpty(t, token)

	assert.Contains(t, token, "eyJ")
}

func TestTokenIsAddedWithDottedUsernameClaimFallback(t *testing.T) {
	sut := newMockAccountResolver(&userv1beta1.User{
		Id:   &userv1beta1.UserId{Idp: "https://idx.example.com", OpaqueId: "123"},
		Mail: "foo@example.com",
	}, nil, "li.un", "username")

	// This tests the . escaping fallback of the readUserIDClaim
	req, rw := mockRequest(map[string]interface{}{
		oidc.Iss: "https://idx.example.com",
		"li.un":  "foo",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(revactx.TokenHeader)
	assert.NotEmpty(t, token)

	assert.Contains(t, token, "eyJ")
}

func TestNSkipOnNoClaims(t *testing.T) {
	sut := newMockAccountResolver(nil, backend.ErrAccountDisabled, oidc.Email, "mail")
	req, rw := mockRequest(nil)

	sut.ServeHTTP(rw, req)

	token := req.Header.Get("x-access-token")
	assert.Empty(t, token)
	assert.Equal(t, http.StatusOK, rw.Code)
}

func TestUnauthorizedOnUserNotFound(t *testing.T) {
	sut := newMockAccountResolver(nil, backend.ErrAccountNotFound, oidc.PreferredUsername, "username")
	req, rw := mockRequest(map[string]interface{}{
		oidc.Iss:               "https://idx.example.com",
		oidc.PreferredUsername: "foo",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(revactx.TokenHeader)
	assert.Empty(t, token)
	assert.Equal(t, http.StatusUnauthorized, rw.Code)
}

func TestUnauthorizedOnUserDisabled(t *testing.T) {
	sut := newMockAccountResolver(nil, backend.ErrAccountDisabled, oidc.PreferredUsername, "username")
	req, rw := mockRequest(map[string]interface{}{
		oidc.Iss:               "https://idx.example.com",
		oidc.PreferredUsername: "foo",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(revactx.TokenHeader)
	assert.Empty(t, token)
	assert.Equal(t, http.StatusUnauthorized, rw.Code)
}

func TestInternalServerErrorOnMissingMailAndUsername(t *testing.T) {
	sut := newMockAccountResolver(nil, backend.ErrAccountNotFound, oidc.Email, "mail")
	req, rw := mockRequest(map[string]interface{}{
		oidc.Iss: "https://idx.example.com",
	})

	sut.ServeHTTP(rw, req)

	token := req.Header.Get(revactx.TokenHeader)
	assert.Empty(t, token)
	assert.Equal(t, http.StatusInternalServerError, rw.Code)
}

func newMockAccountResolver(userBackendResult *userv1beta1.User, userBackendErr error, oidcclaim, cs3claim string) http.Handler {
	tokenManager, _ := jwt.New(map[string]interface{}{
		"secret":  "change-me",
		"expires": int64(60),
	})

	token := ""
	if userBackendResult != nil {
		s, _ := scope.AddOwnerScope(nil)
		token, _ = tokenManager.MintToken(context.Background(), userBackendResult, s)
	}

	ub := mocks.UserBackend{}
	ub.On("GetUserByClaims", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(userBackendResult, token, userBackendErr)
	ub.On("GetUserRoles", mock.Anything, mock.Anything).Return(userBackendResult, nil)

	ra := userRoleMocks.UserRoleAssigner{}
	ra.On("UpdateUserRoleAssignment", mock.Anything, mock.Anything, mock.Anything).Return(userBackendResult, nil)

	return AccountResolver(
		Logger(log.NewLogger()),
		UserProvider(&ub),
		UserRoleAssigner(&ra),
		SkipUserInfo(false),
		UserOIDCClaim(oidcclaim),
		UserCS3Claim(cs3claim),
		AutoprovisionAccounts(false),
	)(mockHandler{})
}

func mockRequest(claims map[string]interface{}) (*http.Request, *httptest.ResponseRecorder) {
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
