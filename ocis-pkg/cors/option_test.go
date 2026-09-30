package cors

import "testing"

func TestAllowsAnyOrigin(t *testing.T) {
	testCases := []struct {
		name    string
		origins []string
		want    bool
	}{
		{"nil list", nil, true},
		{"empty list", []string{}, true},
		{"bare wildcard", []string{"*"}, true},
		{"padded wildcard", []string{"* "}, true},
		{"wildcard scheme pattern", []string{"https://*"}, true},
		{"wildcard subdomain pattern", []string{"*.example.com"}, true},
		{"wildcard among concrete origins", []string{"https://a.example", "https://*"}, true},
		{"single concrete origin", []string{"https://a.example"}, false},
		{"multiple concrete origins", []string{"https://a.example", "https://b.example"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := AllowsAnyOrigin(tc.origins); got != tc.want {
				t.Errorf("AllowsAnyOrigin(%q) = %v, want %v", tc.origins, got, tc.want)
			}
		})
	}
}
