package debug

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/owncloud/ocis/v2/services/proxy/pkg/config/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// secretName matches the names of fields that hold a secret.
var secretName = regexp.MustCompile(`(?i)(password|secret|apikey)$|^creds$`)

const secretMarker = "SECRET-MARKER-"

// setSecrets sets a marker value in every secret string field of v and
// returns the paths of these fields. It allocates nil struct pointers on the way.
func setSecrets(t *testing.T, v reflect.Value, path string, found map[string]string) {
	t.Helper()
	switch v.Kind() {
	case reflect.Ptr:
		if v.Type().Elem().Kind() != reflect.Struct {
			return
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		setSecrets(t, v.Elem(), path, found)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			if !f.IsExported() {
				continue
			}
			fieldPath := path + "." + f.Name
			fv := v.Field(i)
			if fv.Kind() == reflect.String && (f.Tag.Get("mask") == "password" || secretName.MatchString(f.Name)) {
				assert.Equal(t, "password", f.Tag.Get("mask"), "secret field %s has no mask:\"password\" tag", fieldPath)
				fv.SetString(secretMarker + fieldPath)
				found[fieldPath] = fv.String()
				continue
			}
			setSecrets(t, fv, fieldPath, found)
		}
	}
}

func TestConfigDumpMasksSecrets(t *testing.T) {
	cfg := defaults.DefaultConfig()
	defaults.EnsureDefaults(cfg)
	found := map[string]string{}
	setSecrets(t, reflect.ValueOf(cfg), "Config", found)

	// make sure that the test covers the nested secrets
	for _, p := range []string{
		"Config.Events.Creds",
		"Config.Events.AuthPassword",
		"Config.OIDC.UserinfoCache.Creds",
		"Config.PreSignedURL.SigningKeys.Creds",
		"Config.ServiceAccount.ServiceAccountSecret",
	} {
		require.Contains(t, found, p)
	}

	rec := httptest.NewRecorder()
	configDump(cfg)(rec, httptest.NewRequest(http.MethodGet, "/config", nil))
	require.Equal(t, http.StatusOK, rec.Code)

	body := rec.Body.String()
	for p := range found {
		assert.NotContains(t, body, secretMarker+p, "configDump shows the secret %s", p)
	}
	require.False(t, strings.Contains(body, secretMarker), "configDump shows a secret")
	// values that are not secret stay visible
	require.Contains(t, body, cfg.Events.Endpoint)
	require.NotEmpty(t, cfg.Events.Endpoint)
}
