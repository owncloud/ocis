package middleware

import (
	"net/http"

	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

// Trace unpacks the request context looking for an existing trace id.
func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var span *trace.Span

		tc := tracecontext.HTTPFormat{}
		sc, ok := tc.SpanContextFromRequest(r)
		if ok {
			ctx, span = trace.StartSpanWithRemoteParent(r.Context(), r.URL.String(), sc)
			defer span.End()
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
