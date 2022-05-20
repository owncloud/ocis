package cs3

import (
	"context"
	"time"

	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
)

// createSession creates a new Session without the server using the provided
// data.
func createSession(ctx context.Context, u *cs3user.User) *cs3Session {

	if ctx == nil {
		ctx = context.Background()
	}

	sessionCtx, cancel := context.WithCancel(ctx)
	s := &cs3Session{
		ctx:       sessionCtx,
		u:         u,
		ctxCancel: cancel,
	}

	s.when = time.Now()

	return s
}

type cs3Session struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	u         *cs3user.User
	when      time.Time
}

func (s *cs3Session) User() *cs3user.User {
	return s.u
}
