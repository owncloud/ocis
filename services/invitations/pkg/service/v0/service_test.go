package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
	"github.com/stretchr/testify/assert"
)

type fakeBackend struct {
	createErr error
	mailErr   error
	deleteErr error
	canMail   bool
	emptyID   bool

	mailedIDs  []string
	deletedIDs []string
}

func (f *fakeBackend) CreateUser(_ context.Context, inv *invitations.Invitation) (string, error) {
	if f.createErr != nil {
		return "", f.createErr
	}
	inv.InvitedUser = &libregraph.User{}
	if f.emptyID {
		return "", nil
	}
	return "kc-id", nil
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
	mailErr := fmt.Errorf("%w: required action FOO not found", invitations.ErrRejected)
	b := &fakeBackend{canMail: true, mailErr: mailErr}

	inv, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})

	assert.Nil(t, inv)
	assert.ErrorIs(t, err, ErrBackend)
	assert.ErrorContains(t, err, mailErr.Error())
	assert.Equal(t, []string{"kc-id"}, b.deletedIDs)
}

func TestInvite_MailFailureAndDeleteFailureReturnsMailError(t *testing.T) {
	mailErr := fmt.Errorf("%w: mail failed", invitations.ErrRejected)
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
	deleteCtxErr error
	hasDeadline  bool
}

// SendMail models Keycloak answering with an error just as the request
// context gets cancelled.
func (c *ctxBackend) SendMail(ctx context.Context, id string) error {
	c.cancel()
	return fmt.Errorf("%w: 500", invitations.ErrRejected)
}

func (c *ctxBackend) DeleteUser(ctx context.Context, id string) error {
	c.deleteCalled = true
	c.deleteCtxErr = ctx.Err()
	_, c.hasDeadline = ctx.Deadline()
	return ctx.Err()
}

func TestInvite_RollbackSurvivesCancelledRequestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	b := &ctxBackend{fakeBackend: fakeBackend{canMail: true}, cancel: cancel}

	_, err := newTestSvc(b).Invite(ctx, &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})

	assert.ErrorIs(t, err, ErrBackend)
	assert.True(t, b.deleteCalled)
	assert.NoError(t, b.deleteCtxErr, "rollback must not use the cancelled request context")
	assert.True(t, b.hasDeadline, "rollback must be bounded by a timeout")
}

func TestInvite_MailFailureWithUnknownOutcomeKeepsUser(t *testing.T) {
	for name, mailErr := range map[string]error{
		"context canceled":  context.Canceled,
		"deadline exceeded": context.DeadlineExceeded,
		"transport error":   errors.New("connection reset by peer"),
	} {
		t.Run(name, func(t *testing.T) {
			b := &fakeBackend{canMail: true, mailErr: mailErr}
			inv, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})
			assert.Nil(t, inv)
			assert.ErrorIs(t, err, ErrBackend)
			assert.ErrorContains(t, err, mailErr.Error())
			assert.Empty(t, b.deletedIDs)
		})
	}
}

func TestInvite_EmptyIDIsNeverDeleted(t *testing.T) {
	b := &fakeBackend{canMail: true, emptyID: true, mailErr: fmt.Errorf("%w: 500", invitations.ErrRejected)}
	_, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})
	assert.ErrorIs(t, err, ErrBackend)
	assert.Empty(t, b.deletedIDs)
}

func TestInvite_SuccessfulRollbackClearsInvitedUser(t *testing.T) {
	b := &fakeBackend{canMail: true, mailErr: fmt.Errorf("%w: 500", invitations.ErrRejected)}
	inv := &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"}
	_, err := newTestSvc(b).Invite(context.Background(), inv)
	assert.Error(t, err)
	assert.Equal(t, []string{"kc-id"}, b.deletedIDs)
	assert.Nil(t, inv.InvitedUser)
}
