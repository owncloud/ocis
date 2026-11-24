package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/mfa"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	backendmocks "github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend/mocks"
	userrolesmocks "github.com/owncloud/ocis/v2/services/proxy/pkg/userroles/mocks"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go-micro.dev/v4/store"
)

func TestSignedURLAuth_shouldServe(t *testing.T) {
	pua := SignedURLAuthenticator{}
	tests := []struct {
		url      string
		enabled  bool
		expected bool
	}{
		{"https://example.com/example.jpg", true, false},
		{"https://example.com/example.jpg?OC-Signature=something", true, true},
		{"https://example.com/example.jpg", false, false},
		{"https://example.com/example.jpg?OC-Signature=something", false, false},
	}

	for _, tt := range tests {
		pua.PreSignedURLConfig.Enabled = tt.enabled
		r := httptest.NewRequest("", tt.url, nil)
		result := pua.shouldServe(r)

		if result != tt.expected {
			t.Errorf("with %s expected %t got %t", tt.url, tt.expected, result)
		}
	}
}

func TestSignedURLAuth_allRequiredParametersPresent(t *testing.T) {
	pua := SignedURLAuthenticator{}
	baseURL := "https://example.com/example.jpg?"
	tests := []struct {
		params       string
		errorMessage string
	}{
		{"OC-Signature=something&OC-Credential=something&OC-Date=something&OC-Expires=something&OC-Verb=something", ""},
		{"OC-Credential=something&OC-Date=something&OC-Expires=something&OC-Verb=something", "required OC-Signature parameter not found"},
		{"OC-Signature=something&OC-Date=something&OC-Expires=something&OC-Verb=something", "required OC-Credential parameter not found"},
		{"OC-Signature=something&OC-Credential=something&OC-Expires=something&OC-Verb=something", "required OC-Date parameter not found"},
		{"OC-Signature=something&OC-Credential=something&OC-Date=something&OC-Verb=something", "required OC-Expires parameter not found"},
		{"OC-Signature=something&OC-Credential=something&OC-Date=something&OC-Expires=something", "required OC-Verb parameter not found"},
	}
	for _, tt := range tests {
		r := httptest.NewRequest("", baseURL+tt.params, nil)
		err := pua.allRequiredParametersArePresent(r.URL.Query())
		if tt.errorMessage != "" {
			assert.EqualError(t, err, tt.errorMessage, tt.params)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestSignedURLAuth_requestMethodMatches(t *testing.T) {
	pua := SignedURLAuthenticator{}
	tests := []struct {
		method       string
		url          string
		errorMessage string
	}{
		{"GET", "https://example.com/example.jpg?OC-Verb=GET", ""},
		{"GET", "https://example.com/example.jpg?OC-Verb=get", ""},
		{"POST", "https://example.com/example.jpg?OC-Verb=GET", "required OC-Verb parameter did not match request method"},
	}

	for _, tt := range tests {
		r := httptest.NewRequest(tt.method, tt.url, nil)
		err := pua.requestMethodMatches(r.Method, r.URL.Query())
		if tt.errorMessage != "" {
			assert.EqualError(t, err, tt.errorMessage, tt.url)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestSignedURLAuth_requestMethodIsAllowed(t *testing.T) {
	pua := SignedURLAuthenticator{}
	tests := []struct {
		method       string
		allowed      []string
		errorMessage string
	}{
		{"GET", []string{}, "request method is not listed in PreSignedURLConfig AllowedHTTPMethods"},
		{"GET", []string{"POST"}, "request method is not listed in PreSignedURLConfig AllowedHTTPMethods"},
		{"GET", []string{"GET"}, ""},
		{"GET", []string{"get"}, ""},
		{"GET", []string{"POST", "GET"}, ""},
	}

	for _, tt := range tests {
		pua.PreSignedURLConfig.AllowedHTTPMethods = tt.allowed
		err := pua.requestMethodIsAllowed(tt.method)
		if tt.errorMessage != "" {
			assert.EqualError(t, err, tt.errorMessage)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestSignedURLAuth_urlIsExpired(t *testing.T) {
	nowFunc := func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2020-02-02T12:30:00.000Z")
		return t
	}
	pua := SignedURLAuthenticator{
		Now: nowFunc,
	}

	tests := []struct {
		url          string
		errorMessage string
	}{
		// a valid signed url
		{"http://example.com/example.jpg?OC-Date=2020-02-02T12:29:00.000Z&OC-Expires=61", ""},
		// invalid expiry
		{"http://example.com/example.jpg?OC-Date=2020-02-02T12:29:00.000Z&OC-Expires=invalid", "time: invalid duration \"invalids\""},
		// wrong date format on OC-Date
		{"http://example.com/example.jpg?OC-Date=2020-02-02TTT12:29:00.000Z&OC-Expires=5", "parsing time \"2020-02-02TTT12:29:00.000Z\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"TT12:29:00.000Z\" as \"15\""},
		// expired - 12:29:00 + 59s < 12:30
		{"http://example.com/example.jpg?OC-Date=2020-02-02T12:29:00.000Z&OC-Expires=59", "URL is expired"},
		// expired - basically url was created yesterday
		{"http://example.com/example.jpg?OC-Date=2020-02-01T12:29:00.000Z&OC-Expires=59", "URL is expired"},
		// future OC-Date - also valid now
		{"http://example.com/example.jpg?OC-Date=2020-02-03T12:29:00.000Z&OC-Expires=59", ""},
	}

	for _, tt := range tests {
		r := httptest.NewRequest("", tt.url, nil)
		err := pua.urlIsExpired(r.URL.Query())
		if tt.errorMessage != "" {
			assert.EqualError(t, err, tt.errorMessage, tt.url)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestSignedURLAuth_createSignature(t *testing.T) {
	pua := SignedURLAuthenticator{}
	expected := "27d2ebea381384af3179235114801dcd00f91e46f99fca72575301cf3948101d"
	s := pua.createSignature("something", []byte("somerandomkey"))

	if s != expected {
		t.Fail()
	}
}

// TestSignedURLAuth_validateEmptySigningKeyStore guards against an out-of-range
// panic when the signing-key store returns no record without an error. This
// happens with the "noop" store (e.g. when OCIS_CACHE_STORE=noop, which also
// governs the presigned-URL signing-key store). The request must fail
// authentication, not panic.
func TestSignedURLAuth_validateEmptySigningKeyStore(t *testing.T) {
	pua := SignedURLAuthenticator{
		PreSignedURLConfig: config.PreSignedURL{
			AllowedHTTPMethods: []string{"get"},
			Enabled:            true,
		},
		Store:  store.NewNoopStore(),
		Logger: log.NewLogger(),
		Now: func() time.Time {
			tm, _ := time.Parse(time.RFC3339, "2019-05-14T11:02:00.000Z")
			return tm
		},
	}

	u := userpb.User{Id: &userpb.UserId{OpaqueId: "useri"}, DisplayName: "Test User"}
	ctx := revactx.ContextSetUser(context.Background(), &u)
	r := httptest.NewRequest("", "http://cloud.example.net/?OC-Credential=alice&OC-Date=2019-05-14T11%3A01%3A58.135Z&OC-Expires=1200&OC-Verb=GET&OC-Signature=f9e53a1ee23caef10f72ec392c1b537317491b687bfdd224c782be197d9ca2b6", nil).WithContext(ctx)

	var err error
	assert.NotPanics(t, func() { err = pua.validate(r) }, "validate must not panic when the signing key store is empty")
	assert.Error(t, err, "validate must reject the request when no signing key is found")
}

func TestSignedURLAuth_validate(t *testing.T) {
	nowFunc := func() time.Time {
		t, _ := time.Parse(time.RFC3339, "2020-02-02T12:30:00.000Z")
		return t
	}
	cfg := config.PreSignedURL{
		AllowedHTTPMethods: []string{"get"},
		Enabled:            true,
	}
	pua := SignedURLAuthenticator{
		PreSignedURLConfig: cfg,
		Store:              store.NewMemoryStore(),
		Now:                nowFunc,
	}

	pua.Store.Write(&store.Record{
		Key:      "useri",
		Value:    []byte("1234567890"),
		Metadata: nil,
	})

	tests := []struct {
		now          string
		url          string
		errorMessage string
	}{
		{"2020-02-02T12:30:00.000Z", "http://example.com/example.jpg?OC-Date=2020-02-02T12:29:00.000Z&OC-Expires=invalid", "required OC-Signature parameter not found"},
		{"2020-02-02T12:30:00.000Z", "http://cloud.example.net/?OC-Credential=alice&OC-Date=2019-05-14T11%3A01%3A58.135Z&OC-Expires=1200&OC-Verb=GET&OC-Signature=f9e53a1ee23caef10f72ec392c1b537317491b687bfdd224c782be197d9ca2b6", "URL is expired"},
		{"2019-05-14T11:02:00.000Z", "http://cloud.example.net/?OC-Credential=alice&OC-Date=2019-05-14T11%3A01%3A58.135Z&OC-Expires=1200&OC-Verb=GET&OC-Signature=f9e53a1ee23caef10f72ec392c1b537317491b687bfdd224c782be197d9ca2b", "signature mismatch"},
		{"2019-05-14T11:02:00.000Z", "http://cloud.example.net/?OC-Credential=alice&OC-Date=2019-05-14T11%3A01%3A58.135Z&OC-Expires=1200&OC-Verb=GET&OC-Signature=f9e53a1ee23caef10f72ec392c1b537317491b687bfdd224c782be197d9ca2b6", ""},
		{"2019-05-14T11:02:00.000Z", "http://cloud.example.net/?OC-Date=2019-05-14T11%3A01%3A58.135Z&OC-Expires=1200&OC-Verb=GET&OC-Credential=alice&OC-Signature=f9e53a1ee23caef10f72ec392c1b537317491b687bfdd224c782be197d9ca2b6", ""},
		{"2019-05-14T11:02:00.000Z", "http://cloud.example.net/?OC-Algo=PBKDF2%2F10000-SHA512&OC-Date=2019-05-14T11%3A01%3A58.135Z&OC-Expires=1200&OC-Verb=GET&OC-Credential=alice&OC-Signature=f9e53a1ee23caef10f72ec392c1b537317491b687bfdd224c782be197d9ca2b6", ""},
		{"2024-02-07T12:03:11.966Z", "http://localhost:33001/try?id=1&id=2&OC-Credential=user&OC-Date=2024-02-07T12%3A03%3A11.966Z&OC-Expires=2&OC-Verb=GET&OC-Algo=PBKDF2%2F10000-SHA512&OC-Signature=86e21a1efbf0be989a206109cfedf70a22f338dc8995e849ce002032bc6741c5", ""},
	}

	for _, tt := range tests {
		u := userpb.User{
			Id:          &userpb.UserId{OpaqueId: "useri"},
			DisplayName: "Test User",
		}
		ctx := revactx.ContextSetUser(context.Background(), &u)

		pua.Now = func() time.Time {
			t, _ := time.Parse(time.RFC3339, tt.now)
			return t
		}

		r := httptest.NewRequest("", tt.url, nil).WithContext(ctx)
		err := pua.validate(r)
		if tt.errorMessage == "" {
			assert.Nil(t, err)
		} else {
			assert.EqualError(t, err, tt.errorMessage)
		}
	}
}

// signSignedRequest builds a request for the given verb/base URL/params, signed
// with the given key, using the authenticator's own canonicalization
// (buildUrlToSign) so the resulting OC-Signature is valid for signatureIsValid.
// If workaround is true, the request path keeps a literal "$" but the
// signature is computed against the %24-substituted form, simulating a
// reverse proxy that replaces $ with %24 (see signatureMatches).
func signSignedRequest(t *testing.T, pua SignedURLAuthenticator, verb, baseURL string, params url.Values, key []byte, workaround bool) *http.Request {
	t.Helper()

	unsigned := httptest.NewRequest(verb, baseURL+"?"+params.Encode(), nil)
	toSign := pua.buildUrlToSign(unsigned)
	if workaround {
		toSign = strings.Replace(toSign, "$", "%24", 1)
	}
	sig := pua.createSignature(toSign, key)

	signedParams := url.Values{}
	for k, v := range params {
		signedParams[k] = v
	}
	signedParams.Set(_paramOCSignature, sig)

	return httptest.NewRequest(verb, baseURL+"?"+signedParams.Encode(), nil)
}

func signedURLParams(now time.Time) url.Values {
	params := url.Values{}
	params.Set(_paramOCCredential, "useri")
	params.Set(_paramOCDate, now.Format(time.RFC3339))
	params.Set(_paramOCExpires, "1200")
	params.Set(_paramOCVerb, "GET")
	return params
}

func newTestAuthenticator(s store.Store, now time.Time) SignedURLAuthenticator {
	return SignedURLAuthenticator{
		PreSignedURLConfig: config.PreSignedURL{
			AllowedHTTPMethods: []string{"get"},
			Enabled:            true,
		},
		Store:  s,
		Logger: log.NewLogger(),
		Now:    func() time.Time { return now },
	}
}

func userContext(userID string) context.Context {
	u := userpb.User{Id: &userpb.UserId{OpaqueId: userID}, DisplayName: "Test User"}
	return revactx.ContextSetUser(context.Background(), &u)
}

// TestSignedURLAuth_signatureIsValid_mfaKeyMatchSetsBreadcrumb covers a URL
// signed with the MFA-bound key: the signature verifies and the
// request-scoped breadcrumb is set.
func TestSignedURLAuth_signatureIsValid_mfaKeyMatchSetsBreadcrumb(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-02-02T12:30:00.000Z")
	s := store.NewMemoryStore()
	pua := newTestAuthenticator(s, now)

	mfaKey := []byte("mfa-bound-signing-key")
	require.NoError(t, s.Write(&store.Record{Key: mfa.SigningKeyName("useri"), Value: mfaKey}))

	r := signSignedRequest(t, pua, "GET", "http://cloud.example.net/", signedURLParams(now.Add(-time.Minute)), mfaKey, false)
	r = r.WithContext(userContext("useri"))

	err := pua.signatureIsValid(r)
	assert.NoError(t, err)
	assert.True(t, hasMFAKeyBreadcrumb(r.Context()), "a match against the MFA-bound key must set the breadcrumb")
}

// TestSignedURLAuth_signatureIsValid_normalKeyMatchNoBreadcrumb covers a URL
// signed with the normal key: the signature verifies but the breadcrumb is
// not set.
func TestSignedURLAuth_signatureIsValid_normalKeyMatchNoBreadcrumb(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-02-02T12:30:00.000Z")
	s := store.NewMemoryStore()
	pua := newTestAuthenticator(s, now)

	normalKey := []byte("normal-signing-key")
	require.NoError(t, s.Write(&store.Record{Key: "useri", Value: normalKey}))

	r := signSignedRequest(t, pua, "GET", "http://cloud.example.net/", signedURLParams(now.Add(-time.Minute)), normalKey, false)
	r = r.WithContext(userContext("useri"))

	err := pua.signatureIsValid(r)
	assert.NoError(t, err)
	assert.False(t, hasMFAKeyBreadcrumb(r.Context()), "a match against the normal key must not set the breadcrumb")
}

// TestSignedURLAuth_signatureIsValid_noMatch covers a URL signed with
// neither key: authentication fails and the breadcrumb is not set.
func TestSignedURLAuth_signatureIsValid_noMatch(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-02-02T12:30:00.000Z")
	s := store.NewMemoryStore()
	pua := newTestAuthenticator(s, now)

	require.NoError(t, s.Write(&store.Record{Key: mfa.SigningKeyName("useri"), Value: []byte("mfa-bound-signing-key")}))
	require.NoError(t, s.Write(&store.Record{Key: "useri", Value: []byte("normal-signing-key")}))

	r := signSignedRequest(t, pua, "GET", "http://cloud.example.net/", signedURLParams(now.Add(-time.Minute)), []byte("some-other-key"), false)
	r = r.WithContext(userContext("useri"))

	err := pua.signatureIsValid(r)
	assert.Error(t, err)
	assert.False(t, hasMFAKeyBreadcrumb(r.Context()))
}

// TestSignedURLAuth_signatureIsValid_dollarWorkaround covers the existing
// $->%24 reverse-proxy workaround (owncloud/ocis#10180) for both the
// MFA-bound key attempt and the normal key attempt.
func TestSignedURLAuth_signatureIsValid_dollarWorkaround(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-02-02T12:30:00.000Z")

	t.Run("mfa key", func(t *testing.T) {
		s := store.NewMemoryStore()
		pua := newTestAuthenticator(s, now)
		mfaKey := []byte("mfa-bound-signing-key")
		require.NoError(t, s.Write(&store.Record{Key: mfa.SigningKeyName("useri"), Value: mfaKey}))

		r := signSignedRequest(t, pua, "GET", "http://cloud.example.net/a$b", signedURLParams(now.Add(-time.Minute)), mfaKey, true)
		r = r.WithContext(userContext("useri"))

		err := pua.signatureIsValid(r)
		assert.NoError(t, err)
		assert.True(t, hasMFAKeyBreadcrumb(r.Context()))
	})

	t.Run("normal key", func(t *testing.T) {
		s := store.NewMemoryStore()
		pua := newTestAuthenticator(s, now)
		normalKey := []byte("normal-signing-key")
		require.NoError(t, s.Write(&store.Record{Key: "useri", Value: normalKey}))

		r := signSignedRequest(t, pua, "GET", "http://cloud.example.net/a$b", signedURLParams(now.Add(-time.Minute)), normalKey, true)
		r = r.WithContext(userContext("useri"))

		err := pua.signatureIsValid(r)
		assert.NoError(t, err)
		assert.False(t, hasMFAKeyBreadcrumb(r.Context()))
	})
}

// TestSignedURLAuth_signatureIsValid_mismatchDoesNotDiscloseSignature checks
// that a failed verification never puts the server-computed signature into
// the error. That signature is a valid credential for the URL the caller
// chose, and the error is logged.
func TestSignedURLAuth_signatureIsValid_mismatchDoesNotDiscloseSignature(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-02-02T12:30:00.000Z")
	s := store.NewMemoryStore()
	pua := newTestAuthenticator(s, now)

	normalKey := []byte("normal-signing-key")
	require.NoError(t, s.Write(&store.Record{Key: "useri", Value: normalKey}))

	r := signSignedRequest(t, pua, "GET", "http://cloud.example.net/", signedURLParams(now.Add(-time.Minute)), []byte("some-other-key"), false)
	r = r.WithContext(userContext("useri"))
	computed := pua.createSignature(pua.buildUrlToSign(r), normalKey)

	err := pua.signatureIsValid(r)
	require.Error(t, err)
	assert.NotContains(t, err.Error(), computed)
}

// TestSignedURLAuth_Authenticate_failureLogHasNoSignature checks that a
// rejected signed URL is logged without any signature: neither the one the
// server computed nor the one in the request, which may still be valid.
func TestSignedURLAuth_Authenticate_failureLogHasNoSignature(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2020-02-02T12:30:00.000Z")
	normalKey := []byte("normal-signing-key")
	user := &userpb.User{Id: &userpb.UserId{OpaqueId: "useri"}, Username: "useri"}
	path := "/remote.php/dav/files/useri/"

	newAuth := func(t *testing.T, buf *bytes.Buffer) SignedURLAuthenticator {
		s := store.NewMemoryStore()
		require.NoError(t, s.Write(&store.Record{Key: "useri", Value: normalKey}))

		ub := backendmocks.NewUserBackend(t)
		ub.EXPECT().GetUserByClaims(mock.Anything, "username", "useri").Return(user, "", nil)
		ura := userrolesmocks.NewUserRoleAssigner(t)
		ura.EXPECT().ApplyUserRole(mock.Anything, user).Return(user, nil)

		pua := newTestAuthenticator(s, now)
		pua.UserProvider = ub
		pua.UserRoleAssigner = ura
		pua.Logger = log.Logger{Logger: zerolog.New(buf)}
		return pua
	}

	t.Run("wrong signature", func(t *testing.T) {
		var buf bytes.Buffer
		pua := newAuth(t, &buf)
		r := signSignedRequest(t, pua, "GET", "http://cloud.example.net"+path, signedURLParams(now.Add(-time.Minute)), []byte("some-other-key"), false)
		computed := pua.createSignature(pua.buildUrlToSign(r), normalKey)
		supplied := r.URL.Query().Get(_paramOCSignature)

		_, err := pua.Authenticate(r)
		require.ErrorIs(t, err, ErrAuthenticationFailed)
		assert.NotContains(t, buf.String(), computed)
		assert.NotContains(t, buf.String(), supplied)
		assert.Contains(t, buf.String(), path)
	})

	t.Run("valid signature, wrong method", func(t *testing.T) {
		var buf bytes.Buffer
		pua := newAuth(t, &buf)
		signed := signSignedRequest(t, pua, "GET", "http://cloud.example.net"+path, signedURLParams(now.Add(-time.Minute)), normalKey, false)
		r := httptest.NewRequest(http.MethodHead, signed.URL.String(), nil)
		supplied := r.URL.Query().Get(_paramOCSignature)

		_, err := pua.Authenticate(r)
		require.ErrorIs(t, err, ErrAuthenticationFailed)
		assert.NotContains(t, buf.String(), supplied)
		assert.Contains(t, buf.String(), path)
	})
}
