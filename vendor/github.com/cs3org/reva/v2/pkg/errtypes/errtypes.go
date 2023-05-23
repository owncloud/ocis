// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

// Package errtypes contains definitions for common errors.
// It would have nice to call this package errors, err or error
// but errors clashes with github.com/pkg/errors, err is used for any error variable
// and error is a reserved word :)
package errtypes

import (
	"strings"

	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
)

// NotFound is the error to use when a something is not found.
type NotFound string

func (e NotFound) Error() string { return "error: not found: " + string(e) }

// IsNotFound implements the IsNotFound interface.
func (e NotFound) IsNotFound() {}

// InternalError is the error to use when we really don't know what happened. Use with care
type InternalError string

func (e InternalError) Error() string { return "internal error: " + string(e) }

// IsInternalError implements the IsInternalError interface.
func (e InternalError) IsInternalError() {}

// PermissionDenied is the error to use when a resource cannot be access because of missing permissions.
type PermissionDenied string

func (e PermissionDenied) Error() string { return "error: permission denied: " + string(e) }

// IsPermissionDenied implements the IsPermissionDenied interface.
func (e PermissionDenied) IsPermissionDenied() {}

// Locked is the error to use when a resource cannot be modified because of a lock.
type Locked string

func (e Locked) Error() string { return "error: locked by " + string(e) }

// LockID returns the lock ID that caused this error
func (e Locked) LockID() string {
	return string(e)
}

// IsLocked implements the IsLocked interface.
func (e Locked) IsLocked() {}

// Aborted is the error to use when a client should retry at a higher level
// (e.g., when a client-specified test-and-set fails, indicating the
// client should restart a read-modify-write sequence) request fails
// because a requested etag or lock ID mismatches.
//
// HTTP Mapping: 412 Precondition Failed
type Aborted string

func (e Aborted) Error() string { return "error: aborted: " + string(e) }

// IsAborted implements the IsAborted interface.
func (e Aborted) IsAborted() {}

// PreconditionFailed is the error to use when a client should not retry until
// the system state has been explicitly fixed.  E.g., if an "rmdir"
// fails because the directory is non-empty, PreconditionFailed
// should be returned since the client should not retry unless
// the files are deleted from the directory. PreconditionFailed should also be
// returned when an intermediate directory for an MKCOL or PUT is missing.
//
// # FIXME rename to FailedPrecondition to make it less confusable with the http status Precondition Failed
//
// HTTP Mapping: 400 Bad Request, 405 Method Not Allowed, 409 Conflict
type PreconditionFailed string

func (e PreconditionFailed) Error() string { return "error: precondition failed: " + string(e) }

// IsPreconditionFailed implements the IsPreconditionFailed interface.
func (e PreconditionFailed) IsPreconditionFailed() {}

// AlreadyExists is the error to use when a resource something is not found.
type AlreadyExists string

func (e AlreadyExists) Error() string { return "error: already exists: " + string(e) }

// IsAlreadyExists implements the IsAlreadyExists interface.
func (e AlreadyExists) IsAlreadyExists() {}

// UserRequired represents an error when a resource is not found.
type UserRequired string

func (e UserRequired) Error() string { return "error: user required: " + string(e) }

// IsUserRequired implements the IsUserRequired interface.
func (e UserRequired) IsUserRequired() {}

// InvalidCredentials is the error to use when receiving invalid credentials.
type InvalidCredentials string

func (e InvalidCredentials) Error() string { return "error: invalid credentials: " + string(e) }

// IsInvalidCredentials implements the IsInvalidCredentials interface.
func (e InvalidCredentials) IsInvalidCredentials() {}

// NotSupported is the error to use when an action is not supported.
type NotSupported string

func (e NotSupported) Error() string { return "error: not supported: " + string(e) }

// IsNotSupported implements the IsNotSupported interface.
func (e NotSupported) IsNotSupported() {}

// PartialContent is the error to use when the client request has partial data.
type PartialContent string

func (e PartialContent) Error() string { return "error: partial content: " + string(e) }

// IsPartialContent implements the IsPartialContent interface.
func (e PartialContent) IsPartialContent() {}

// BadRequest is the error to use when the server cannot or will not process the request (due to a client error). Reauthenticating won't help.
type BadRequest string

func (e BadRequest) Error() string { return "error: bad request: " + string(e) }

// IsBadRequest implements the IsBadRequest interface.
func (e BadRequest) IsBadRequest() {}

// ChecksumMismatch is the error to use when the sent hash does not match the calculated hash.
type ChecksumMismatch string

func (e ChecksumMismatch) Error() string { return "error: checksum mismatch: " + string(e) }

// IsChecksumMismatch implements the IsChecksumMismatch interface.
func (e ChecksumMismatch) IsChecksumMismatch() {}

