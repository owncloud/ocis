package theme_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"

	"github.com/owncloud/ocis/v2/ocis-pkg/x/io/fsx"
	"github.com/owncloud/ocis/v2/services/graph/mocks"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
	"github.com/owncloud/ocis/v2/services/web/pkg/theme"
)

func TestNewService(t *testing.T) {
	t.Run("fails if the options are invalid", func(t *testing.T) {
		_, err := theme.NewService(theme.ServiceOptions{})
		assert.Error(t, err)
	})

	t.Run("success if the options are valid", func(t *testing.T) {
		_, err := theme.NewService(
			theme.ServiceOptions{}.
				WithThemeFS(fsx.NewFallbackFS(fsx.NewMemMapFs(), fsx.NewMemMapFs())).
				WithGatewaySelector(mocks.NewSelectable[gateway.GatewayAPIClient](t)),
		)
		assert.NoError(t, err)
	})
}

func TestService_Get(t *testing.T) {
	primaryFS := fsx.NewMemMapFs()
	fallbackFS := fsx.NewFallbackFS(primaryFS, fsx.NewMemMapFs())

	add := func(filename string, content interface{}) {
		b, err := json.Marshal(content)
		assert.Nil(t, err)

		assert.Nil(t, afero.WriteFile(primaryFS, filename, b, 0644))
	}

	// baseTheme
	add("base/theme.json", map[string]interface{}{
		"base": "base",
	})
	// brandingTheme
	add("_branding/theme.json", map[string]interface{}{
		"_branding": "_branding",
	})

	service, _ := theme.NewService(
		theme.ServiceOptions{}.
			WithThemeFS(fallbackFS).
			WithGatewaySelector(mocks.NewSelectable[gateway.GatewayAPIClient](t)),
	)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.SetPathValue("id", "base")

	w := httptest.NewRecorder()
	service.Get(w, r)

	jsonData := gjson.Parse(w.Body.String())
	// baseTheme
	assert.Equal(t, jsonData.Get("base").String(), "base")
	// brandingTheme
	assert.Equal(t, jsonData.Get("_branding").String(), "_branding")
	// themeDefaults
	assert.Equal(t, jsonData.Get("common.shareRoles."+unifiedrole.UnifiedRoleViewerID+".label").String(), "Viewer")
}

// A theme (including the built-in "owncloud" one, whose dark and high-contrast
// variants each define their own logo) can define its own logo per variant
// under clients.web.themes[].logo. An admin-uploaded logo (branding) must
// override those variant-specific logos too, not just the clients.web.defaults
// ones, otherwise the active theme's own logo silently wins over the admin's
// upload.
func TestService_Get_BrandingLogoOverridesThemeVariants(t *testing.T) {
	primaryFS := fsx.NewMemMapFs()
	fallbackFS := fsx.NewFallbackFS(primaryFS, fsx.NewMemMapFs())

	add := func(filename string, content interface{}) {
		b, err := json.Marshal(content)
		assert.Nil(t, err)

		assert.Nil(t, afero.WriteFile(primaryFS, filename, b, 0644))
	}

	// baseTheme: a theme with its own per-variant logo, as set by the theme author.
	// The first variant has no logo of its own, mirroring the built-in "owncloud"
	// theme's default variant, which falls back to clients.web.defaults.logo.
	add("base/theme.json", map[string]interface{}{
		"clients": map[string]interface{}{
			"web": map[string]interface{}{
				"themes": []interface{}{
					map[string]interface{}{
						"isDark": false,
						"name":   "Light Theme",
					},
					map[string]interface{}{
						"isDark": true,
						"name":   "Dark Theme",
						"logo": map[string]interface{}{
							"topbar":  "themes/base/assets/logo.svg",
							"favicon": "themes/base/assets/favicon.jpg",
							"login":   "themes/base/assets/login.svg",
						},
					},
				},
			},
		},
	})
	// brandingTheme: an admin-uploaded logo
	add("_branding/theme.json", map[string]interface{}{
		"common": map[string]interface{}{
			"logo": "themes/_branding/custom-logo.svg",
		},
		"clients": map[string]interface{}{
			"web": map[string]interface{}{
				"defaults": map[string]interface{}{
					"logo": map[string]interface{}{
						"topbar": "themes/_branding/custom-logo.svg",
						"login":  "themes/_branding/custom-logo.svg",
					},
				},
			},
		},
	})

	service, _ := theme.NewService(
		theme.ServiceOptions{}.
			WithThemeFS(fallbackFS).
			WithGatewaySelector(mocks.NewSelectable[gateway.GatewayAPIClient](t)),
	)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.SetPathValue("id", "base")

	w := httptest.NewRecorder()
	service.Get(w, r)

	jsonData := gjson.Parse(w.Body.String())
	// variant without its own logo: untouched, still absent, no crash
	assert.False(t, jsonData.Get("clients.web.themes.0.logo").Exists())
	// variant with its own logo: overridden by the admin upload
	assert.Equal(t, "themes/_branding/custom-logo.svg", jsonData.Get("clients.web.themes.1.logo.topbar").String())
	assert.Equal(t, "themes/_branding/custom-logo.svg", jsonData.Get("clients.web.themes.1.logo.login").String())
	// favicon is not part of the admin logo upload, so the theme's own favicon must be left untouched
	assert.Equal(t, "themes/base/assets/favicon.jpg", jsonData.Get("clients.web.themes.1.logo.favicon").String())
	// the defaults themselves also carry the admin upload
	assert.Equal(t, "themes/_branding/custom-logo.svg", jsonData.Get("clients.web.defaults.logo.topbar").String())
}

// Without an admin-uploaded logo, a theme's own per-variant logo must be left
// exactly as the theme author defined it.
func TestService_Get_NoBrandingLogoPreservesThemeVariants(t *testing.T) {
	primaryFS := fsx.NewMemMapFs()
	fallbackFS := fsx.NewFallbackFS(primaryFS, fsx.NewMemMapFs())

	add := func(filename string, content interface{}) {
		b, err := json.Marshal(content)
		assert.Nil(t, err)

		assert.Nil(t, afero.WriteFile(primaryFS, filename, b, 0644))
	}

	// baseTheme only, no _branding/theme.json at all (admin never uploaded a logo)
	add("base/theme.json", map[string]interface{}{
		"clients": map[string]interface{}{
			"web": map[string]interface{}{
				"themes": []interface{}{
					map[string]interface{}{
						"isDark": true,
						"name":   "Dark Theme",
						"logo": map[string]interface{}{
							"topbar":  "themes/base/assets/logo.svg",
							"favicon": "themes/base/assets/favicon.jpg",
							"login":   "themes/base/assets/login.svg",
						},
					},
				},
			},
		},
	})

	service, _ := theme.NewService(
		theme.ServiceOptions{}.
			WithThemeFS(fallbackFS).
			WithGatewaySelector(mocks.NewSelectable[gateway.GatewayAPIClient](t)),
	)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.SetPathValue("id", "base")

	w := httptest.NewRecorder()
	service.Get(w, r)

	jsonData := gjson.Parse(w.Body.String())
	assert.Equal(t, "themes/base/assets/logo.svg", jsonData.Get("clients.web.themes.0.logo.topbar").String())
	assert.Equal(t, "themes/base/assets/login.svg", jsonData.Get("clients.web.themes.0.logo.login").String())
	assert.Equal(t, "themes/base/assets/favicon.jpg", jsonData.Get("clients.web.themes.0.logo.favicon").String())
}
