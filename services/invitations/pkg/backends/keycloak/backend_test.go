// Package keycloak offers an invitation backend for the invitation service.
// TODO: Maybe move this outside of the invitation service and make it more generic?

package keycloak_test

import (
	"context"
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/backends/keycloak"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/mock"
)

const (
	clientID     = "test-id"
	clientSecret = "test-secret"
	clientRealm  = "client-realm"
	userRealm    = "user-realm"
	jwtToken     = "test-token"
)

func TestBackend_CreateUser(t *testing.T) {
	type args struct {
		invitation *invitations.Invitation
	}
	type mockInputs struct {
		funcName string
		args     []interface{}
		returns  []interface{}
	}
	tests := []struct {
		name          string
		args          args
		want          string
		keycloakMocks []mockInputs
		assertion     assert.ErrorAssertionFunc
	}{
		{
			name: "Test without diplay name",
			args: args{
				invitation: &invitations.Invitation{
					InvitedUserEmailAddress: "test@example.org",
				},
			},
			want: "test-id",
			keycloakMocks: []mockInputs{
				{
					funcName: "LoginClient",
					args: []interface{}{
						mock.Anything,
						clientID,
						clientSecret,
						clientRealm,
					},
					returns: []interface{}{
						&gocloak.JWT{
							AccessToken: jwtToken,
						},
						nil,
					},
				},
				{
					funcName: "RetrospectToken",
					args: []interface{}{
						mock.Anything,
						jwtToken,
						clientID,
						clientSecret,
						clientRealm,
					},
					returns: []interface{}{
						&gocloak.IntroSpectTokenResult{
							Active: gocloak.BoolP(true),
						},
						nil,
					},
				},
				{
					funcName: "CreateUser",
					args: []interface{}{
						mock.Anything,
						jwtToken,
						userRealm,
						mock.Anything, // can't match on the user because it generates a UUID internally.
						// might be worth refactoring the UUID generation to outside of the func
					},
					returns: []interface{}{
						"test-id",
						nil,
					},
				},
			},
			assertion: func(t assert.TestingT, err error, args ...interface{}) bool {
				return assert.Nil(t, err, args...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			c := &mocks.GoCloak{}
			for _, m := range tt.keycloakMocks {
				c.On(m.funcName, m.args...).Return(m.returns...)
			}
			b := keycloak.NewWithClient(log.NopLogger(), c, clientID, clientSecret, clientRealm, userRealm)
			got, err := b.CreateUser(ctx, tt.args.invitation)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBackend_SendMail(t *testing.T) {
	type args struct {
		id string
	}
	type mockInputs struct {
		funcName string
		args     []interface{}
		returns  []interface{}
	}
	tests := []struct {
		name          string
		args          args
		keycloakMocks []mockInputs
		assertion     assert.ErrorAssertionFunc
	}{
		{
			name: "Mail successfully sent",
			args: args{
				id: "test-id",
			},
			keycloakMocks: []mockInputs{
				{
					funcName: "LoginClient",
					args: []interface{}{
						mock.Anything,
						clientID,
						clientSecret,
						clientRealm,
					},
					returns: []interface{}{
						&gocloak.JWT{
							AccessToken: jwtToken,
						},
						nil,
					},
				},
				{
					funcName: "RetrospectToken",
					args: []interface{}{
						mock.Anything,
						jwtToken,
						clientID,
						clientSecret,
						clientRealm,
					},
					returns: []interface{}{
						&gocloak.IntroSpectTokenResult{
							Active: gocloak.BoolP(true),
						},
						nil,
					},
				},
				{
					funcName: "ExecuteActionsEmail",
					args: []interface{}{
						mock.Anything,
						jwtToken,
						userRealm,
						gocloak.ExecuteActionsEmail{
							UserID:  gocloak.StringP("test-id"),
							Actions: &[]string{"UPDATE_PASSWORD", "VERIFY_EMAIL"},
						},
					},
					returns: []interface{}{
						nil,
					},
				},
			},
			assertion: func(t assert.TestingT, err error, args ...interface{}) bool {
				return assert.Nil(t, err, args...)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			c := &mocks.GoCloak{}
			for _, m := range tt.keycloakMocks {
				c.On(m.funcName, m.args...).Return(m.returns...)
			}
			b := keycloak.NewWithClient(log.NopLogger(), c, clientID, clientSecret, clientRealm, userRealm)
			tt.assertion(t, b.SendMail(ctx, tt.args.id))
		})
	}
}
