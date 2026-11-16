package svc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/mfa"
	"github.com/owncloud/reva/v2/pkg/autoprop"
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

// serveSigningKey drives GetSigningKey the way production does, deliberately
// NOT setting the MFA status on the request context directly:
//   - mfaOnWire carries the proxy's auto-propagated MFA status on the request
//     HEADERS, exactly as reva's outbound RoundTripper serializes it;
//   - withAutopropHandler wraps the handler in autoprop.NewHttpHandler(), which
//     is what the real OCS router installs (see service.go) to reconstruct that
//     metadata into the request context. Only then does revactx.HasMFA see it.
//
// The user is injected as the ExtractAccountUUID middleware would, after the
// autoprop handler has run.
func serveSigningKey(t *testing.T, o Ocs, user *cs3user.User, mfaOnWire, withAutopropHandler bool) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/ocs/v2.php/cloud/user/signing-key", nil)
	req.Header.Set("Accept", "application/json")
	if mfaOnWire {
		autoprop.AppendToHttpRequest(req, revactx.SetMFA(context.Background()))
	}

	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o.GetSigningKey(w, r.WithContext(revactx.ContextSetUser(r.Context(), user)))
	})
	if withAutopropHandler {
		h = autoprop.NewHttpHandler()(h)
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func decodeSigningKey(t *testing.T, w *httptest.ResponseRecorder) signingKeyResponse {
	t.Helper()
	var body signingKeyResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	return body
}

// TestGetSigningKey_mfaRequestIssuesMFABoundKeyOnly asserts that a request
// carrying MFA on the wire (reconstructed by the autoprop handler the router
// installs) is issued the MFA-bound key, and never the normal key.
func TestGetSigningKey_mfaRequestIssuesMFABoundKeyOnly(t *testing.T) {
	s := store.NewMemoryStore()
	o := newTestOcs(s)
	u := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "useri"}}

	w := serveSigningKey(t, o, u, true /* mfaOnWire */, true /* autoprop handler */)
	require.Equal(t, http.StatusOK, w.Code)
	body := decodeSigningKey(t, w)
	require.NotEmpty(t, body.OCS.Data.SigningKey)

	mfaRecord, err := s.Read(mfa.SigningKeyName("useri"))
	require.NoError(t, err)
	require.Len(t, mfaRecord, 1)
	assert.Equal(t, string(mfaRecord[0].Value), body.OCS.Data.SigningKey)

	_, err = s.Read("useri")
	assert.ErrorIs(t, err, store.ErrNotFound, "an MFA request must never create/return the normal key")
}

// TestGetSigningKey_nonMFARequestIssuesNormalKeyOnly asserts that a request
// without MFA keeps getting the normal key, and never the MFA-bound key.
func TestGetSigningKey_nonMFARequestIssuesNormalKeyOnly(t *testing.T) {
	s := store.NewMemoryStore()
	o := newTestOcs(s)
	u := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "useri"}}

	w := serveSigningKey(t, o, u, false /* mfaOnWire */, true /* autoprop handler */)
	require.Equal(t, http.StatusOK, w.Code)
	body := decodeSigningKey(t, w)
	require.NotEmpty(t, body.OCS.Data.SigningKey)

	normalRecord, err := s.Read("useri")
	require.NoError(t, err)
	require.Len(t, normalRecord, 1)
	assert.Equal(t, string(normalRecord[0].Value), body.OCS.Data.SigningKey)

	_, err = s.Read(mfa.SigningKeyName("useri"))
	assert.ErrorIs(t, err, store.ErrNotFound, "a non-MFA request must never create/return the MFA-bound key")
}

// TestGetSigningKey_mfaOnWireNeedsAutopropHandler is the regression guard for
// the real defect: if the router does not reconstruct the auto-propagated MFA
// status (i.e. autoprop.NewHttpHandler is not installed), then even a request
// that carries MFA on the wire falls back to the normal key — which would
// silently break every vault signed URL. With the handler absent the MFA-bound
// key must NOT be issued; with it present (covered above) it must.
func TestGetSigningKey_mfaOnWireNeedsAutopropHandler(t *testing.T) {
	s := store.NewMemoryStore()
	o := newTestOcs(s)
	u := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "useri"}}

	w := serveSigningKey(t, o, u, true /* mfaOnWire */, false /* NO autoprop handler */)
	require.Equal(t, http.StatusOK, w.Code)

	// without the handler the MFA status never reaches the context, so the
	// normal key is issued and the MFA-bound key is never created.
	_, err := s.Read("useri")
	require.NoError(t, err)
	_, err = s.Read(mfa.SigningKeyName("useri"))
	assert.ErrorIs(t, err, store.ErrNotFound)
}

// TestGetSigningKey_mfaAndNormalKeysDiffer proves the two keys are distinct:
// an MFA request and a non-MFA request for the same user must never be handed
// the same key value.
func TestGetSigningKey_mfaAndNormalKeysDiffer(t *testing.T) {
	s := store.NewMemoryStore()
	o := newTestOcs(s)
	u := &cs3user.User{Id: &cs3user.UserId{OpaqueId: "useri"}}

	mfaBody := decodeSigningKey(t, serveSigningKey(t, o, u, true, true))
	normalBody := decodeSigningKey(t, serveSigningKey(t, o, u, false, true))

	assert.NotEqual(t, mfaBody.OCS.Data.SigningKey, normalBody.OCS.Data.SigningKey)
}

// TestGetSigningKey_guestRefused asserts the existing guest refusal is
// unaffected by the MFA-bound key change, with or without MFA.
func TestGetSigningKey_guestRefused(t *testing.T) {
	for _, mfaOnWire := range []bool{false, true} {
		s := store.NewMemoryStore()
		o := newTestOcs(s)
		u := &cs3user.User{
			Id: &cs3user.UserId{OpaqueId: "guesti", Type: cs3user.UserType_USER_TYPE_GUEST},
		}
		w := serveSigningKey(t, o, u, mfaOnWire, true)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	}
}

// TestGetSigningKey_publicShareRefused asserts the existing public-share guest
// refusal (OCISDEV-720) is unaffected by the MFA-bound key change.
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
	w := serveSigningKey(t, o, u, false, true)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
