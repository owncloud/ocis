package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/router"
)

type stubAuth struct {
	name  string
	err   error // non-nil => this authenticator fails with this error
	calls *[]string
}

func (s stubAuth) Authenticate(r *http.Request) (*http.Request, error) {
	*s.calls = append(*s.calls, s.name)
	if s.err != nil {
		return nil, s.err
	}
	return r, nil
}

// A transient OIDC failure followed by a succeeding authenticator must serve 200,
// not 503 — the deferral flag exists so authenticator order does not matter.
func TestTransientThenSuccessServes200(t *testing.T) {
	var calls []string
	auths := []Authenticator{
		stubAuth{name: "oidc-transient", err: oidc.ErrTemporarilyUnavailable, calls: &calls},
		stubAuth{name: "public-share-ok", err: nil, calls: &calls},
	}

	served := false
	handler := Authentication(auths, EnableBasicAuth(false))(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { served = true }),
	)

	req := httptest.NewRequest(http.MethodGet, "http://example.com/dav/public-files/x", http.NoBody)
	req = req.WithContext(router.SetRoutingInfo(req.Context(), router.RoutingInfo{}))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	t.Logf("authenticators called: %v, status: %d", calls, rr.Code)
	if !served || rr.Code != http.StatusOK {
		t.Fatalf("want 200 served by later authenticator, got status=%d served=%v (calls=%v)", rr.Code, served, calls)
	}
}
