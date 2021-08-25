package response

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
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

// APIVersion retrieves the api version from the context.
func APIVersion(ctx context.Context) string {
	value := ctx.Value(apiVersionKey)
	if value != nil {
		return value.(string)
	}
	return ""
}

// OcsV1StatusCodes returns the http status codes for the OCS API v1.
func OcsV1StatusCodes(meta data.Meta) int {
	if meta.StatusCode == data.MetaUnauthorized.StatusCode {
		return http.StatusUnauthorized
	}
	return http.StatusOK
}

// OcsV2StatusCodes maps the OCS codes to http status codes for the ocs API v2.
// see https://github.com/owncloud/core/blob/c08baf580927ecb8ec179028dda255fdd85b4568/lib/private/legacy/api.php#L528
// also HTTP status codes for apps are the same as OCS codes
// see https://github.com/owncloud/core/blob/b9ff4c93e051c94adfb301545098ae627e52ef76/lib/public/AppFramework/OCSController.php#L142-L150
// I think this is a bug in the ocs v2 api, but since we are going to mimic bugs in ocis ... here goes
func OcsV2StatusCodes(meta data.Meta) int {
	sc := meta.StatusCode
	switch sc {
	case data.MetaNotFound.StatusCode:
		return http.StatusNotFound
	case data.MetaUnknownError.StatusCode:
		fallthrough
	case data.MetaServerError.StatusCode:
		return http.StatusInternalServerError
	case data.MetaUnauthorized.StatusCode:
		return http.StatusUnauthorized
	case data.MetaOK.StatusCode:
		// TODO mustn't data.Meta be a pointer so this assignment has an effect
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

	// TODO change this status code? yes, align with oc10 core mapStatusCodes
	return http.StatusOK
}

// VersionCtx middleware is used to determine the response mapper from
// the URL parameters passed through as the request. In case
// the Version is unknown, we stop here and return a 404.
func VersionCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := chi.URLParam(r, "version")
		if version == "" {
			mustNotFail(render.Render(w, r, ErrRender(data.MetaBadRequest.StatusCode, "unknown ocs api version")))
			return
		}
		w.Header().Set("Ocs-Api-Version", version)

		// store version in context so handlers can access it
		ctx := context.WithValue(r.Context(), apiVersionKey, version)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func mustNotFail(err error) {
	if err != nil {
		panic(err)
	}
}
