package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
)

// OCSFormatCtx middleware is used to determine the content type from
// the format URL parameter passed in an ocs request. Defaults to XML
func OCSFormatCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("format") {
		case "", "xml":
			r.Header.Set("Accept", "application/xml")
			r = r.WithContext(context.WithValue(r.Context(), render.ContentTypeCtxKey, render.ContentTypeXML))
		case "json":
			r.Header.Set("Accept", "application/json")
			r = r.WithContext(context.WithValue(r.Context(), render.ContentTypeCtxKey, render.ContentTypeJSON))
		}
		next.ServeHTTP(w, r)
	})
}
