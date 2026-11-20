package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/stretchr/testify/assert"
	"go-micro.dev/v4/store"
)

var mfaTestUser = &userpb.User{Id: &userpb.UserId{OpaqueId: "user-1"}}

// mfaChain mirrors the proxy order: MultiFactor runs after the user is resolved.
func mfaChain(cfg config.MFAConfig, s store.Store, hasMFA *bool) http.Handler {
	resolveUser := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(revactx.ContextSetUser(r.Context(), mfaTestUser)))
		})
	}
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*hasMFA = revactx.HasMFA(r.Context())
	})

	return resolveUser(MultiFactor(cfg, Logger(log.NopLogger()), MFAStore(s))(final))
}

func oidcRequest(acr string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "https://example.com/dav/spaces/vault", nil)
	return req.WithContext(oidc.NewContext(req.Context(), map[string]any{"acr": acr}))
}

// signedURLRequest has no OIDC claims and no MFA-bound-key breadcrumb. It
// represents a claims-less request that was not authenticated with the
// MFA-bound signing key: a normal-key-signed URL, basic auth, or an app
// token.
func signedURLRequest() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "https://example.com/dav/spaces/vault?OC-Signature=abc", nil)
	return req.WithContext(revactx.ContextSetUser(req.Context(), mfaTestUser))
}

// mfaKeySignedURLRequest has no OIDC claims but does carry the MFA-bound
// signing-key breadcrumb, as set by the signed-URL authenticator when a
// signature matches the MFA-bound key.
func mfaKeySignedURLRequest() *http.Request {
	req := signedURLRequest()
	return req.WithContext(setMFAKeyBreadcrumb(req.Context()))
}

func TestMultiFactor_mfaKeySignedURLInheritsMFAFromOIDCRequestWhileRecordFresh(t *testing.T) {
	cfg := config.MFAConfig{Enabled: true, AuthLevelNames: []string{"advanced"}, SessionDuration: 60}
	s := store.NewMemoryStore()
	var hasMFA bool

	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), oidcRequest("advanced"))
	assert.True(t, hasMFA)

	records, err := s.Read(key(mfaTestUser.GetId().GetOpaqueId()))
	assert.NoError(t, err)
	assert.Len(t, records, 1)
	assert.Greater(t, records[0].Expiry, time.Duration(0))
	assert.LessOrEqual(t, records[0].Expiry, 60*time.Second)

	// A request carrying the MFA-bound-key breadcrumb inherits MFA from the
	// fresh stored record.
	hasMFA = false
	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), mfaKeySignedURLRequest())
	assert.True(t, hasMFA)
}

// TestMultiFactor_mfaKeySignedURLWithStaleRecordHasNoMFA proves that even a
// request carrying the MFA-bound-key breadcrumb does not inherit MFA once
// the stored per-user record is stale.
func TestMultiFactor_mfaKeySignedURLWithStaleRecordHasNoMFA(t *testing.T) {
	cfg := config.MFAConfig{Enabled: true, AuthLevelNames: []string{"advanced"}, SessionDuration: 60}
	s := store.NewMemoryStore()
	var hasMFA bool

	verifiedAt := time.Now().Add(-61 * time.Second).Unix()
	assert.NoError(t, s.Write(&store.Record{
		Key:   key(mfaTestUser.GetId().GetOpaqueId()),
		Value: []byte(strconv.FormatInt(verifiedAt, 10)),
	}))

	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), mfaKeySignedURLRequest())
	assert.False(t, hasMFA)
}

// TestMultiFactor_signedURLWithoutMFAKeyBreadcrumbNeverInheritsMFA is the
// regression test for OCISDEV-1699 / SECOPS-322 F-1.1: a claims-less
// request without the MFA-bound-key breadcrumb (e.g. a normal-key-signed
// URL, basic auth, or an app token) must never inherit MFA, even though the
// user's stored per-user record is fresh because another one of their
// sessions recently passed MFA.
func TestMultiFactor_signedURLWithoutMFAKeyBreadcrumbNeverInheritsMFA(t *testing.T) {
	cfg := config.MFAConfig{Enabled: true, AuthLevelNames: []string{"advanced"}, SessionDuration: 60}
	s := store.NewMemoryStore()
	var hasMFA bool

	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), oidcRequest("advanced"))
	assert.True(t, hasMFA)

	hasMFA = false
	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), signedURLRequest())
	assert.False(t, hasMFA)
}

// TestMultiFactor_mfaKeySignedVaultDownloadRegressionOCISDEV1503 guards
// against reopening OCISDEV-1503: a vault signed-URL download, signed with
// the MFA-bound key, must still succeed while the stored MFA record is
// fresh.
func TestMultiFactor_mfaKeySignedVaultDownloadRegressionOCISDEV1503(t *testing.T) {
	cfg := config.MFAConfig{Enabled: true, AuthLevelNames: []string{"advanced"}, SessionDuration: 60}
	s := store.NewMemoryStore()
	var hasMFA bool

	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), oidcRequest("advanced"))
	assert.True(t, hasMFA)

	hasMFA = false
	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), mfaKeySignedURLRequest())
	assert.True(t, hasMFA)
}

func TestMultiFactor_noMFAPersistedWithoutRequiredAcr(t *testing.T) {
	cfg := config.MFAConfig{Enabled: true, AuthLevelNames: []string{"advanced"}, SessionDuration: 60}
	s := store.NewMemoryStore()
	var hasMFA bool

	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), oidcRequest("regular"))
	assert.False(t, hasMFA)

	_, err := s.Read(key(mfaTestUser.GetId().GetOpaqueId()))
	assert.ErrorIs(t, err, store.ErrNotFound)

	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), signedURLRequest())
	assert.False(t, hasMFA)
}

