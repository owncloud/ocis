package keycloak

import (
	"testing"

	kcpkg "github.com/owncloud/ocis/v2/ocis-pkg/keycloak"
	"github.com/stretchr/testify/assert"
)

func TestParseUserActions(t *testing.T) {
	defaults := []kcpkg.UserAction{kcpkg.UserActionUpdatePassword, kcpkg.UserActionVerifyEmail}
	tests := []struct {
		name string
		in   []string
		want []kcpkg.UserAction
	}{
		{"nil falls back to defaults", nil, defaults},
		{"empty falls back to defaults", []string{}, defaults},
		{"whitespace-only entries fall back to defaults", []string{"", "   ", "\t"}, defaults},
		{
			"entries are trimmed",
			[]string{" UPDATE_PASSWORD ", " VERIFY_EMAIL "},
			[]kcpkg.UserAction{kcpkg.UserActionUpdatePassword, kcpkg.UserActionVerifyEmail},
		},
		{
			"unknown values are passed through unchanged",
			[]string{"FOO"},
			[]kcpkg.UserAction{kcpkg.UserAction("FOO")},
		},
		{
			"case is preserved",
			[]string{"update_password"},
			[]kcpkg.UserAction{kcpkg.UserAction("update_password")},
		},
		{
			"empty entries between valid ones are dropped",
			[]string{"UPDATE_PASSWORD", " ", "FOO"},
			[]kcpkg.UserAction{kcpkg.UserActionUpdatePassword, kcpkg.UserAction("FOO")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseUserActions(tt.in))
		})
	}
}
