package svc

import (
	"html"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/owncloud/ocis/v2/services/idp/pkg/config"
	"go.opentelemetry.io/otel/trace"
)

func newTestService(t *testing.T, cfgFn func(*config.Config)) (Service, *httptest.Server) {
	t.Helper()
	tmpDir, _ := os.MkdirTemp("", "idp-test-*")
	t.Cleanup(func() { os.RemoveAll(tmpDir) })

	encryptionKey := filepath.Join(tmpDir, "encryption.key")
	os.WriteFile(encryptionKey, []byte("test-encryption-secret-32-bytes!"), 0600)

	cfg := &config.Config{
		IDP: config.Settings{
			Iss:                        "https://localhost:9200",
			IdentityManager:            "guest",
			Insecure:                   true,
			AllowClientGuests:          true,
			EncryptionSecretFile:       encryptionKey,
			SigningMethod:              "PS256",
			IdentifierRegistrationConf: filepath.Join(tmpDir, "clients.yaml"),
		},
		Clients: []config.Client{{
			ID: "test", Trusted: true,
			RedirectURIs: []string{"https://localhost:9200/callback"},
		}},
		Commons: &shared.Commons{OcisURL: "https://localhost:9200"},
	}

	if cfgFn != nil {
		cfgFn(cfg)
	}

	svc := NewService(
		Config(cfg),
		Logger(log.NewLogger(log.Level("error"))),
		TraceProvider(trace.NewNoopTracerProvider()),
	)

	srv := httptest.NewServer(svc)
	t.Cleanup(srv.Close)
	return svc, srv
}

func getPage(t *testing.T, srv *httptest.Server, path string, acceptLang string) (int, string, http.Header) {
	t.Helper()
	req, err := http.NewRequest("GET", srv.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if acceptLang != "" {
		req.Header.Set("Accept-Language", acceptLang)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body), resp.Header
}

func TestService(t *testing.T) {
	_, srv := newTestService(t, nil)

	tests := []struct {
		name     string
		path     string
		wantCode int
		contains string
	}{
		{"identifier_route", "/signin/v1/identifier", 200, ""},
		{"identifier_index", "/signin/v1/identifier/index.html", 200, ""},
		{"welcome_route", "/signin/v1/welcome", 200, ""},
		{"goodbye_route", "/signin/v1/goodbye", 200, ""},
		{"oidc_discovery", "/.well-known/openid-configuration", 200, "issuer"},
		{"jwks_endpoint", "/konnect/v1/jwks.json", 200, "keys"},
		{"handles_204", "/signin/v1/identifier", 200, "res.status === 204"},
		{"logon_mode", "/signin/v1/identifier", 200, `params: [username, password, "1"]`},
		{"query_params", "/signin/v1/identifier", 200, "window.location.search"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, body, _ := getPage(t, srv, tt.path, "")
			if code != tt.wantCode {
				t.Errorf("status = %d, want %d", code, tt.wantCode)
			}
			if tt.contains != "" && !strings.Contains(body, tt.contains) {
				t.Errorf("missing %q", tt.contains)
			}
		})
	}
}

func TestTemplateRendering(t *testing.T) {
	_, srv := newTestService(t, nil)

	rawPlaceholder := regexp.MustCompile(`__[A-Z_]+__`)

	pages := []struct {
		name        string
		path        string
		mustContain []string
	}{
		{"index", "/signin/v1/identifier", []string{
			"Username", "Password", "Sign in", `<html lang="en">`,
		}},
		{"welcome", "/signin/v1/welcome", []string{
			"You are signed in", `<html lang="en">`,
		}},
		{"goodbye", "/signin/v1/goodbye", []string{
			"You are signed out", `<html lang="en">`,
		}},
	}

	for _, tt := range pages {
		t.Run(tt.name, func(t *testing.T) {
			code, body, headers := getPage(t, srv, tt.path, "")
			if code != 200 {
				t.Fatalf("status = %d, want 200", code)
			}

			ct := headers.Get("Content-Type")
			if ct != "text/html; charset=utf-8" {
				t.Errorf("Content-Type = %q, want text/html; charset=utf-8", ct)
			}

			for _, s := range tt.mustContain {
				if !strings.Contains(body, s) {
					t.Errorf("missing %q in body", s)
				}
			}

			if matches := rawPlaceholder.FindAllString(body, -1); len(matches) > 0 {
				t.Errorf("raw placeholders remain in output: %v", matches)
			}
		})
	}
}

