package middleware

import (
	"net/http"
	"strings"

	"go.opencensus.io/trace"
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
			ctx, span := trace.StartSpan(r.Context(), "serve static asset")
			defer span.End()
			r = r.WithContext(ctx)

			// serve the static assets for the identifier web app
			if strings.HasPrefix(r.URL.Path, "/signin/v1/static/") {
				if strings.HasSuffix(r.URL.Path, "/") {
					// but no listing of folders
					span.AddAttributes(trace.StringAttribute("asset not found", r.URL.Path))
					span.SetStatus(trace.Status{
						Code:    1,
						Message: "asset not found",
					})
					http.NotFound(w, r)
				} else {
					r.URL.Path = strings.Replace(r.URL.Path, "/signin/v1/static/", "/signin/v1/identifier/static/", 1)
					span.AddAttributes(trace.StringAttribute("served", r.URL.Path))
					static.ServeHTTP(w, r)
				}
				return
			}
			span.AddAttributes(trace.StringAttribute("served", r.URL.Path))
			span.SetStatus(trace.Status{
				Code:    0,
				Message: "ok",
			})
			next.ServeHTTP(w, r)
		})
	}
}
