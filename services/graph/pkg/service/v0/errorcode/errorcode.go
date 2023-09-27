// Package errorcode allows to deal with graph error codes
package errorcode

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

// ErrorCode defines code as used in MS Graph - see https://docs.microsoft.com/en-us/graph/errors?context=graph%2Fapi%2F1.0&view=graph-rest-1.0
type ErrorCode int

// Error defines a custom error struct, containing and MS Graph error code an a textual error message
type Error struct {
	errorCode ErrorCode
	msg       string
}

const (
	// AccessDenied defines the error if the caller doesn't have permission to perform the action.
	AccessDenied ErrorCode = iota
	// ActivityLimitReached defines the error if the app or user has been throttled.
	ActivityLimitReached
	// GeneralException defines the error if an unspecified error has occurred.
	GeneralException
	// InvalidAuthenticationToken defines the error if the access token is missing
	InvalidAuthenticationToken
	// InvalidRange defines the error if the specified byte range is invalid or unavailable.
	InvalidRange
	// InvalidRequest defines the error if the request is malformed or incorrect.
	InvalidRequest
	// ItemNotFound defines the error if the resource could not be found.
	ItemNotFound
	// MalwareDetected defines the error if malware was detected in the requested resource.
	MalwareDetected
	// NameAlreadyExists defines the error if the specified item name already exists.
	NameAlreadyExists
	// NotAllowed defines the error if the action is not allowed by the system.
	NotAllowed
	// NotSupported defines the error if the request is not supported by the system.
	NotSupported
	// ResourceModified defines the error if the resource being updated has changed since the caller last read it, usually an eTag mismatch.
	ResourceModified
	// ResyncRequired defines the error if the delta token is no longer valid, and the app must reset the sync state.
	ResyncRequired
	// ServiceNotAvailable defines the error if the service is not available. Try the request again after a delay. There may be a Retry-After header.
	ServiceNotAvailable
	// QuotaLimitReached the user has reached their quota limit.
	QuotaLimitReached
	// Unauthenticated the caller is not authenticated.
	Unauthenticated
	// PreconditionFailed the request cannot be made and this error response is sent back
	PreconditionFailed
)

var errorCodes = [...]string{
	"accessDenied",
	"activityLimitReached",
	"generalException",
	"InvalidAuthenticationToken",
	"invalidRange",
	"invalidRequest",
	"itemNotFound",
	"malwareDetected",
	"nameAlreadyExists",
	"notAllowed",
	"notSupported",
	"resourceModified",
	"resyncRequired",
	"serviceNotAvailable",
	"quotaLimitReached",
	"unauthenticated",
	"preconditionFailed",
}

// New constructs a new errorcode.Error
func New(e ErrorCode, msg string) Error {
	return Error{
		errorCode: e,
		msg:       msg,
	}
}

// Render writes an Graph ErrorCode	object to the response writer
func (e ErrorCode) Render(w http.ResponseWriter, r *http.Request, status int, msg string) {
	innererror := map[string]interface{}{
		"date": time.Now().UTC().Format(time.RFC3339),
	}

	innererror["request-id"] = middleware.GetReqID(r.Context())
	resp := &libregraph.OdataError{
		Error: libregraph.OdataErrorMain{
			Code:       e.String(),
			Message:    msg,
			Innererror: innererror,
		},
	}
	render.Status(r, status)
	render.JSON(w, r, resp)
}

// Render writes an Graph Error object to the response writer
func (e Error) Render(w http.ResponseWriter, r *http.Request) {
	var status int
	switch e.errorCode {
	case ItemNotFound:
		status = http.StatusNotFound
	case NameAlreadyExists:
		status = http.StatusConflict
	case NotAllowed:
		status = http.StatusForbidden
	default:
		status = http.StatusInternalServerError
	}
	e.errorCode.Render(w, r, status, e.msg)
}

// String returns the string corresponding to the ErrorCode
func (e ErrorCode) String() string {
	return errorCodes[e]
}

// Error return the concatenation of the error string and optinal message
func (e Error) Error() string {
	errString := errorCodes[e.errorCode]
	if e.msg != "" {
		errString += ": " + e.msg
	}
	return errString
}

// RenderError render the Graph Error based on a code or default one
func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	var errcode Error
	if errors.As(err, &errcode) {
		errcode.Render(w, r)
	} else {
		GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
	}
}
