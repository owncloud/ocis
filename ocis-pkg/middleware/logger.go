package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/owncloud/ocis/ocis-pkg/log"
)

// Logger is a middleware to log http requests. It uses debug level logging and should be used by all services save the proxy (which uses info level logging).
func Logger(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrap := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(wrap, r)

			logger.Debug().
				Str("request", r.Header.Get("X-Request-ID")).
				Str("proto", r.Proto).
				Str("method", r.Method).
				Int("status", wrap.Status()).
				Str("path", r.URL.Path).
				Dur("duration", time.Since(start)).
				Int("bytes", wrap.BytesWritten()).
				Msg("")
		})
	}
}
