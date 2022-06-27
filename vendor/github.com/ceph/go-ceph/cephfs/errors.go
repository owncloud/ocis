package cephfs

/*
#include <errno.h>
*/
import "C"

import (
	"errors"

	"github.com/ceph/go-ceph/internal/errutil"
)

// cephFSError represents an error condition returned from the CephFS APIs.
type cephFSError int

// Error returns the error string for the cephFSError type.
func (e cephFSError) Error() string {
	return errutil.FormatErrorCode("cephfs", int(e))
}

func (e cephFSError) ErrorCode() int {
	return int(e)
}

func getError(e C.int) error {
	if e == 0 {
		return nil
	}
	return cephFSError(e)
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
	// ErrEmptyArgument may be returned if a function argument is passed
	// a zero-length slice or map.
	ErrEmptyArgument = errors.New("Argument must contain at least one item")
)

// Public CephFSErrors:

const (
	// ErrNotConnected may be returned when client is not connected
	// to a cluster.
	ErrNotConnected = cephFSError(-C.ENOTCONN)
)

// Private errors:

const (
	errInvalid     = cephFSError(-C.EINVAL)
	errNameTooLong = cephFSError(-C.ENAMETOOLONG)
	errNoEntry     = cephFSError(-C.ENOENT)
	errRange       = cephFSError(-C.ERANGE)
)
