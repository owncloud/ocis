// Package keycloak offers an invitation backend for the invitation service.
// TODO: Maybe move this outside of the invitation service and make it more generic?

package keycloak_test

import (
	"context"
	"testing"

	kcpkg "github.com/owncloud/ocis/v2/ocis-pkg/keycloak"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/backends/keycloak"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/backends/keycloak/mocks"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		name        string
		args        args
		want        string
		clientMocks []mockInputs
		assertion   assert.ErrorAssertionFunc
	}{
		{
			name: "Test without diplay name",
			args: args{
				invitation: &invitations.Invitation{
					InvitedUserEmailAddress: "test@example.org",
				},
			},
			want: "test-id",
			clientMocks: []mockInputs{
				{
					funcName: "CreateUser",
					args: []interface{}{
						mock.Anything,
						userRealm,
						mock.Anything, // can't match on the user because it generates a UUID internally.
						[]kcpkg.UserAction{
							kcpkg.UserActionUpdatePassword,
							kcpkg.UserActionVerifyEmail,
						},
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
			c := &mocks.Client{}
			for _, m := range tt.clientMocks {
				c.On(m.funcName, m.args...).Return(m.returns...)
			}
			b := keycloak.NewWithClient(log.NopLogger(), c, userRealm)
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
		name        string
		args        args
		clientMocks []mockInputs
		assertion   assert.ErrorAssertionFunc
	}{
		{
			name: "Mail successfully sent",
			args: args{
				id: "test-id",
			},
			clientMocks: []mockInputs{
				{
					funcName: "SendActionsMail",
					args: []interface{}{
						mock.Anything,
						userRealm,
						"test-id",
						[]kcpkg.UserAction{
							kcpkg.UserActionUpdatePassword,
							kcpkg.UserActionVerifyEmail,
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
			c := &mocks.Client{}
			for _, m := range tt.clientMocks {
				c.On(m.funcName, m.args...).Return(m.returns...)
			}
			b := keycloak.NewWithClient(log.NopLogger(), c, userRealm)
			tt.assertion(t, b.SendMail(ctx, tt.args.id))
		})
	}
}
