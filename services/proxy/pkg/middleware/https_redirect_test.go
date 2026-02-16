package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func TestHTTPSRedirect_UsesTrustedHost(t *testing.T) {
	downstream := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mw := HTTPSRedirect("https://trusted.ocis.local")(downstream)

	req := httptest.NewRequest(http.MethodGet, "/foo?bar=1", nil)
	req.Host = "non-trusted.example"
	req.Header.Set("X-Forwarded-Proto", "http")

	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusPermanentRedirect)
	location := rr.Header().Get("Location")
	assert.Equal(t, location, "https://trusted.ocis.local/foo?bar=1")
}
