package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// AccessLog is a middleware to log http requests at info level logging.
func AccessLog(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrap := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(wrap, r)

			logger.Info().
				Str("proto", r.Proto).
				Str("request", r.Header.Get("X-Request-ID")).
				Str("remote-addr", r.RemoteAddr).
				Str("method", r.Method).
				Int("status", wrap.Status()).
				Str("path", r.URL.Path).
				Dur("duration", time.Since(start)).
				Int("bytes", wrap.BytesWritten()).
				Msg("access-log")
		})
	}
}
