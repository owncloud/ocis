package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/stretchr/testify/assert"
)

// TestSelectorCookie tests the core functionality of the selector cookie middleware.
func TestSelectorCookie(t *testing.T) {
	const expectedCookieName = "test-selector"
	const expectedCookieValue = "test-value"
	tests := []struct {
		name           string
		hasOIDCContext bool
		config         config.PolicySelector
		checkCookie    func(*testing.T, []*http.Cookie)
	}{
		{
			name:           "successful cookie set with claims selector",
			hasOIDCContext: true,
			config: config.PolicySelector{
				Claims: &config.ClaimsSelectorConf{
					SelectorCookieName: expectedCookieName,
					DefaultPolicy:      expectedCookieValue,
				},
			},
			checkCookie: func(t *testing.T, cookies []*http.Cookie) {
				assert.Len(t, cookies, 1)
				cookie := cookies[0]
				assert.Equal(t, expectedCookieName, cookie.Name)
				assert.Equal(t, expectedCookieValue, cookie.Value)
				assert.Equal(t, "/", cookie.Path)
				assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite, "SameSite=Strict prevents CSRF attacks by blocking cross-site cookie manipulation (OWASP: https://owasp.org/www-project-cheat-sheets/cheatsheets/Session_Management_Cheat_Sheet.html#samesite-attribute)")
				assert.True(t, cookie.Secure, "Secure flag prevents MITM attacks by ensuring cookie only sent over HTTPS (OWASP: https://owasp.org/www-project-cheat-sheets/cheatsheets/Session_Management_Cheat_Sheet.html#secure-attribute)")
				assert.True(t, cookie.HttpOnly, "HttpOnly prevents XSS attacks by blocking JavaScript access to cookie (OWASP: https://owasp.org/www-project-cheat-sheets/cheatsheets/Session_Management_Cheat_Sheet.html#httponly-attribute)")
			},
		},
		{
			name:           "no cookie set without OIDC context",
			hasOIDCContext: false,
			config: config.PolicySelector{
				Claims: &config.ClaimsSelectorConf{
					SelectorCookieName: expectedCookieName,
					DefaultPolicy:      expectedCookieValue,
				},
			},
			checkCookie: func(t *testing.T, cookies []*http.Cookie) {
				assert.Empty(t, cookies)
			},
		},
		{
			name:           "no cookie set without selector config",
			hasOIDCContext: true,
			config:         config.PolicySelector{},
			checkCookie: func(t *testing.T, cookies []*http.Cookie) {
				assert.Empty(t, cookies)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.NewLogger()
			options := []Option{
				Logger(logger),
				PolicySelectorConfig(tt.config),
			}

			handler := SelectorCookie(options...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// This is a no-op handler since we're testing the middleware's cookie setting behavior,
				// not the actual request handling. The middleware should set cookies before this handler is called.
			}))

			req := httptest.NewRequest("GET", "https://example.com", nil)
			if tt.hasOIDCContext {
				req = req.WithContext(oidc.NewContext(req.Context(), map[string]interface{}{
					oidc.OcisRoutingPolicy: expectedCookieValue,
				}))
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			tt.checkCookie(t, w.Result().Cookies())
		})
	}
}
