package svc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

// buildRateLimitedRouter mirrors the middleware composition applied to the
// exportPersonalData route in service.go, so the rate-limiting behavior can be
// exercised in isolation from the full Graph service.
func buildRateLimitedRouter() http.Handler {
	r := chi.NewRouter()
	r.Route("/users/{userID}", func(r chi.Router) {
		r.With(httprate.LimitBy(exportPersonalDataLimit, exportPersonalDataWindow, httprate.KeyByEndpoint)).
			Post("/exportPersonalData", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusAccepted)
			})
	})
	return r
}

func doExportRequest(t *testing.T, h http.Handler, userID string) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/users/"+userID+"/exportPersonalData", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	return rr.Code
}

func TestExportPersonalDataRateLimit(t *testing.T) {
	t.Run("returns 429 after the limit is exceeded", func(t *testing.T) {
		h := buildRateLimitedRouter()

		for i := range exportPersonalDataLimit {
			if code := doExportRequest(t, h, "alice"); code != http.StatusAccepted {
				t.Fatalf("request %d: got status %d, want %d", i+1, code, http.StatusAccepted)
			}
		}

		if code := doExportRequest(t, h, "alice"); code != http.StatusTooManyRequests {
			t.Fatalf("request past limit: got status %d, want %d", code, http.StatusTooManyRequests)
		}
	})

	t.Run("keys per endpoint so one user cannot starve another", func(t *testing.T) {
		h := buildRateLimitedRouter()

		// Exhaust the first user's bucket.
		for range exportPersonalDataLimit {
			doExportRequest(t, h, "alice")
		}
		if code := doExportRequest(t, h, "alice"); code != http.StatusTooManyRequests {
			t.Fatalf("first user should be limited, got %d", code)
		}

		// A different userID resolves to a different endpoint path and bucket.
		if code := doExportRequest(t, h, "bob"); code != http.StatusAccepted {
			t.Fatalf("second user should not be limited, got %d", code)
		}
	})
}
