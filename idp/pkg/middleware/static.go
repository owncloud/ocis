package middleware

import (
	"net/http"
	"strings"

	idpTracing "github.com/owncloud/ocis/idp/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Static is a middleware that serves static assets.
func Static(root string, fs http.FileSystem) func(http.Handler) http.Handler {
	if !strings.HasSuffix(root, "/") {
		root = root + "/"
	}

	static := http.StripPrefix(
		root,
		http.FileServer(
			fs,
		),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, span := idpTracing.TraceProvider.Tracer("idp").Start(r.Context(), "serve static asset")
			defer span.End()
			r = r.WithContext(ctx)

			// serve the static assets for the identifier web app
			if strings.HasPrefix(r.URL.Path, "/signin/v1/static/") {
				if strings.HasSuffix(r.URL.Path, "/") {
					// but no listing of folders
					span.SetAttributes(attribute.KeyValue{Key: "path", Value: attribute.StringValue(r.URL.Path)})
					span.SetStatus(codes.Error, "asset not found")
					http.NotFound(w, r)
				} else {
					r.URL.Path = strings.Replace(r.URL.Path, "/signin/v1/static/", "/signin/v1/identifier/static/", 1)
					span.SetAttributes(attribute.KeyValue{Key: "server", Value: attribute.StringValue(r.URL.Path)})
					static.ServeHTTP(w, r)
				}
				return
			}
			span.SetAttributes(attribute.KeyValue{Key: "path", Value: attribute.StringValue(r.URL.Path)})
			span.SetStatus(codes.Ok, "ok")
		})
	}
}
