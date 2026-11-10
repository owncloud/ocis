package service

import (
	"context"
	"errors"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
	"github.com/stretchr/testify/assert"
)

type fakeBackend struct {
	createErr error
	mailErr   error
	deleteErr error
	canMail   bool

	mailedIDs  []string
	deletedIDs []string
}

func (f *fakeBackend) CreateUser(context.Context, *invitations.Invitation) (string, error) {
	return "kc-id", f.createErr
}
func (f *fakeBackend) CanSendMail() bool { return f.canMail }
func (f *fakeBackend) SendMail(_ context.Context, id string) error {
	f.mailedIDs = append(f.mailedIDs, id)
	return f.mailErr
}
func (f *fakeBackend) DeleteUser(_ context.Context, id string) error {
	f.deletedIDs = append(f.deletedIDs, id)
	return f.deleteErr
}

func newTestSvc(b Backend) svc {
	return svc{log: log.NopLogger(), backend: b}
}

func TestInvite_MailFailureRollsBackUser(t *testing.T) {
	mailErr := errors.New("required action FOO not found")
	b := &fakeBackend{canMail: true, mailErr: mailErr}

	inv, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})

	assert.Nil(t, inv)
	assert.ErrorIs(t, err, ErrBackend)
	assert.ErrorContains(t, err, mailErr.Error())
	assert.Equal(t, []string{"kc-id"}, b.deletedIDs)
}

func TestInvite_MailFailureAndDeleteFailureReturnsMailError(t *testing.T) {
	mailErr := errors.New("mail failed")
	b := &fakeBackend{canMail: true, mailErr: mailErr, deleteErr: errors.New("delete failed")}

	inv, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})

	assert.Nil(t, inv)
	assert.ErrorIs(t, err, ErrBackend)
	assert.ErrorContains(t, err, "mail failed")
	assert.NotContains(t, err.Error(), "delete failed")
	assert.Equal(t, []string{"kc-id"}, b.deletedIDs)
}

func TestInvite_SuccessDoesNotDelete(t *testing.T) {
	b := &fakeBackend{canMail: true}
	inv, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})
	assert.NoError(t, err)
	assert.NotNil(t, inv)
	assert.Empty(t, b.deletedIDs)
}

func TestInvite_CreateFailureDoesNotDelete(t *testing.T) {
	b := &fakeBackend{canMail: true, createErr: errors.New("409")}
	_, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})
	assert.ErrorIs(t, err, ErrBackend)
	assert.Empty(t, b.deletedIDs)
	assert.Empty(t, b.mailedIDs)
}

// ctxBackend records the context DeleteUser receives and cancels the request
// context while the mail is being sent.
type ctxBackend struct {
	fakeBackend
	cancel       context.CancelFunc
	deleteCalled bool
	deleteErr    error
	hasDeadline  bool
}

func (c *ctxBackend) SendMail(ctx context.Context, id string) error {
	c.cancel()
	return context.Canceled
}

func (c *ctxBackend) DeleteUser(ctx context.Context, id string) error {
	c.deleteCalled = true
	c.deleteErr = ctx.Err()
	_, c.hasDeadline = ctx.Deadline()
	return ctx.Err()
}

func TestInvite_RollbackSurvivesCancelledRequestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	b := &ctxBackend{fakeBackend: fakeBackend{canMail: true}, cancel: cancel}

	_, err := newTestSvc(b).Invite(ctx, &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})

	assert.ErrorIs(t, err, ErrBackend)
	assert.True(t, b.deleteCalled)
	assert.NoError(t, b.deleteErr, "rollback must not use the cancelled request context")
	assert.True(t, b.hasDeadline, "rollback must be bounded by a timeout")
}
