package errorcode

import (
	"net/http"

	"github.com/go-chi/render"
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

// ErrorCode defines code as used in MS Graph - see https://docs.microsoft.com/en-us/graph/errors?context=graph%2Fapi%2F1.0&view=graph-rest-1.0
type ErrorCode int

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
}

// Render writes an Graph ErrorObject to the response writer
func (e ErrorCode) Render(w http.ResponseWriter, r *http.Request, status int, msg string) {
	resp := &msgraph.ErrorObject{
		Code:    e.String(),
		Message: msg,
	}
	render.Status(r, status)
	render.JSON(w, r, resp)
}

func (e ErrorCode) String() string {
	return errorCodes[e]
}
