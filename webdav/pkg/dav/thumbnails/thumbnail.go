package thumbnail

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

const (
	// DefaultWidth defines the default width of a thumbnail
	DefaultWidth = 32
	// DefaultHeight defines the default height of a thumbnail
	DefaultHeight = 32
)

// Request combines all parameters provided when requesting a thumbnail
type Request struct {
	Filepath      string
	Filetype      string
	Etag          string
	Width         int
	Height        int
	Authorization string
	Username      string
}

// NewRequest extracts all required parameters from a http request.
func NewRequest(r *http.Request) (Request, error) {
	path := extractFilePath(r)
	query := r.URL.Query()
	width, err := strconv.Atoi(query.Get("x"))
	if err != nil {
		width = DefaultWidth
	}
	height, err := strconv.Atoi(query.Get("y"))
	if err != nil {
		height = DefaultHeight
	}

	etag := query.Get("c")
	if strings.TrimSpace(etag) == "" {
		return Request{}, fmt.Errorf("c (etag) is missing in query")
	}

	authorization := r.Header.Get("Authorization")

	tr := Request{
		Filepath:      path,
		Filetype:      strings.Replace(filepath.Ext(path), ".", "", 1),
		Etag:          etag,
		Width:         width,
		Height:        height,
		Authorization: authorization,
		Username:      chi.URLParam(r, "user"),
	}

	return tr, nil
}

// the url looks as followed
//
// /remote.php/dav/files/<user>/<filepath>
//
// User and filepath are dynamic and filepath can contain slashes
// So using the URLParam function is not possible.
func extractFilePath(r *http.Request) string {
	user := chi.URLParam(r, "user")
	parts := strings.SplitN(r.URL.Path, user, 2)
	return parts[1]
}
