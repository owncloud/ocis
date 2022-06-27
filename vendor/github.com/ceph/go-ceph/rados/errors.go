package rados

/*
#include <errno.h>
*/
import "C"

import (
	"errors"

	"github.com/ceph/go-ceph/internal/errutil"
)

// radosError represents an error condition returned from the Ceph RADOS APIs.
type radosError int

// Error returns the error string for the radosError type.
func (e radosError) Error() string {
	return errutil.FormatErrorCode("rados", int(e))
}

func (e radosError) ErrorCode() int {
	return int(e)
}

func getError(e C.int) error {
	if e == 0 {
		return nil
	}
	return radosError(e)
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

// Public go errors:

var (
	// ErrNotConnected is returned when functions are called
	// without a RADOS connection.
	ErrNotConnected = errors.New("RADOS not connected")
	// ErrEmptyArgument may be returned if a function argument is passed
	// a zero-length slice or map.
	ErrEmptyArgument = errors.New("Argument must contain at least one item")
	// ErrInvalidIOContext may be returned if an api call requires an IOContext
	// but IOContext is not ready for use.
	ErrInvalidIOContext = errors.New("IOContext is not ready for use")
	// ErrOperationIncomplete is returned from write op or read op steps for
	// which the operation has not been performed yet.
	ErrOperationIncomplete = errors.New("Operation has not been performed yet")
)

// Public radosErrors:

const (
	// ErrNotFound indicates a missing resource.
	ErrNotFound = radosError(-C.ENOENT)
	// ErrPermissionDenied indicates a permissions issue.
	ErrPermissionDenied = radosError(-C.EPERM)
	// ErrObjectExists indicates that an exclusive object creation failed.
	ErrObjectExists = radosError(-C.EEXIST)

	// RadosErrorNotFound indicates a missing resource.
	//
	// Deprecated: use ErrNotFound instead
	RadosErrorNotFound = ErrNotFound
	// RadosErrorPermissionDenied indicates a permissions issue.
	//
	// Deprecated: use ErrPermissionDenied instead
	RadosErrorPermissionDenied = ErrPermissionDenied
)

// Private errors:

const (
	errNameTooLong = radosError(-C.ENAMETOOLONG)

	errRange = radosError(-C.ERANGE)
)
