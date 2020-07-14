package svc

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type key int

const (
	apiVersionKey key = iota
	ocsVersion1       = "1"
	ocsVersion2       = "2"
)

var (
	defaultStatusCodeMapper = OcsV2StatusCodes
)

// Meta holds response metadata
type Meta struct {
	Status       string `json:"status" xml:"status"`
	StatusCode   int    `json:"statuscode" xml:"statuscode"`
	Message      string `json:"message" xml:"message"`
	TotalItems   string `json:"totalitems,omitempty" xml:"totalitems,omitempty"`
	ItemsPerPage string `json:"itemsperpage,omitempty" xml:"itemsperpage,omitempty"`
}

// MetaOK is the default ok response
var MetaOK = Meta{Status: "ok", StatusCode: 100, Message: "OK"}

// MetaBadRequest is used for unknown errors
var MetaBadRequest = Meta{Status: "error", StatusCode: 400, Message: "Bad Request"}

// MetaServerError is returned on server errors
var MetaServerError = Meta{Status: "error", StatusCode: 996, Message: "Server Error"}

// MetaUnauthorized is returned on unauthorized requests
var MetaUnauthorized = Meta{Status: "error", StatusCode: 997, Message: "Unauthorised"}

// MetaNotFound is returned when trying to access not existing resources
var MetaNotFound = Meta{Status: "error", StatusCode: 998, Message: "Not Found"}

// MetaUnknownError is used for unknown errors
var MetaUnknownError = Meta{Status: "error", StatusCode: 999, Message: "Unknown Error"}

// OcsV1StatusCodes returns the http status codes for the OCS API v1.
func OcsV1StatusCodes(meta Meta) int {
	return http.StatusOK
}

// OcsV2StatusCodes maps the OCS codes to http status codes for the ocs API v2.
func OcsV2StatusCodes(meta Meta) int {
	sc := meta.StatusCode
	switch sc {
	case MetaNotFound.StatusCode:
		return http.StatusNotFound
	case MetaUnknownError.StatusCode:
		fallthrough
	case MetaServerError.StatusCode:
		return http.StatusInternalServerError
	case MetaUnauthorized.StatusCode:
		return http.StatusUnauthorized
	case 100:
		meta.StatusCode = http.StatusOK
		return http.StatusOK
	}
	// any 2xx, 4xx and 5xx will be used as is
	if sc >= 200 && sc < 600 {
		return sc
	}

	// any error codes > 100 are treated as client errors
	if sc > 100 && sc < 200 {
		return http.StatusBadRequest
	}

	// TODO change this status code?
	return http.StatusOK
}

// APIVersion retrieves the api version from the context.
func APIVersion(ctx context.Context) string {
	value := ctx.Value(apiVersionKey)
	if value != nil {
		return value.(string)
	}
	return ""
}

// VersionCtx middleware is used to determine the response mapper from
// the URL parameters passed through as the request. In case
// the Version is unknown, we stop here and return a 404.
func (g Ocs) VersionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := chi.URLParam(r, "version")
		if version == "" {
			render.Render(w, r, ErrRender(MetaBadRequest.StatusCode, "unknown ocs api version"))
			return
		}
		w.Header().Set("Ocs-Api-Version", version)

		// store version in context so handlers can access it
		ctx := context.WithValue(r.Context(), apiVersionKey, version)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