// StatusChecksumMismatch 419 is an unofficial http status code in an unassigned range that is used for checksum mismatches
// Proposed by https://stackoverflow.com/a/35665694
// Official HTTP status code registry: https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml
// Note: TUS uses unassigned 460 Checksum-Mismatch
// RFC proposal for checksum digest uses a `Want-Digest` header: https://tools.ietf.org/html/rfc3230
// oc clienst issue: https://github.com/owncloud/core/issues/22711
const StatusChecksumMismatch = 419

// InsufficientStorage is the error to use when there is insufficient storage.
type InsufficientStorage string

func (e InsufficientStorage) Error() string { return "error: insufficient storage: " + string(e) }

// IsInsufficientStorage implements the IsInsufficientStorage interface.
func (e InsufficientStorage) IsInsufficientStorage() {}

// StatusCode returns StatusInsufficientStorage, this implementation is needed to allow TUS to cast the correct http errors.
func (e InsufficientStorage) StatusCode() int {
	return StatusInsufficientStorage
}

// Body returns the error body. This implementation is needed to allow TUS to cast the correct http errors
func (e InsufficientStorage) Body() []byte {
	return []byte(e.Error())
}

// StatusInsufficientStorage 507 is an official HTTP status code to indicate that there is insufficient storage
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/507
const StatusInsufficientStorage = 507

// IsNotFound is the interface to implement
// to specify that an a resource is not found.
type IsNotFound interface {
	IsNotFound()
}

// IsAlreadyExists is the interface to implement
// to specify that a resource already exists.
type IsAlreadyExists interface {
	IsAlreadyExists()
}

// IsInternalError is the interface to implement
// to specify that there was some internal error
type IsInternalError interface {
	IsInternalError()
}

// IsUserRequired is the interface to implement
// to specify that a user is required.
type IsUserRequired interface {
	IsUserRequired()
}

// IsInvalidCredentials is the interface to implement
// to specify that credentials were wrong.
type IsInvalidCredentials interface {
	IsInvalidCredentials()
}

// IsNotSupported is the interface to implement
// to specify that an action is not supported.
type IsNotSupported interface {
	IsNotSupported()
}

// IsPermissionDenied is the interface to implement
// to specify that an action is denied.
type IsPermissionDenied interface {
	IsPermissionDenied()
}

// IsLocked is the interface to implement
// to specify that an resource is locked.
type IsLocked interface {
	IsLocked()
}

// IsAborted is the interface to implement
// to specify that a request was aborted.
type IsAborted interface {
	IsAborted()
}

// IsPreconditionFailed is the interface to implement
// to specify that a precondition failed.
type IsPreconditionFailed interface {
	IsPreconditionFailed()
}

// IsPartialContent is the interface to implement
// to specify that the client request has partial data.
type IsPartialContent interface {
	IsPartialContent()
}

// IsBadRequest is the interface to implement
// to specify that the server cannot or will not process the request.
type IsBadRequest interface {
	IsBadRequest()
}

// IsChecksumMismatch is the interface to implement
// to specify that a checksum does not match.
type IsChecksumMismatch interface {
	IsChecksumMismatch()
}

// IsInsufficientStorage is the interface to implement
// to specify that there is insufficient storage.
type IsInsufficientStorage interface {
	IsInsufficientStorage()
}

// NewErrtypeFromStatus maps an rpc status to an errtype
func NewErrtypeFromStatus(status *rpc.Status) error {
	switch status.Code {
	case rpc.Code_CODE_OK:
		return nil
	case rpc.Code_CODE_NOT_FOUND:
		return NotFound(status.Message)
	case rpc.Code_CODE_ALREADY_EXISTS:
		return AlreadyExists(status.Message)
		// case rpc.Code_CODE_FAILED_PRECONDITION: ?
		// return UserRequired(status.Message)
		// case rpc.Code_CODE_PERMISSION_DENIED: ?
		// IsInvalidCredentials
	case rpc.Code_CODE_UNIMPLEMENTED:
		return NotSupported(status.Message)
	case rpc.Code_CODE_PERMISSION_DENIED:
		// FIXME add locked status!
		if strings.HasPrefix(status.Message, "set lock: error: locked by ") {
			return Locked(strings.TrimPrefix(status.Message, "set lock: error: locked by "))
		}
		return PermissionDenied(status.Message)
	case rpc.Code_CODE_LOCKED:
		return Locked(status.Message)
	// case rpc.Code_CODE_DATA_LOSS: ?
	//	IsPartialContent
	case rpc.Code_CODE_ABORTED:
		return Aborted(status.Message)
	case rpc.Code_CODE_FAILED_PRECONDITION:
		return PreconditionFailed(status.Message)
	case rpc.Code_CODE_INSUFFICIENT_STORAGE:
		return InsufficientStorage(status.Message)
	case rpc.Code_CODE_INVALID_ARGUMENT, rpc.Code_CODE_OUT_OF_RANGE:
		return BadRequest(status.Message)
	default:
		return InternalError(status.Message)
	}
}
