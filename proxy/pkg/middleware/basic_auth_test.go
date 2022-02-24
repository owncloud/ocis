package middleware

import (
	"net/http/httptest"
	"testing"
)

/**/

func TestBasicAuth__isPublicLink(t *testing.T) {
	tests := []struct {
		url      string
		username string
		expected bool
	}{
		{url: "/remote.php/dav/public-files/", username: "", expected: false},
		{url: "/remote.php/dav/public-files/", username: "abc", expected: false},
		{url: "/remote.php/dav/public-files/", username: "private", expected: false},
		{url: "/remote.php/dav/public-files/", username: "public", expected: true},
		{url: "/ocs/v1.php/cloud/capabilities", username: "", expected: false},
		{url: "/ocs/v1.php/cloud/capabilities", username: "abc", expected: false},
		{url: "/ocs/v1.php/cloud/capabilities", username: "private", expected: false},
		{url: "/ocs/v1.php/cloud/capabilities", username: "public", expected: true},
		{url: "/ocs/v1.php/cloud/users/admin", username: "public", expected: false},
	}
	ba := basicAuth{}

	for _, tt := range tests {
		req := httptest.NewRequest("", tt.url, nil)

		if tt.username != "" {
			req.SetBasicAuth(tt.username, "")
		}

		result := ba.isPublicLink(req)
		if result != tt.expected {
			t.Errorf("with %s expected %t got %t", tt.url, tt.expected, result)
		}
	}
}
