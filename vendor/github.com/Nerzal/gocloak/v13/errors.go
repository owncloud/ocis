package gocloak

import (
	"strings"
)

// HTTPErrorResponse is a model of an error response
type HTTPErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Message     string `json:"errorMessage,omitempty"`
	Description string `json:"error_description,omitempty"`
}

// String returns a string representation of an error
func (e HTTPErrorResponse) String() string {
	var res strings.Builder
	if len(e.Error) > 0 {
		res.WriteString(e.Error)
	}
	if len(e.Message) > 0 {
		if res.Len() > 0 {
			res.WriteString(": ")
		}
		res.WriteString(e.Message)
	}
	if len(e.Description) > 0 {
		if res.Len() > 0 {
			res.WriteString(": ")
		}
		res.WriteString(e.Description)
	}
	return res.String()
}

// NotEmpty validates that error is not emptyp
func (e HTTPErrorResponse) NotEmpty() bool {
	return len(e.Error) > 0 || len(e.Message) > 0 || len(e.Description) > 0
}
