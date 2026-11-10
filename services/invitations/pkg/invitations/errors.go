package invitations

import "errors"

// ErrRejected is wrapped by backends around errors where the backend certainly
// refused the request, so that the requested change was not applied. Errors
// without it (transport failures, timeouts, cancelled contexts) leave the
// outcome unknown.
var ErrRejected = errors.New("request rejected by backend")
