package middleware

import (
	"errors"
)

var (
	// ErrInvalidToken is returned when the request token is invalid.
	ErrInvalidToken = errors.New("invalid or missing token")

	// ErrUnauthorized is returned if the request is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInternal is returned if something went wrong
	ErrInternal     = errors.New("internal error")
)
