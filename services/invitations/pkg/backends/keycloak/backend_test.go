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
		userActions []kcpkg.UserAction
		want        string
		clientMocks []mockInputs
		assertion   assert.ErrorAssertionFunc
	}{
		{
			name: "nil actions fall back to the default actions",
			args: args{
				invitation: &invitations.Invitation{
					InvitedUserEmailAddress: "test@example.org",
				},
			},
			userActions: nil,
			want:        "test-id",
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
		{
			name: "configured actions are passed through to keycloak",
			args: args{
				invitation: &invitations.Invitation{
					InvitedUserEmailAddress: "test@example.org",
				},
			},
			userActions: []kcpkg.UserAction{kcpkg.UserActionUpdatePassword},
			want:        "test-id",
			clientMocks: []mockInputs{
				{
					funcName: "CreateUser",
					args: []interface{}{
						mock.Anything,
						userRealm,
						mock.Anything,
						[]kcpkg.UserAction{
							kcpkg.UserActionUpdatePassword,
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
			b := keycloak.NewWithClient(log.NopLogger(), c, userRealm, tt.userActions)
			got, err := b.CreateUser(ctx, tt.args.invitation)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
			// On success the created user is recorded on the invitation as invitedUser.
			if err == nil && assert.NotNil(t, tt.args.invitation.InvitedUser) {
				assert.Equal(t, tt.args.invitation.InvitedUserEmailAddress, tt.args.invitation.InvitedUser.GetMail())
				assert.Equal(t, "Guest", tt.args.invitation.InvitedUser.GetUserType())
				assert.NotEmpty(t, tt.args.invitation.InvitedUser.GetId())
			}
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
		userActions []kcpkg.UserAction
		clientMocks []mockInputs
		assertion   assert.ErrorAssertionFunc
	}{
		{
			name: "mail successfully sent with default actions",
			args: args{
				id: "test-id",
			},
			userActions: nil,
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
		{
			name: "mail sent with only the configured action",
			args: args{
				id: "test-id",
			},
			userActions: []kcpkg.UserAction{kcpkg.UserActionUpdatePassword},
			clientMocks: []mockInputs{
				{
					funcName: "SendActionsMail",
					args: []interface{}{
						mock.Anything,
						userRealm,
						"test-id",
						[]kcpkg.UserAction{
							kcpkg.UserActionUpdatePassword,
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
			b := keycloak.NewWithClient(log.NopLogger(), c, userRealm, tt.userActions)
			tt.assertion(t, b.SendMail(ctx, tt.args.id))
		})
	}
}
