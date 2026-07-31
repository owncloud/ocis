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

func TestResolveUserType(t *testing.T) {
	const (
		instanceID1 = "ec730a6c-1b63-4b45-b83b-9e2311afdf85"
		instanceID2 = "8d24cb5f-6ee6-4b98-86df-c4c268dddb46"
		masterID    = "11111111-1111-1111-1111-111111111111"
	)

	testCases := []struct {
		name           string
		multiInstance  bool
		instanceID     string
		masterID       string
		memberClaim    string
		guestClaim     string
		claims         map[string]interface{}
		expectMember   bool
		expectGuest    bool
		description    string
	}{
		{
			name:          "single instance mode - always member",
			multiInstance: false,
			instanceID:    instanceID1,
			masterID:      "",
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims:        map[string]interface{}{},
			expectMember:  true,
			expectGuest:   false,
			description:   "In single instance mode, any user should be considered a member",
		},
		{
			name:          "multi-instance: user has matching instance ID in memberClaim",
			multiInstance: true,
			instanceID:    instanceID1,
			masterID:      "",
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims: map[string]interface{}{
				"owncloudMemberOf": instanceID1,
			},
			expectMember: true,
			expectGuest:  false,
			description:  "User with matching instance ID should be granted member access",
		},
		{
			name:          "multi-instance: user has matching instance ID in guestClaim",
			multiInstance: true,
			instanceID:    instanceID1,
			masterID:      "",
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims: map[string]interface{}{
				"owncloudGuestOf": instanceID1,
			},
			expectMember: false,
			expectGuest:  true,
			description:  "User with matching instance ID in guestClaim should be granted guest access",
		},
		{
			name:          "multi-instance: user has no matching instance ID",
			multiInstance: true,
			instanceID:    instanceID1,
			masterID:      "",
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims: map[string]interface{}{
				"owncloudMemberOf": instanceID2,
			},
			expectMember: false,
			expectGuest:  false,
			description:  "User without matching instance ID should be denied access",
		},
		{
			name:          "multi-instance with master-id: user has master-id in memberClaim",
			multiInstance: true,
			instanceID:    instanceID1,
			masterID:      masterID,
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims: map[string]interface{}{
				"owncloudMemberOf": masterID,
			},
			expectMember: true,
			expectGuest:  false,
			description:  "User with master-id in memberClaim should be granted member access to any instance",
		},
		{
			name:          "multi-instance with master-id disabled: user has master-id value but feature disabled",
			multiInstance: true,
			instanceID:    instanceID1,
			masterID:      "", // Empty = disabled
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims: map[string]interface{}{
				"owncloudMemberOf": masterID,
			},
			expectMember: false,
			expectGuest:  false,
			description:  "When master-id is empty/disabled, having the master-id value should not grant access",
		},
		{
			name:          "multi-instance with master-id: user has both master-id and instance-id",
			multiInstance: true,
			instanceID:    instanceID1,
			masterID:      masterID,
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims: map[string]interface{}{
				"owncloudMemberOf": []string{masterID, instanceID1},
			},
			expectMember: true,
			expectGuest:  false,
			description:  "User with both master-id and instance-id should be granted access (master-id matches first)",
		},
		{
			name:          "multi-instance with master-id: array with master-id in middle",
			multiInstance: true,
			instanceID:    instanceID1,
			masterID:      masterID,
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims: map[string]interface{}{
				"owncloudMemberOf": []string{instanceID2, masterID, "other-id"},
			},
			expectMember: true,
			expectGuest:  false,
			description:  "Master-id should be found when it's in an array with other values",
		},
		{
			name:          "multi-instance with master-id: user on wrong instance but has master-id",
			multiInstance: true,
			instanceID:    instanceID2, // Different instance
			masterID:      masterID,
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims: map[string]interface{}{
				"owncloudMemberOf": masterID,
			},
			expectMember: true,
			expectGuest:  false,
			description:  "Master-id should grant access to ANY instance, regardless of which instance is checking",
		},
		{
			name:          "multi-instance: empty claim values",
			multiInstance: true,
			instanceID:    instanceID1,
			masterID:      masterID,
			memberClaim:   "owncloudMemberOf",
			guestClaim:    "owncloudGuestOf",
			claims: map[string]interface{}{
				"owncloudMemberOf": "",
			},
			expectMember: false,
			expectGuest:  false,
			description:  "Empty string values should not grant access",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resolver := accountResolver{
				multiInstance: tc.multiInstance,
				instanceID:    tc.instanceID,
				masterID:      tc.masterID,
				memberClaim:   tc.memberClaim,
				guestClaim:    tc.guestClaim,
			}

			isMember, isGuest := resolver.resolveUserType(tc.claims)

			assert.Equal(t, tc.expectMember, isMember, "isMember mismatch: %s", tc.description)
			assert.Equal(t, tc.expectGuest, isGuest, "isGuest mismatch: %s", tc.description)
		})
	}
}
