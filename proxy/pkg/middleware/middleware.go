package middleware

import (
	"errors"
)

var (
	ErrInvalidToken = errors.New("invalid or missing token")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInternal     = errors.New("internal error")
)
