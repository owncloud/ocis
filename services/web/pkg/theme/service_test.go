package theme_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
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

func TestService_GetConcurrent(t *testing.T) {
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
	add("base2/theme.json", map[string]interface{}{
		"base2": "base2",
		"common": map[string]interface{}{
			"shareRoles": map[string]interface{}{
				unifiedrole.UnifiedRoleViewerID: map[string]interface{}{
					"label": "Observer",
				},
			},
		},
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

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			baseTheme := "base"
			roleName := "Viewer"
			if i%2 == 0 {
				baseTheme = "base2"
				roleName = "Observer"
			}
			r.SetPathValue("id", baseTheme)

			w := httptest.NewRecorder()
			service.Get(w, r)

			jsonData := gjson.Parse(w.Body.String())
			// baseTheme
			assert.Equal(t, jsonData.Get(baseTheme).String(), baseTheme)
			// brandingTheme
			assert.Equal(t, jsonData.Get("_branding").String(), "_branding")
			// themeDefaults
			assert.Equal(t, jsonData.Get("common.shareRoles."+unifiedrole.UnifiedRoleViewerID+".label").String(), roleName)
		}()
	}
	wg.Wait()

	// check one last run
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
