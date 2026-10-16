package http

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireIdentity_RejectsWhenAccountMiddlewareStopsTheRequest(t *testing.T) {
	// stands in for ExtractAccountUUID returning early on an invalid token
	extract := func(next nethttp.Handler) nethttp.Handler {
		return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {})
	}
	called := false
	h := requireIdentity(extract)(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) { called = true }))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(nethttp.MethodPost, "/", nil))

	if rec.Code != nethttp.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
	if called {
		t.Fatal("next handler must not run")
	}
}

func TestRequireIdentity_PassesThroughWhenAccountMiddlewareContinues(t *testing.T) {
	extract := func(next nethttp.Handler) nethttp.Handler { return next }
	h := requireIdentity(extract)(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusNoContent)
	}))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(nethttp.MethodPost, "/", nil))

	if rec.Code != nethttp.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}
