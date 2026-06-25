package keycloak_test

import (
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/keycloak"
	"github.com/stretchr/testify/assert"
)

func TestUserActionFromString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   keycloak.UserAction
		wantOk bool
	}{
		{name: "known uppercase action", input: "UPDATE_PASSWORD", want: keycloak.UserActionUpdatePassword, wantOk: true},
		{name: "known verify email", input: "VERIFY_EMAIL", want: keycloak.UserActionVerifyEmail, wantOk: true},
		{name: "known lowercase provider id", input: "delete_account", want: keycloak.UserActionDeleteAccount, wantOk: true},
		{name: "known hyphenated provider id", input: "webauthn-register", want: keycloak.UserActionWebauthnRegister, wantOk: true},
		{name: "unknown action", input: "NOT_A_REAL_ACTION", wantOk: false},
		{name: "empty string", input: "", wantOk: false},
		{name: "wrong case is not matched", input: "update_password", wantOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := keycloak.UserActionFromString(tt.input)
			assert.Equal(t, tt.wantOk, ok)
			if tt.wantOk {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
