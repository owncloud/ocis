package keycloak_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/owncloud/ocis/v2/ocis-pkg/keycloak"
	"github.com/stretchr/testify/assert"
)

// fakeGoCloak implements keycloak.GoCloak for the token and mail calls.
type fakeGoCloak struct {
	keycloak.GoCloak
	loginErr      error
	retrospectErr error
	active        bool
	mailErr       error
	mailCalled    bool
}

func (f *fakeGoCloak) LoginClient(context.Context, string, string, string, ...string) (*gocloak.JWT, error) {
	if f.loginErr != nil {
		return nil, f.loginErr
	}
	return &gocloak.JWT{AccessToken: "token"}, nil
}

func (f *fakeGoCloak) RetrospectToken(context.Context, string, string, string, string) (*gocloak.IntroSpectTokenResult, error) {
	if f.retrospectErr != nil {
		return nil, f.retrospectErr
	}
	return &gocloak.IntroSpectTokenResult{Active: &f.active}, nil
}

func (f *fakeGoCloak) ExecuteActionsEmail(context.Context, string, string, gocloak.ExecuteActionsEmail) error {
	f.mailCalled = true
	return f.mailErr
}

func TestSendActionsMail_TokenFailuresAreTokenErrors(t *testing.T) {
	cause := errors.New("connection refused")
	tests := []struct {
		name string
		gc   *fakeGoCloak
		msg  string
	}{
		{"login fails", &fakeGoCloak{loginErr: cause}, "failed to get token: connection refused"},
		{"retrospect fails", &fakeGoCloak{retrospectErr: cause}, "failed to retrospect token: connection refused"},
		{"token not active", &fakeGoCloak{active: false}, "token is not active"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := keycloak.NewWithClient(tt.gc, "https://kc", "id", "secret", "realm")
			err := c.SendActionsMail(context.Background(), "realm", "user", nil)
			var tokenErr *keycloak.TokenError
			assert.ErrorAs(t, err, &tokenErr)
			assert.EqualError(t, err, tt.msg)
			assert.False(t, tt.gc.mailCalled)
		})
	}
}

func TestSendActionsMail_MailErrorIsNoTokenError(t *testing.T) {
	gc := &fakeGoCloak{active: true, mailErr: &gocloak.APIError{Code: 500}}
	c := keycloak.NewWithClient(gc, "https://kc", "id", "secret", "realm")
	err := c.SendActionsMail(context.Background(), "realm", "user", nil)
	var tokenErr *keycloak.TokenError
	assert.False(t, errors.As(err, &tokenErr))
	assert.True(t, gc.mailCalled)
}
