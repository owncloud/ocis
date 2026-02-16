package cs3

import (
	"time"

	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
)

// createSession creates a new Session without the server using the provided
// data.
func createSession(u *cs3user.User) *cs3Session {
	s := &cs3Session{
		u: u,
	}

	s.when = time.Now()

	return s
}

type cs3Session struct {
	u    *cs3user.User
	when time.Time
}

// User returns the cs3 user of the session
func (s *cs3Session) User() *cs3user.User {
	return s.u
}
