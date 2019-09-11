package static

import (
	"net/http"

	"github.com/owncloud/ocis-phoenix/pkg/assets"
)

// static gets initialized by New and provides the handler.
type static struct {
	root string
	path string

	handler http.Handler
}

// ServeHTTP just implements the http.Handler interface.
func (s static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.handler.ServeHTTP(w, r)
}

// Handler returns the handler for static endpoint.
func Handler(opts ...Option) http.Handler {
	s := new(static)

	for _, opt := range opts {
		opt(s)
	}

	s.handler = http.StripPrefix(
		s.root,
		http.FileServer(
			assets.New(
				assets.WithPath(s.path),
			),
		),
	)

	return s
}
