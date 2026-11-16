package svc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cs3types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"

	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-micro.dev/v4/store"
)

// signingKeyResponse mirrors response.Response/data.SigningKey enough to
// decode the JSON body produced by GetSigningKey in tests.
type signingKeyResponse struct {
	OCS struct {
		Meta struct {
			StatusCode int `json:"statuscode"`
		} `json:"meta"`
		Data struct {
			User       string `json:"user"`
			SigningKey string `json:"signing-key"`
		} `json:"data"`
	} `json:"ocs"`
}

func newTestOcs(s store.Store) Ocs {
	return Ocs{
		logger: log.NewLogger(),
		store:  s,
	}
}

// TestGetSigningKey_mfaContextIssuesMFABoundKeyOnly asserts that a request
// whose context carries MFA is issued the MFA-bound key, stored under the
// MFA-bound key name, and never the normal key.
func TestGetSigningKey_mfaContextIssuesMFABoundKeyOnly(t *testing.T) {
	s := store.NewMemoryStore()
	o := newTestOcs(s)

	u := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "useri"}}
	ctx := revactx.ContextSetUser(t.Context(), u)
	ctx = revactx.SetMFA(ctx)

	r := httptest.NewRequest(http.MethodGet, "/ocs/v2.php/cloud/user/signing-key", nil).WithContext(ctx)
	r.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	o.GetSigningKey(w, r)

	require.Equal(t, http.StatusOK, w.Code)

	var body signingKeyResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.NotEmpty(t, body.OCS.Data.SigningKey)

	// the returned key must be the one stored under the MFA-bound key name ...
	mfaRecord, err := s.Read(mfaSigningKeyName("useri"))
	require.NoError(t, err)
	require.Len(t, mfaRecord, 1)
	assert.Equal(t, string(mfaRecord[0].Value), body.OCS.Data.SigningKey)

	// ... and never under the normal key name.
	_, err = s.Read("useri")
	assert.ErrorIs(t, err, store.ErrNotFound, "an MFA request must never create/return the normal key")
}

// TestGetSigningKey_nonMFAContextIssuesNormalKeyOnly asserts that a request
// whose context does not carry MFA keeps getting the normal key, stored
// under the normal key name, and never the MFA-bound key.
func TestGetSigningKey_nonMFAContextIssuesNormalKeyOnly(t *testing.T) {
	s := store.NewMemoryStore()
	o := newTestOcs(s)

	u := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "useri"}}
	ctx := revactx.ContextSetUser(t.Context(), u)
	// no revactx.SetMFA: this request does not carry MFA.

	r := httptest.NewRequest(http.MethodGet, "/ocs/v2.php/cloud/user/signing-key", nil).WithContext(ctx)
	r.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	o.GetSigningKey(w, r)

	require.Equal(t, http.StatusOK, w.Code)

	var body signingKeyResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.NotEmpty(t, body.OCS.Data.SigningKey)

	// the returned key must be the one stored under the normal key name ...
	normalRecord, err := s.Read("useri")
	require.NoError(t, err)
	require.Len(t, normalRecord, 1)
	assert.Equal(t, string(normalRecord[0].Value), body.OCS.Data.SigningKey)

	// ... and never under the MFA-bound key name.
	_, err = s.Read(mfaSigningKeyName("useri"))
	assert.ErrorIs(t, err, store.ErrNotFound, "a non-MFA request must never create/return the MFA-bound key")
}

// TestGetSigningKey_mfaAndNormalKeysDiffer proves the two keys are
// distinct: an MFA request and a non-MFA request for the same user must
// never be handed the same key value.
func TestGetSigningKey_mfaAndNormalKeysDiffer(t *testing.T) {
	s := store.NewMemoryStore()
	o := newTestOcs(s)
	u := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "useri"}}

	mfaCtx := revactx.SetMFA(revactx.ContextSetUser(t.Context(), u))
	r1 := httptest.NewRequest(http.MethodGet, "/ocs/v2.php/cloud/user/signing-key", nil).WithContext(mfaCtx)
	r1.Header.Set("Accept", "application/json")
	w1 := httptest.NewRecorder()
	o.GetSigningKey(w1, r1)
	var body1 signingKeyResponse
	require.NoError(t, json.Unmarshal(w1.Body.Bytes(), &body1))

	normalCtx := revactx.ContextSetUser(t.Context(), u)
	r2 := httptest.NewRequest(http.MethodGet, "/ocs/v2.php/cloud/user/signing-key", nil).WithContext(normalCtx)
	r2.Header.Set("Accept", "application/json")
	w2 := httptest.NewRecorder()
	o.GetSigningKey(w2, r2)
	var body2 signingKeyResponse
	require.NoError(t, json.Unmarshal(w2.Body.Bytes(), &body2))

	assert.NotEqual(t, body1.OCS.Data.SigningKey, body2.OCS.Data.SigningKey)
}

// TestGetSigningKey_guestRefused asserts the existing guest refusal is
// unaffected by the MFA-bound key change, with or without MFA.
func TestGetSigningKey_guestRefused(t *testing.T) {
	for _, mfa := range []bool{false, true} {
		s := store.NewMemoryStore()
		o := newTestOcs(s)
		u := &cs3user.User{
			Id: &cs3user.UserId{OpaqueId: "guesti", Type: cs3user.UserType_USER_TYPE_GUEST},
		}
		ctx := revactx.ContextSetUser(t.Context(), u)
		if mfa {
			ctx = revactx.SetMFA(ctx)
		}

		r := httptest.NewRequest(http.MethodGet, "/ocs/v2.php/cloud/user/signing-key", nil).WithContext(ctx)
		r.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()
		o.GetSigningKey(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	}
}

// TestGetSigningKey_publicShareRefused asserts the existing public-share
// guest refusal (OCISDEV-720) is unaffected by the MFA-bound key change.
func TestGetSigningKey_publicShareRefused(t *testing.T) {
	s := store.NewMemoryStore()
	o := newTestOcs(s)
	u := &cs3user.User{
		Id: &cs3user.UserId{OpaqueId: "shareri"},
		Opaque: &cs3types.Opaque{
			Map: map[string]*cs3types.OpaqueEntry{
				"public-share-role": {
					Decoder: "plain",
					Value:   []byte("viewer"),
				},
			},
		},
	}
	ctx := revactx.ContextSetUser(t.Context(), u)

	r := httptest.NewRequest(http.MethodGet, "/ocs/v2.php/cloud/user/signing-key", nil).WithContext(ctx)
	r.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()
	o.GetSigningKey(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
