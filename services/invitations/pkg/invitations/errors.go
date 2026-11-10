package invitations

import "errors"

// ErrRejected matches errors where the backend certainly refused the request,
// so that the requested change was not applied. Errors without it (transport
// failures, timeouts, cancelled contexts) leave the outcome unknown. Backends
// mark an error with Rejected.
var ErrRejected = errors.New("request rejected by backend")

// rejectedError marks err as a definite rejection without changing its
// message, which is passed on to API clients.
type rejectedError struct {
	err error
}

// Rejected returns err marked as a definite rejection: errors.Is(err,
// ErrRejected) reports true, while Error() returns the message of err
// unchanged. It returns nil if err is nil.
func Rejected(err error) error {
	if err == nil {
		return nil
	}
	return &rejectedError{err: err}
}

func (e *rejectedError) Error() string { return e.err.Error() }

func (e *rejectedError) Unwrap() error { return e.err }

func (e *rejectedError) Is(target error) bool { return target == ErrRejected }
