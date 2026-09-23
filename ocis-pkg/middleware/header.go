package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/cors"

	rscors "github.com/rs/cors"
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

	allowCredentials := options.AllowCredentials
	// A wildcard origin must not be combined with credentials.
	if allowCredentials && cors.AllowsAnyOrigin(options.AllowedOrigins) {
		logger.Warn().
			Strs("allowed_origins", options.AllowedOrigins).
			Msg("cors: refusing to allow credentials together with a wildcard origin, disabling allow_credentials")
		allowCredentials = false
	}

	logger.Debug().
		Str("allowed_origins", strings.Join(options.AllowedOrigins, ", ")).
		Str("allowed_methods", strings.Join(options.AllowedMethods, ", ")).
		Str("allowed_headers", strings.Join(options.AllowedHeaders, ", ")).
		Bool("allow_credentials", allowCredentials).
		Msg("setup cors middleware")
	c := rscors.New(rscors.Options{
		AllowedOrigins:   options.AllowedOrigins,
		AllowedMethods:   options.AllowedMethods,
		AllowedHeaders:   options.AllowedHeaders,
		AllowCredentials: allowCredentials,
	})
	return c.Handler
}
