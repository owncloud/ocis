package commands

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrStatusNotEmpty may be returned if a call should not have a status
	// string set but one is.
	ErrStatusNotEmpty = errors.New("response status not empty")
	// ErrBodyNotEmpty may be returned if a call should have an empty body but
	// a body value is present.
	ErrBodyNotEmpty = errors.New("response body not empty")
)

const (
	deprecatedSuffix = "call is deprecated and will be removed in a future release"
	missingPrefix    = "No handler found"
	einval           = -22
)

type cephError interface {
	ErrorCode() int
}

// NotImplementedError error values will be returned in the case that an API
// call is not available in the version of Ceph that is running in the target
// cluster.
type NotImplementedError struct {
	Response
}

// Error implements the error interface.
func (e NotImplementedError) Error() string {
	return fmt.Sprintf("API call not implemented server-side: %s", e.status)
}

// Response encapsulates the data returned by ceph and supports easy processing
// pipelines.
type Response struct {
	body   []byte
	status string
	err    error
}

// Ok returns true if the response contains no error.
func (r Response) Ok() bool {
	return r.err == nil
}

// Error implements the error interface.
func (r Response) Error() string {
	if r.status == "" {
		return r.err.Error()
	}
	return fmt.Sprintf("%s: %q", r.err, r.status)
}

// Unwrap returns the error this response contains.
func (r Response) Unwrap() error {
	return r.err
}

// Status returns the status string value.
func (r Response) Status() string {
	return r.status
}

// Body returns the response body as a raw byte-slice.
func (r Response) Body() []byte {
	return r.body
}

// End returns an error if the response contains an error or nil, indicating
// that response is no longer needed for processing.
func (r Response) End() error {
	if !r.Ok() {
		if ce, ok := r.err.(cephError); ok {
			if ce.ErrorCode() == einval && strings.HasPrefix(r.status, missingPrefix) {
				return NotImplementedError{Response: r}
			}
		}
		return r
	}
	return nil
}

// NoStatus asserts that the input response has no status value.
func (r Response) NoStatus() Response {
	if !r.Ok() {
		return r
	}
	if r.status != "" {
		return Response{r.body, r.status, ErrStatusNotEmpty}
	}
	return r
}

// NoBody asserts that the input response has no body value.
func (r Response) NoBody() Response {
	if !r.Ok() {
		return r
	}
	if len(r.body) != 0 {
		return Response{r.body, r.status, ErrBodyNotEmpty}
	}
	return r
}

// EmptyBody is similar to NoBody but also accepts an empty JSON object.
func (r Response) EmptyBody() Response {
	if !r.Ok() {
		return r
	}
	if len(r.body) != 0 {
		d := map[string]interface{}{}
		if err := json.Unmarshal(r.body, &d); err != nil {
			return Response{r.body, r.status, err}
		}
		if len(d) != 0 {
			return Response{r.body, r.status, ErrBodyNotEmpty}
		}
	}
	return r
}

// NoData asserts that the input response has no status or body values.
func (r Response) NoData() Response {
	return r.NoStatus().NoBody()
}

// FilterPrefix sets the status value to an empty string if the status
// value contains the given prefix string.
func (r Response) FilterPrefix(p string) Response {
	if !r.Ok() {
		return r
	}
	if strings.HasPrefix(r.status, p) {
		return Response{r.body, "", r.err}
	}
	return r
}

// FilterSuffix sets the status value to an empty string if the status
// value contains the given suffix string.
func (r Response) FilterSuffix(s string) Response {
	if !r.Ok() {
		return r
	}
	if strings.HasSuffix(r.status, s) {
		return Response{r.body, "", r.err}
	}
	return r
}

// FilterBodyPrefix sets the body value equivalent to an empty string if the
// body value contains the given prefix string.
func (r Response) FilterBodyPrefix(p string) Response {
	if !r.Ok() {
		return r
	}
	if bytes.HasPrefix(r.body, []byte(p)) {
		return Response{[]byte(""), r.status, r.err}
	}
	return r
}

// FilterDeprecated removes deprecation warnings from the response status.
// Use it when checking the response from calls that may be deprecated in ceph
// if you want those calls to continue working if the warning is present.
func (r Response) FilterDeprecated() Response {
	return r.FilterSuffix(deprecatedSuffix)
}

// Unmarshal data from the response body into v.
func (r Response) Unmarshal(v interface{}) Response {
	if !r.Ok() {
		return r
	}
	if err := json.Unmarshal(r.body, v); err != nil {
		return Response{body: r.body, err: err}
	}
	return r
}

// NewResponse returns a response.
func NewResponse(b []byte, s string, e error) Response {
	return Response{b, s, e}
}
