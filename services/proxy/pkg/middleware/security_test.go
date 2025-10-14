package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"gotest.tools/v3/assert"
)

func TestLoadCSPConfig(t *testing.T) {
	// setup test env
	yaml := `
directives:
  frame-src:
    - '''self'''
    - 'https://embed.diagrams.net/'
    - 'https://${ONLYOFFICE_DOMAIN|onlyoffice.owncloud.test}/'
    - 'https://${COLLABORA_DOMAIN|collabora.owncloud.test}/'
`

	config, err := loadCSPConfig([]byte(yaml))
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, config.Directives["frame-src"][0], "'self'")
	assert.Equal(t, config.Directives["frame-src"][1], "https://embed.diagrams.net/")
	assert.Equal(t, config.Directives["frame-src"][2], "https://onlyoffice.owncloud.test/")
	assert.Equal(t, config.Directives["frame-src"][3], "https://collabora.owncloud.test/")
}

func TestStrictTransportSecurity(t *testing.T) {
	// Create test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Use production security middleware config
	cspConfig := &config.CSP{
		Directives: map[string][]string{
			"default-src": {"'none'"},
		},
	}
	securityMiddleware := Security(cspConfig)

	// Test HTTPS request, url not important, only headers will be checked
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Forwarded-Proto", "https")

	rr := httptest.NewRecorder()
	securityMiddleware(handler).ServeHTTP(rr, req)

	hstsHeader := rr.Header().Get("Strict-Transport-Security")

	// HSTS header should contain includeSubDomains
	expected := "max-age=315360000; includeSubDomains; preload"
	assert.Equal(t, hstsHeader, expected, "HSTS header missing includeSubDomains directive - subdomains not protected")
}
