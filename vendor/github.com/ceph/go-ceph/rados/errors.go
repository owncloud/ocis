package rados

/*
#include <errno.h>
*/
import "C"

import (
	"errors"

	"github.com/ceph/go-ceph/internal/errutil"
)

// Public go errors:

var (
	// ErrNotConnected is returned when functions are called
	// without a RADOS connection.
	ErrNotConnected = getError(-C.ENOTCONN)
	// ErrEmptyArgument may be returned if a function argument is passed
	// a zero-length slice or map.
	ErrEmptyArgument = errors.New("Argument must contain at least one item")
	// ErrInvalidIOContext may be returned if an api call requires an IOContext
	// but IOContext is not ready for use.
	ErrInvalidIOContext = errors.New("IOContext is not ready for use")
	// ErrOperationIncomplete is returned from write op or read op steps for
	// which the operation has not been performed yet.
	ErrOperationIncomplete = errors.New("Operation has not been performed yet")

	// ErrNotFound indicates a missing resource.
	ErrNotFound = getError(-C.ENOENT)
	// ErrPermissionDenied indicates a permissions issue.
	ErrPermissionDenied = getError(-C.EPERM)
	// ErrObjectExists indicates that an exclusive object creation failed.
	ErrObjectExists = getError(-C.EEXIST)

	// RadosErrorNotFound indicates a missing resource.
	//
	// Deprecated: use ErrNotFound instead
	RadosErrorNotFound = ErrNotFound
	// RadosErrorPermissionDenied indicates a permissions issue.
	//
	// Deprecated: use ErrPermissionDenied instead
	RadosErrorPermissionDenied = ErrPermissionDenied

	// Private errors:

	errNameTooLong = getError(-C.ENAMETOOLONG)
	errRange       = getError(-C.ERANGE)
)

func getError(errno C.int) error {
	return errutil.GetError("rados", int(errno))
}

// getErrorIfNegative converts a ceph return code to error if negative.
// This is useful for functions that return a usable positive value on
// success but a negative error number on error.
func getErrorIfNegative(ret C.int) error {
	if ret >= 0 {
		return nil
	}
	return getError(ret)
}
