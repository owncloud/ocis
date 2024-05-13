package theme_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/owncloud/ocis/v2/services/web/pkg/theme"
)

func TestAllowedLogoFileTypes(t *testing.T) {
	type test struct {
		filename string
		mimetype string
		allowed  bool
	}

	tests := []test{
		{filename: "foo.jpg", mimetype: "image/jpeg", allowed: true},
		{filename: "foo.jpeg", mimetype: "image/jpeg", allowed: true},
		{filename: "foo.png", mimetype: "image/png", allowed: true},
		{filename: "foo.gif", mimetype: "image/gif", allowed: true},
		{filename: "foo.tiff", mimetype: "image/tiff", allowed: false},
	}

	for _, tc := range tests {
		assert.Equal(t, theme.IsFiletypePermitted(theme.AllowedLogoFileTypes, tc.filename, tc.mimetype), tc.allowed)
	}
}
