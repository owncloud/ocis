package autoprop

import (
	"context"
	"net/http"
)

// NewHttpHandler creates a new HTTP server handler to inject the custom
// HTTP autoprogation headers in the request context.
func NewHttpHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := moveHttpHeadersToOcisMeta(r, r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// autoPropRoundTripper will inject the autopropagation data in the context
// inside the HTTP headers. A base RoundTripper is needed, so use
// NewHttpRoundTripper with a base RoundTripper to create an instance
type autoPropRoundTripper struct {
	base http.RoundTripper
}

// RoundTrip implements the RoundTripper interface. This method will inject
// the autopropagation data from the context into the request headers before
// sending the modified request to the base RoundTripper
func (rt *autoPropRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	r2 := r.Clone(r.Context())
	moveOcisMetaToHttpHeaders(r2, r.Context())
	return rt.base.RoundTrip(r2)
}

// NewHttpRoundTripper creates a new instance of the autoPropRoundTripper.
// This is used by the HTTP clients to inject the autopropagation headers
// in the HTTP requests.
func NewHttpRoundTripper(base http.RoundTripper) http.RoundTripper {
	return &autoPropRoundTripper{
		base: base,
	}
}

// AppendToHttpRequest adds the autopropagation data as HTTP headers in the
// provided request.
func AppendToHttpRequest(r *http.Request, ctx context.Context) {
	moveOcisMetaToHttpHeaders(r, ctx)
}