func TestLocalization(t *testing.T) {
	_, srv := newTestService(t, nil)

	tests := []struct {
		lang         string
		path         string
		wantContains []string
		wantLang     string
	}{
		{"de", "/signin/v1/identifier", []string{"Benutzername", "Passwort"}, "de"},
		{"fr", "/signin/v1/identifier", []string{"Utilisateur", "Mot de passe"}, "fr"},
		{"nl", "/signin/v1/identifier", []string{"Gebruikersnaam", "Wachtwoord"}, "nl"},
		{"de", "/signin/v1/welcome", []string{"Sie sind angemeldet"}, "de"},
		{"de", "/signin/v1/goodbye", []string{"dieses Fenster"}, "de"},
		{"ja", "/signin/v1/identifier", []string{"Username", "Password"}, "en"},
		{"", "/signin/v1/identifier", []string{"Username", "Password"}, "en"},
	}

	for _, tt := range tests {
		name := tt.lang + "_" + tt.path
		if tt.lang == "" {
			name = "empty_" + tt.path
		}
		t.Run(name, func(t *testing.T) {
			_, body, _ := getPage(t, srv, tt.path, tt.lang)
			for _, s := range tt.wantContains {
				if !strings.Contains(body, s) {
					t.Errorf("lang=%q path=%q: missing %q", tt.lang, tt.path, s)
				}
			}
			langAttr := `<html lang="` + tt.wantLang + `">`
			if !strings.Contains(body, langAttr) {
				t.Errorf("lang=%q: missing %q", tt.lang, langAttr)
			}
		})
	}
}

func TestHTMLStructure(t *testing.T) {
	_, srv := newTestService(t, nil)

	t.Run("index_form_elements", func(t *testing.T) {
		_, body, _ := getPage(t, srv, "/signin/v1/identifier", "")
		required := []string{
			"<form", `id="login-form"`, `id="oc-login-username"`, `id="oc-login-password"`, `id="submit-btn"`,
			`fetchApi("/_/logon"`,
		}
		for _, s := range required {
			if !strings.Contains(body, s) {
				t.Errorf("missing %q", s)
			}
		}
	})

	t.Run("welcome_structure", func(t *testing.T) {
		_, body, _ := getPage(t, srv, "/signin/v1/welcome", "")
		if !strings.Contains(body, "oc-message") {
			t.Error("missing oc-message class")
		}
		if !strings.Contains(body, "oc-footer-message") {
			t.Error("missing oc-footer-message class")
		}
	})

	t.Run("goodbye_structure", func(t *testing.T) {
		_, body, _ := getPage(t, srv, "/signin/v1/goodbye", "")
		if !strings.Contains(body, "oc-message") {
			t.Error("missing oc-message class")
		}
		if !strings.Contains(body, "oc-footer-message") {
			t.Error("missing oc-footer-message class")
		}
	})
}

func TestPasswordResetLink(t *testing.T) {
	t.Run("with_reset_uri", func(t *testing.T) {
		_, srv := newTestService(t, func(cfg *config.Config) {
			cfg.Service.PasswordResetURI = "https://example.com/reset"
		})
		_, body, _ := getPage(t, srv, "/signin/v1/identifier", "")
		if !strings.Contains(body, `<a href="https://example.com/reset"`) {
			t.Error("missing password reset link")
		}
		if !strings.Contains(body, "Reset password") {
			t.Error("missing reset password label")
		}
	})

	t.Run("without_reset_uri", func(t *testing.T) {
		_, srv := newTestService(t, nil)
		_, body, _ := getPage(t, srv, "/signin/v1/identifier", "")
		if strings.Contains(body, `<a href=`) {
			t.Error("unexpected link found when no reset URI configured")
		}
	})

	t.Run("reset_link_localized", func(t *testing.T) {
		_, srv := newTestService(t, func(cfg *config.Config) {
			cfg.Service.PasswordResetURI = "https://example.com/reset"
		})
		_, body, _ := getPage(t, srv, "/signin/v1/identifier", "de")
		if !strings.Contains(body, "Passwort zurücksetzen") {
			t.Error("reset label not translated to German")
		}
	})
}

func TestCSPNonceConsistency(t *testing.T) {
	_, srv := newTestService(t, nil)

	_, body, _ := getPage(t, srv, "/signin/v1/identifier", "")

	metaNonce := regexp.MustCompile(`<meta property="csp-nonce" content="([^"]+)"`)
	scriptNonce := regexp.MustCompile(`<script nonce="([^"]+)"`)

	metaMatch := metaNonce.FindStringSubmatch(body)
	scriptMatch := scriptNonce.FindStringSubmatch(body)

	if len(metaMatch) < 2 {
		t.Fatal("could not find CSP nonce in <meta> tag")
	}
	if len(scriptMatch) < 2 {
		t.Fatal("could not find nonce in <script> tag")
	}
	if metaMatch[1] != scriptMatch[1] {
		t.Errorf("nonce mismatch: meta=%q script=%q", metaMatch[1], scriptMatch[1])
	}
	if len(metaMatch[1]) < 16 {
		t.Errorf("nonce too short: %q", metaMatch[1])
	}
}

