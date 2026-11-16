package staticroutes

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	oidcmocks "github.com/owncloud/ocis/v2/ocis-pkg/oidc/mocks"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	backendmocks "github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend/mocks"
	"github.com/shamaton/msgpack/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go-micro.dev/v4/store"
)

// backchannelLogoutTestUser is the user resolved from the logout token's
// session in all tests below.
var backchannelLogoutTestUser = &userpb.User{Id: &userpb.UserId{OpaqueId: "user-1"}}

// newBackchannelLogoutHandler builds a StaticRouteHandler wired with memory
// stores, seeds the userinfo cache so the handler can resolve
// backchannelLogoutTestUser from the session, and seeds the MFA status
// record and MFA-bound signing key for that user so their deletion is
// observable.
func newBackchannelLogoutHandler(t *testing.T) (*StaticRouteHandler, string) {
	t.Helper()

	const sessionID = "session-1"
	const userinfoKey = "userinfo-key-1"

	userInfoCache := store.NewMemoryStore()
	require.NoError(t, userInfoCache.Write(&store.Record{Key: sessionID, Value: []byte(userinfoKey)}))

	claims, err := msgpack.MarshalAsMap(map[string]any{"oidc-claim": "the-claim-value"})
	require.NoError(t, err)
	require.NoError(t, userInfoCache.Write(&store.Record{Key: userinfoKey, Value: claims}))

	mfaStore := store.NewMemoryStore()
	require.NoError(t, mfaStore.Write(&store.Record{Key: mfaRecordKey(backchannelLogoutTestUser.GetId().GetOpaqueId()), Value: []byte("1700000000")}))

	signingKeyStore := store.NewMemoryStore()
	require.NoError(t, signingKeyStore.Write(&store.Record{Key: mfaSigningKeyName(backchannelLogoutTestUser.GetId().GetOpaqueId()), Value: []byte("some-mfa-bound-key")}))

	oidcClient := &oidcmocks.OIDCClient{}
	oidcClient.On("VerifyLogoutToken", mock.Anything, mock.Anything).Return(&oidc.LogoutToken{SessionId: sessionID}, nil)

	userProvider := &backendmocks.UserBackend{}
	userProvider.On("GetUserByClaims", mock.Anything, "userid", "the-claim-value").Return(backchannelLogoutTestUser, "", nil)

	handler := &StaticRouteHandler{
		Logger: log.NopLogger(),
		Config: config.Config{
			UserOIDCClaim: "oidc-claim",
			UserCS3Claim:  "userid",
		},
		OidcClient:      oidcClient,
		UserInfoCache:   userInfoCache,
		UserProvider:    userProvider,
		MFAStore:        mfaStore,
		SigningKeyStore: signingKeyStore,
	}

	return handler, backchannelLogoutTestUser.GetId().GetOpaqueId()
}

func postBackchannelLogout(t *testing.T, handler *StaticRouteHandler) *httptest.ResponseRecorder {
	t.Helper()

	form := url.Values{"logout_token": {"irrelevant-because-VerifyLogoutToken-is-mocked"}}
	req := httptest.NewRequest(http.MethodPost, "https://example.com/backchannel_logout", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.backchannelLogout(w, req)
	return w
}

func TestBackchannelLogout_deletesMFAStatusAndSigningKey(t *testing.T) {
	handler, userID := newBackchannelLogoutHandler(t)

	w := postBackchannelLogout(t, handler)

	assert.Equal(t, http.StatusOK, w.Code)

	records, err := handler.MFAStore.Read(mfaRecordKey(userID))
	assert.True(t, err != nil || len(records) == 0, "expected the per-user MFA status record to be deleted at backchannel logout")

	records, err = handler.SigningKeyStore.Read(mfaSigningKeyName(userID))
	assert.True(t, err != nil || len(records) == 0, "expected the MFA-bound signing key to be deleted at backchannel logout")
}

func TestBackchannelLogout_stillCleansUpUserInfoCache(t *testing.T) {
	handler, _ := newBackchannelLogoutHandler(t)

	w := postBackchannelLogout(t, handler)

	assert.Equal(t, http.StatusOK, w.Code)

	records, err := handler.UserInfoCache.Read("session-1")
	assert.True(t, err != nil || len(records) == 0, "expected the session lookup entry to still be cleaned up from the userinfo cache")

	records, err = handler.UserInfoCache.Read("userinfo-key-1")
	assert.True(t, err != nil || len(records) == 0, "expected the userinfo record to still be cleaned up from the cache")
}
