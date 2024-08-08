package unifiedrole

import (
	"errors"
)

var (
	// ErrUnknownRole is returned when an unknown unified role is requested.
	ErrUnknownRole = errors.New("unknown role, check if the role is enabled")
)