func TestMultiFactor_signedURLWithoutPriorOIDCRequestHasNoMFA(t *testing.T) {
	cfg := config.MFAConfig{Enabled: true, AuthLevelNames: []string{"advanced"}, SessionDuration: 60}
	var hasMFA bool

	mfaChain(cfg, store.NewMemoryStore(), &hasMFA).ServeHTTP(httptest.NewRecorder(), signedURLRequest())
	assert.False(t, hasMFA)
}

func TestMultiFactor_expiredStoreEntryHasNoMFA(t *testing.T) {
	cfg := config.MFAConfig{Enabled: true, AuthLevelNames: []string{"advanced"}, SessionDuration: 60}
	s := store.NewMemoryStore()
	var hasMFA bool

	verifiedAt := time.Now().Add(-61 * time.Second).Unix()
	assert.NoError(t, s.Write(&store.Record{
		Key:   key(mfaTestUser.GetId().GetOpaqueId()),
		Value: []byte(strconv.FormatInt(verifiedAt, 10)),
	}))

	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), signedURLRequest())
	assert.False(t, hasMFA)
}

func TestMultiFactor_disabledDoesNotPersist(t *testing.T) {
	cfg := config.MFAConfig{Enabled: false, AuthLevelNames: []string{"advanced"}}
	s := store.NewMemoryStore()
	var hasMFA bool

	mfaChain(cfg, s, &hasMFA).ServeHTTP(httptest.NewRecorder(), oidcRequest("advanced"))
	// MultiFactor always sets MFA when disabled, but nothing needs to be persisted
	assert.True(t, hasMFA)

	_, err := s.Read(key(mfaTestUser.GetId().GetOpaqueId()))
	assert.ErrorIs(t, err, store.ErrNotFound)
}

func enforcedMFAConfig() config.MFAConfig {
	return config.MFAConfig{
		Enabled:                true,
		AuthLevelNames:         []string{"advanced"},
		SessionDuration:        60,
		EnforceSessionDuration: true,
		IssuedAtClaim:          "auth_time",
	}
}

func oidcRequestWithClaims(claims map[string]any) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "https://example.com/dav/spaces/vault", nil)
	return req.WithContext(oidc.NewContext(req.Context(), claims))
}

func TestMultiFactor_enforcedSessionDuration(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name   string
		claims map[string]any
		want   bool
	}{
		{name: "fresh", claims: map[string]any{"acr": "advanced", "auth_time": float64(now.Add(-30 * time.Second).Unix())}, want: true},
		{name: "expired", claims: map[string]any{"acr": "advanced", "auth_time": float64(now.Add(-61 * time.Second).Unix())}, want: false},
		{name: "missing", claims: map[string]any{"acr": "advanced"}, want: false},
		{name: "not a number", claims: map[string]any{"acr": "advanced", "auth_time": "yesterday"}, want: false},
		{name: "zero", claims: map[string]any{"acr": "advanced", "auth_time": float64(0)}, want: false},
		{name: "in the future", claims: map[string]any{"acr": "advanced", "auth_time": float64(now.Add(time.Minute).Unix())}, want: false},
		{name: "fresh without mfa acr", claims: map[string]any{"acr": "regular", "auth_time": float64(now.Unix())}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var hasMFA bool
			mfaChain(enforcedMFAConfig(), store.NewMemoryStore(), &hasMFA).ServeHTTP(httptest.NewRecorder(), oidcRequestWithClaims(tt.claims))
			assert.Equal(t, tt.want, hasMFA)
		})
	}
}

func TestMultiFactor_customIssuedAtClaim(t *testing.T) {
	cfg := enforcedMFAConfig()
	cfg.IssuedAtClaim = "loa_times.advanced"
	var hasMFA bool

	claims := map[string]any{
		"acr":       "advanced",
		"auth_time": float64(time.Now().Add(-time.Hour).Unix()),
		"loa_times": map[string]any{"advanced": float64(time.Now().Add(-10 * time.Second).Unix())},
	}
	mfaChain(cfg, store.NewMemoryStore(), &hasMFA).ServeHTTP(httptest.NewRecorder(), oidcRequestWithClaims(claims))
	assert.True(t, hasMFA)
}

func TestMultiFactor_enforcementDisabledIgnoresMFATime(t *testing.T) {
	cfg := enforcedMFAConfig()
	cfg.EnforceSessionDuration = false
	var hasMFA bool

	claims := map[string]any{"acr": "advanced", "auth_time": float64(time.Now().Add(-24 * time.Hour).Unix())}
	mfaChain(cfg, store.NewMemoryStore(), &hasMFA).ServeHTTP(httptest.NewRecorder(), oidcRequestWithClaims(claims))
	assert.True(t, hasMFA)
}

// TestMultiFactor_storeFallbackDoesNotOutliveMFASession proves that a signed-URL download can't
// extend MFA beyond the original session: the store record holds the MFA time from the claim.
func TestMultiFactor_storeFallbackDoesNotOutliveMFASession(t *testing.T) {
	s := store.NewMemoryStore()
	var hasMFA bool

	verifiedAt := time.Now().Add(-50 * time.Second)
	claims := map[string]any{"acr": "advanced", "auth_time": float64(verifiedAt.Unix())}
	mfaChain(enforcedMFAConfig(), s, &hasMFA).ServeHTTP(httptest.NewRecorder(), oidcRequestWithClaims(claims))
	assert.True(t, hasMFA)

	records, err := s.Read(key(mfaTestUser.GetId().GetOpaqueId()))
	assert.NoError(t, err)
	assert.Len(t, records, 1)
	assert.Equal(t, strconv.FormatInt(verifiedAt.Unix(), 10), string(records[0].Value))
	assert.LessOrEqual(t, records[0].Expiry, 10*time.Second)
}
