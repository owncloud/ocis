package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/cors"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// TestCorsWildcardOriginDropsCredentials is a regression test for OCISDEV-55
// (F-1.3 Non-restrictive CORS policy): a wildcard/allow-all origin must never be
// combined with Access-Control-Allow-Credentials: true.
func TestCorsWildcardOriginDropsCredentials(t *testing.T) {
	testCases := []struct {
		name             string
		allowedOrigins   []string
		allowCredentials bool
		requestOrigin    string
		wantAllowOrigin  string
		wantCredentials  bool
	}{
		{
			name:             "wildcard origin with credentials drops credentials",
			allowedOrigins:   []string{"*"},
			allowCredentials: true,
			requestOrigin:    "https://sec-consult.com",
			wantAllowOrigin:  "*",
			wantCredentials:  false,
		},
		{
			name:             "empty origins (allow all) with credentials drops credentials",
			allowedOrigins:   []string{},
			allowCredentials: true,
			requestOrigin:    "https://sec-consult.com",
			wantAllowOrigin:  "*",
			wantCredentials:  false,
		},
		{
			name:             "wildcard origin without credentials stays wildcard",
			allowedOrigins:   []string{"*"},
			allowCredentials: false,
			requestOrigin:    "https://sec-consult.com",
			wantAllowOrigin:  "*",
			wantCredentials:  false,
		},
		{
			name:             "explicit whitelisted origin keeps credentials",
			allowedOrigins:   []string{"https://cloud.example.com"},
			allowCredentials: true,
			requestOrigin:    "https://cloud.example.com",
			wantAllowOrigin:  "https://cloud.example.com",
			wantCredentials:  true,
		},
		{
			name:             "non-whitelisted origin is not reflected",
			allowedOrigins:   []string{"https://cloud.example.com"},
			allowCredentials: true,
			requestOrigin:    "https://sec-consult.com",
			wantAllowOrigin:  "",
			wantCredentials:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := Cors(
				cors.Logger(log.NewLogger()),
				cors.AllowedOrigins(tc.allowedOrigins),
				cors.AllowedMethods([]string{"GET", "OPTIONS"}),
				cors.AllowedHeaders([]string{"Authorization"}),
				cors.AllowCredentials(tc.allowCredentials),
			)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			// Exercise the actual (non-preflight) request path.
			req := httptest.NewRequest(http.MethodGet, "/ocs/v1.php/apps/files_sharing/api/v1/shares", nil)
			req.Header.Set("Origin", tc.requestOrigin)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			gotOrigin := rec.Header().Get("Access-Control-Allow-Origin")
			if gotOrigin != tc.wantAllowOrigin {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", gotOrigin, tc.wantAllowOrigin)
			}

			gotCredentials := rec.Header().Get("Access-Control-Allow-Credentials") == "true"
			if gotCredentials != tc.wantCredentials {
				t.Errorf("Access-Control-Allow-Credentials = %v, want %v", gotCredentials, tc.wantCredentials)
			}

			// The insecure combination the pentest flagged must never occur.
			if gotOrigin == "*" && gotCredentials {
				t.Error("insecure CORS: wildcard origin combined with credentials")
			}
		})
	}
}