func TestNoInlineJSTranslations(t *testing.T) {
	_, srv := newTestService(t, nil)

	_, body, _ := getPage(t, srv, "/signin/v1/identifier", "")

	// L10n strings must live in data-* attributes, never in raw JS string literals.
	// Check that none of the MSG_* values appear as JS var assignments.
	jsVarPattern := regexp.MustCompile(`var MSG_\w+\s*=\s*"[^"]*[a-zA-Z]`)
	if matches := jsVarPattern.FindAllString(body, -1); len(matches) > 0 {
		t.Errorf("l10n strings found in inline JS (should use data-* attributes): %v", matches)
	}

	// Verify data-* attributes are present on the form
	dataAttrs := []string{
		`data-msg-required="`,
		`data-msg-invalid="`,
		`data-msg-failed="`,
		`data-msg-default="`,
		`data-msg-signing-in="`,
		`data-msg-sign-in="`,
	}
	for _, attr := range dataAttrs {
		if !strings.Contains(body, attr) {
			t.Errorf("missing %s on form element", attr)
		}
	}

	// Verify JS reads from dataset
	if !strings.Contains(body, "form.dataset.") {
		t.Error("JS should read l10n strings from form.dataset")
	}
}

func TestQuotesInTranslationsEscaped(t *testing.T) {
	// Quotes in translations (e.g. French «Nom d'utilisateur») can break
	// data-* attributes if not escaped. Verify rendered output is safe.
	tpl := []byte(`<form data-msg="__MSG__"></form>`)
	req := httptest.NewRequest("GET", "/signin/v1/identifier", nil)

	idp := &IDP{
		logger:     log.NewLogger(log.Level("error")),
		config:     &config.Config{},
		translator: l10n.NewTranslator("en", "idp", nil),
	}

	replacements := map[string]string{
		"__MSG__": html.EscapeString(`Nom d'utilisateur "obligatoire"`),
	}
	rendered, _ := idp.renderTemplate(tpl, "", req, replacements)
	out := string(rendered)

	if strings.Contains(out, `"obligatoire"`) {
		t.Error("unescaped double quote in data attribute")
	}
	if !strings.Contains(out, `&#34;obligatoire&#34;`) {
		t.Error("double quotes not escaped to &#34;")
	}
}

func TestStaticAssets(t *testing.T) {
	_, srv := newTestService(t, nil)

	code, body, headers := getPage(t, srv, "/signin/v1/static/theme.css", "")
	if code != 200 {
		t.Fatalf("theme.css status = %d, want 200", code)
	}
	ct := headers.Get("Content-Type")
	if !strings.Contains(ct, "css") {
		t.Errorf("Content-Type = %q, want css", ct)
	}
	if len(body) == 0 {
		t.Error("theme.css is empty")
	}
}

func TestSecurityHeaders(t *testing.T) {
	_, srv := newTestService(t, nil)

	pages := []string{
		"/signin/v1/identifier",
		"/signin/v1/welcome",
		"/signin/v1/goodbye",
	}

	nonceRe := regexp.MustCompile(`nonce="([^"]+)"`)

	for _, p := range pages {
		t.Run(p, func(t *testing.T) {
			_, body, headers := getPage(t, srv, p, "")
			nonceMatch := nonceRe.FindStringSubmatch(body)
			if len(nonceMatch) < 2 {
				t.Fatal("could not find nonce in HTML")
			}
			nonce := nonceMatch[1]

			for name, want := range map[string]string{
				"X-Frame-Options":        "DENY",
				"X-Content-Type-Options": "nosniff",
				"Referrer-Policy":        "origin",
			} {
				if got := headers.Get(name); got != want {
					t.Errorf("%s = %q, want %q", name, got, want)
				}
			}

			csp := headers.Get("Content-Security-Policy")
			if csp == "" {
				t.Fatal("missing Content-Security-Policy header")
			}
			for _, substr := range []string{
				"script-src 'nonce-" + nonce + "'",
				"frame-ancestors 'none'",
				"base-uri 'none'",
			} {
				if !strings.Contains(csp, substr) {
					t.Errorf("CSP missing %q in %q", substr, csp)
				}
			}
		})
	}
}

func TestNoRawPlaceholders(t *testing.T) {
	_, srv := newTestService(t, nil)
	rawPlaceholder := regexp.MustCompile(`__[A-Z_]+__`)

	paths := []string{
		"/signin/v1/identifier",
		"/signin/v1/welcome",
		"/signin/v1/goodbye",
	}

	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			_, body, _ := getPage(t, srv, p, "")
			if matches := rawPlaceholder.FindAllString(body, -1); len(matches) > 0 {
				t.Errorf("raw placeholders in %s: %v", p, matches)
			}
		})
	}
}
