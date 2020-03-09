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
		// reconstruct span context from request
		if sc, ok := tc.SpanContextFromRequest(r); ok {
			// if there is one, add it to the new span
			ctx, span = trace.StartSpanWithRemoteParent(r.Context(), r.URL.String(), sc)
			defer span.End()
		} else {
			// create a new span if there is no context
			ctx, span = trace.StartSpan(r.Context(), r.URL.String())
			defer span.End()
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
