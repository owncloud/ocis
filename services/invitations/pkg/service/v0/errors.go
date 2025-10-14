package service

import "errors"

var (
	ErrNotFound     = errors.New("query target not found")
	ErrBadRequest   = errors.New("bad request")
	ErrMissingEmail = errors.New("missing email address")
	ErrBackend      = errors.New("backend error")
)
