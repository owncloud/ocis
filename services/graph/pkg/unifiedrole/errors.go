package unifiedrole

import (
	"errors"
)

var (
	// ErrUnknownUnifiedRole is returned when an unknown unified role is requested.
	ErrUnknownUnifiedRole = errors.New("unknown unified role, check if the role is enabled")
)
