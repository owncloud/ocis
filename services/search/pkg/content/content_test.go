package content_test

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/search/pkg/content"
)

func TestCleanContent(t *testing.T) {
	tests := []struct {
		given  string
		expect string
	}{
		{
			given:  "find can keeper should keeper will",
			expect: "keeper keeper",
		},
		{
			given:  "user1 shares the file to Marie",
			expect: "user1 shares file marie",
		},
		{
			given:  "content contains https://localhost/remote.php/dav/files/admin/Photos/San%20Francisco.jpg and stop word",
			expect: "content contains https://localhost/remote.php/dav/files/admin/photos/san%20francisco.jpg stop word",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.given, func(t *testing.T) {
			Equal(t, tc.expect, content.CleanString(tc.given, "en"))
		})
	}
}
