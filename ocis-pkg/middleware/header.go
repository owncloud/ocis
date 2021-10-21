package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/cors"

	chicors "github.com/go-chi/cors"
)

// NoCache writes required cache headers to all requests.
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		w.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

		next.ServeHTTP(w, r)
	})
}

// Cors writes required cors headers to all requests.
func Cors(opts ...cors.Option) func(http.Handler) http.Handler {
	options := cors.NewOptions(opts...)
	logger := options.Logger
	logger.Debug().
		Str("allowed_origins", strings.Join(options.AllowedOrigins, ", ")).
		Str("allowed_methods", strings.Join(options.AllowedMethods, ", ")).
		Str("allowed_headers", strings.Join(options.AllowedHeaders, ", ")).
		Bool("allow_credentials", options.AllowCredentials).
		Msg("setup cors middleware")
	return chicors.Handler(chicors.Options{
		AllowedOrigins:   options.AllowedOrigins,
		AllowedMethods:   options.AllowedMethods,
		AllowedHeaders:   options.AllowedHeaders,
		AllowCredentials: options.AllowCredentials,
	})
}

// Secure writes required access headers to all requests.
func Secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Indicates whether the browser is allowed to render this page in a <frame>, <iframe>, <embed> or <object>.
		w.Header().Set("X-Frame-Options", "DENY")
		// Does basically the same as X-Frame-Options.
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'")
		// This header inidicates that MIME types advertised in the Content-Type headers should not be changed and be followed.
		w.Header().Set("X-Content-Type-Options", "nosniff")

		if r.TLS != nil {
			// Tell browsers that the website should only be accessed  using HTTPS.
			w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		}

		next.ServeHTTP(w, r)
	})
}
