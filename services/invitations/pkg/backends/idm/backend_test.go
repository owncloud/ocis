package idm_test

import (
	"context"
	"errors"
	"testing"

	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/backends/idm"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeProvisioner records the user it was asked to create and returns a
// configurable result, standing in for *identity.LDAP in unit tests.
type fakeProvisioner struct {
	gotUser libregraph.User
	err     error
}

func (f *fakeProvisioner) CreateUser(_ context.Context, user libregraph.User) (*libregraph.User, error) {
	f.gotUser = user
	if f.err != nil {
		return nil, f.err
	}
	created := user // the directory backend echoes the created user (id preserved)
	return &created, nil
}

func TestBackend_CreateUser_ProvisionsGuest(t *testing.T) {
	f := &fakeProvisioner{}
	b := idm.New(log.NopLogger(), f)
	inv := &invitations.Invitation{InvitedUserEmailAddress: "guest@example.org"}

	id, err := b.CreateUser(context.Background(), inv)
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// the provisioner received a Guest user built from the invitation
	assert.Equal(t, "guest@example.org", f.gotUser.GetMail())
	assert.Equal(t, "Guest", f.gotUser.GetUserType())
	assert.Equal(t, "guest@example.org", f.gotUser.GetOnPremisesSamAccountName())
	assert.True(t, f.gotUser.GetAccountEnabled())

	// the invitation now carries the created user as invitedUser
	require.NotNil(t, inv.InvitedUser)
	assert.Equal(t, id, inv.InvitedUser.GetId())
	assert.Equal(t, "guest@example.org", inv.InvitedUser.GetMail())
}

func TestBackend_CreateUser_DisplayNameFallsBackToEmail(t *testing.T) {
	f := &fakeProvisioner{}
	b := idm.New(log.NopLogger(), f)

	// no display name -> falls back to the email
	_, err := b.CreateUser(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "g@example.org"})
	require.NoError(t, err)
	assert.Equal(t, "g@example.org", f.gotUser.GetDisplayName())

	// explicit display name is kept
	_, err = b.CreateUser(context.Background(), &invitations.Invitation{
		InvitedUserEmailAddress: "g@example.org",
		InvitedUserDisplayName:  "Guest Person",
	})
	require.NoError(t, err)
	assert.Equal(t, "Guest Person", f.gotUser.GetDisplayName())
}

func TestBackend_CreateUser_ProvisionError(t *testing.T) {
	f := &fakeProvisioner{err: errors.New("ldap write failed")}
	b := idm.New(log.NopLogger(), f)
	inv := &invitations.Invitation{InvitedUserEmailAddress: "g@example.org"}

	_, err := b.CreateUser(context.Background(), inv)
	require.Error(t, err)
	assert.Nil(t, inv.InvitedUser)
}

func TestBackend_CanSendMail_IsFalse(t *testing.T) {
	b := idm.New(log.NopLogger(), &fakeProvisioner{})
	assert.False(t, b.CanSendMail())
	assert.NoError(t, b.SendMail(context.Background(), "some-id"))
}
