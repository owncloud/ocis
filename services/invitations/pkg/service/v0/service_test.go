package service

import (
	"bytes"
	"context"
	"errors"
	"testing"

	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/invitations/pkg/invitations"
	"github.com/rs/zerolog"
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
	mailErr := invitations.Rejected(errors.New("required action FOO not found"))
	b := &fakeBackend{canMail: true, mailErr: mailErr}

	inv, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})

	assert.Nil(t, inv)
	assert.ErrorIs(t, err, ErrBackend)
	assert.ErrorContains(t, err, mailErr.Error())
	assert.Equal(t, []string{"kc-id"}, b.deletedIDs)
}

func TestInvite_MailFailureAndDeleteFailureReturnsMailError(t *testing.T) {
	mailErr := invitations.Rejected(errors.New("mail failed"))
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

// ctxBackend records the contexts the backend calls receive. cancelIn names
// the call during which the request context is cancelled.
type ctxBackend struct {
	fakeBackend
	cancel   context.CancelFunc
	cancelIn string

	ctxErrs   map[string]error
	deadlines map[string]bool
}

func (c *ctxBackend) record(call string, ctx context.Context) {
	if c.cancelIn == call {
		c.cancel()
	}
	if c.ctxErrs == nil {
		c.ctxErrs, c.deadlines = map[string]error{}, map[string]bool{}
	}
	c.ctxErrs[call] = ctx.Err()
	_, c.deadlines[call] = ctx.Deadline()
}

func (c *ctxBackend) CreateUser(ctx context.Context, inv *invitations.Invitation) (string, error) {
	c.record("CreateUser", ctx)
	return c.fakeBackend.CreateUser(ctx, inv)
}

func (c *ctxBackend) SendMail(ctx context.Context, id string) error {
	c.record("SendMail", ctx)
	return c.fakeBackend.SendMail(ctx, id)
}

func (c *ctxBackend) DeleteUser(ctx context.Context, id string) error {
	c.record("DeleteUser", ctx)
	return c.fakeBackend.DeleteUser(ctx, id)
}

func TestInvite_BackendCallsAreDetachedFromTheRequestContext(t *testing.T) {
	for _, cancelIn := range []string{"before", "CreateUser", "SendMail"} {
		t.Run("request cancelled "+cancelIn, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if cancelIn == "before" {
				cancel()
			}
			b := &ctxBackend{
				fakeBackend: fakeBackend{canMail: true, mailErr: invitations.Rejected(errors.New("500"))},
				cancel:      cancel,
				cancelIn:    cancelIn,
			}

			_, err := newTestSvc(b).Invite(ctx, &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})

			assert.ErrorIs(t, err, ErrBackend)
			for _, call := range []string{"CreateUser", "SendMail", "DeleteUser"} {
				if assert.Contains(t, b.ctxErrs, call, "%s must be called", call) {
					assert.NoError(t, b.ctxErrs[call], "%s must not get the cancelled request context", call)
					assert.True(t, b.deadlines[call], "%s must be bounded by a timeout", call)
				}
			}
			assert.Equal(t, []string{"kc-id"}, b.deletedIDs)
		})
	}
}

func TestInvite_ResponseMessageHidesClassification(t *testing.T) {
	orig := errors.New("400 Bad Request: Provided invalid required actions")
	b := &fakeBackend{canMail: true, mailErr: invitations.Rejected(orig)}

	_, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})

	// Byte-identical to the message before the rollback was introduced.
	assert.EqualError(t, err, "backend error: 400 Bad Request: Provided invalid required actions")
	assert.Equal(t, []string{"kc-id"}, b.deletedIDs)
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

func TestInvite_EmptyIDFailsWithoutMailOrDelete(t *testing.T) {
	b := &fakeBackend{canMail: true, emptyID: true}
	inv, err := newTestSvc(b).Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})
	assert.Nil(t, inv)
	assert.ErrorIs(t, err, ErrBackend)
	assert.Empty(t, b.mailedIDs)
	assert.Empty(t, b.deletedIDs)
}

func TestInvite_UnknownOutcomeTellsOperatorsWhatToDo(t *testing.T) {
	// Some dependency raises the global zerolog level in an init function.
	previous := zerolog.GlobalLevel()
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	t.Cleanup(func() { zerolog.SetGlobalLevel(previous) })
	var buf bytes.Buffer
	b := &fakeBackend{canMail: true, mailErr: context.DeadlineExceeded}
	s := svc{log: log.Logger{Logger: zerolog.New(&buf)}, backend: b}

	_, err := s.Invite(context.Background(), &invitations.Invitation{InvitedUserEmailAddress: "a@example.org"})

	assert.ErrorIs(t, err, ErrBackend)
	assert.Empty(t, b.deletedIDs)
	out := buf.String()
	assert.Contains(t, out, `"level":"error"`)
	assert.Contains(t, out, `"userID":"kc-id"`)
	assert.Contains(t, out, `"email":"a@example.org"`)
	assert.Contains(t, out, "was kept")
	assert.Contains(t, out, "may or may not have been sent")
	assert.Contains(t, out, "remove the guest a@example.org (user id kc-id) in Keycloak before inviting")
}
