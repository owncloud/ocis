package middleware

import (
	"net/http"

	"github.com/ascarter/requestid"
)

// RequestID is a convenient middleware to inject a request id.
func RequestID(next http.Handler) http.Handler {
	return requestid.RequestIDHandler(next)
}
