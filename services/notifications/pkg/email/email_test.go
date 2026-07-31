package email

import (
	"strings"
	"testing"
)

// TestEscapeStringMapDoesNotMutateInput guards the root cause: escapeStringMap
// must return a new map and leave its argument untouched. The render loop reuses
// the same vars map for every recipient and also feeds the unescaped values to
// the plain-text template, so in-place mutation corrupts later renders.
func TestEscapeStringMapDoesNotMutateInput(t *testing.T) {
	in := map[string]string{
		"SpaceName":   "R&D",
		"SpaceSharer": "<Alice>",
		"ShareLink":   "https://example.test/s/x?a=1&b=2",
	}
	want := map[string]string{
		"SpaceName":   "R&D",
		"SpaceSharer": "<Alice>",
		"ShareLink":   "https://example.test/s/x?a=1&b=2",
	}

	out := escapeStringMap(in)

	// Input must be unchanged.
	for k, v := range want {
		if in[k] != v {
			t.Errorf("escapeStringMap mutated input[%q]: got %q, want %q", k, in[k], v)
		}
	}

	// Output must be escaped exactly once.
	if got, want := out["SpaceName"], "R&amp;D"; got != want {
		t.Errorf("out[SpaceName] = %q, want %q", got, want)
	}
	if got, want := out["SpaceSharer"], "&lt;Alice&gt;"; got != want {
		t.Errorf("out[SpaceSharer] = %q, want %q", got, want)
	}

	// Calling again with the same (still-unescaped) input must be stable.
	out2 := escapeStringMap(in)
	for k := range out {
		if out[k] != out2[k] {
			t.Errorf("escapeStringMap not stable for %q: %q vs %q", k, out[k], out2[k])
		}
	}
}

// TestRenderEmailTemplateRepeatableWithSharedVars reproduces the user-facing
// symptom: the notifications render loop calls RenderEmailTemplate once per
// recipient with the SAME vars map. Rendering twice with the same map must yield
// identical output. Before the fix the map was escaped in place on the first
// render, so the second render's text body and subject carried already-escaped
// values and its HTML body was double-escaped.
func TestRenderEmailTemplateRepeatableWithSharedVars(t *testing.T) {
	vars := map[string]string{
		"SpaceSharer":  "A&B",
		"SpaceGrantee": "Bob",
		"SpaceName":    "R&D",
		"ShareLink":    "https://example.test/s/abc?a=1&b=2",
	}

	first, err := RenderEmailTemplate(SharedSpace, "en", "en", "", "", vars)
	if err != nil {
		t.Fatalf("first render failed: %v", err)
	}
	second, err := RenderEmailTemplate(SharedSpace, "en", "en", "", "", vars)
	if err != nil {
		t.Fatalf("second render failed: %v", err)
	}

	if first.Subject != second.Subject {
		t.Errorf("subject differs across renders with a shared vars map:\n first  = %q\n second = %q", first.Subject, second.Subject)
	}
	if first.TextBody != second.TextBody {
		t.Errorf("text body differs across renders with a shared vars map:\n first  = %q\n second = %q", first.TextBody, second.TextBody)
	}
	if first.HTMLBody != second.HTMLBody {
		t.Errorf("html body differs across renders with a shared vars map:\n first  = %q\n second = %q", first.HTMLBody, second.HTMLBody)
	}

	// The plain-text body must contain the raw ampersand, not an HTML entity.
	if !strings.Contains(first.TextBody, "R&D") {
		t.Errorf("text body should contain the unescaped space name %q, got: %q", "R&D", first.TextBody)
	}
	// The HTML body must contain the value escaped exactly once (not double-escaped).
	if !strings.Contains(first.HTMLBody, "R&amp;D") {
		t.Errorf("html body should contain the single-escaped space name %q, got: %q", "R&amp;D", first.HTMLBody)
	}
	if strings.Contains(first.HTMLBody, "R&amp;amp;D") {
		t.Errorf("html body is double-escaped, contains %q: %q", "R&amp;amp;D", first.HTMLBody)
	}
}
