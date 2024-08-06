package unifiedrole

import (
	"errors"
)

var (
	// ErrUnknownUnifiedRole is returned when an unknown unified role is requested.
	ErrUnknownUnifiedRole = errors.New("unknown unified role, check if the role is enabled")

	// ErrTooManyResults is returned when a filter returns too many results.
	ErrTooManyResults = errors.New("too many results, consider using a more specific filter")
)
