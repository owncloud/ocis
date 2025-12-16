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
	ra.On("UpdateUserRoleAssignment", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(userBackendResult, nil)

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

func TestReadClaim(t *testing.T) {
	var testCases = []struct {
		Alias    string
		Key      string
		Claims   map[string]any
		Expected []string
	}{
		{
			Alias: "single value",
			Key:   "testkey",
			Claims: map[string]any{
				"testkey": "testvalue",
			},
			Expected: []string{"testvalue"},
		},
		{
			Alias: "multivalue",
			Key:   "testkey",
			Claims: map[string]any{
				"testkey": []string{"testvalue1", "testvalue2"},
			},
			Expected: []string{"testvalue1", "testvalue2"},
		},
		{
			Alias: "empty value 1",
			Key:   "testkey",
			Claims: map[string]any{
				"testkey": "",
			},
			Expected: []string{""},
		},
		{
			Alias: "empty value 2",
			Key:   "testkey",
			Claims: map[string]any{
				"testkey": []string{},
			},
			Expected: []string{},
		},
		{
			Alias: "no value",
			Key:   "testkey",
			Claims: map[string]any{
				"testkey": nil,
			},
			Expected: nil,
		},
		{
			Alias: "no key",
			Key:   "testkey",
			Claims: map[string]any{
				"anotherkey": "withvalue",
			},
			Expected: nil,
		},
		{
			Alias: "wrong type 1",
			Key:   "testkey",
			Claims: map[string]any{
				"anotherkey": true,
			},
			Expected: nil,
		},
		{
			Alias: "wrong type 2",
			Key:   "testkey",
			Claims: map[string]any{
				"anotherkey": 123,
			},
			Expected: nil,
		},
	}

	for _, tc := range testCases {
		s := readClaim(tc.Key, tc.Claims)
		assert.Equal(t, tc.Expected, s)
	}

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
