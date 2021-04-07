package requests

import (
	"errors"
	"net/http"
	"net/url"
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
type ThumbnailRequest struct {
	Filepath        string
	Extension       string
	Width           int32
	Height          int32
	PublicLinkToken string
}

// NewRequest extracts all required parameters from a http request.
func ParseThumbnailRequest(r *http.Request) (ThumbnailRequest, error) {
	fp := extractFilePath(r)
	q := r.URL.Query()

	width, height, err := parseDimensions(q)
	if err != nil {
		return ThumbnailRequest{}, err
	}

	tr := ThumbnailRequest{
		Filepath:        fp,
		Extension:       filepath.Ext(fp),
		Width:           int32(width),
		Height:          int32(height),
		PublicLinkToken: chi.URLParam(r, "token"),
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
	if user != "" {
		parts := strings.SplitN(r.URL.Path, user, 2)
		return parts[1]
	}
	token := chi.URLParam(r, "token")
	if token != "" {
		parts := strings.SplitN(r.URL.Path, token, 2)
		return parts[1]
	}
	return ""
}

func parseDimensions(q url.Values) (int64, int64, error) {
	width, err := parseDimension(q.Get("x"), DefaultWidth)
	if err != nil {
		return 0, 0, err
	}
	height, err := parseDimension(q.Get("y"), DefaultHeight)
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func parseDimension(d string, defaultValue int64) (int64, error) {
	if d == "" {
		return defaultValue, nil
	}
	result, err := strconv.ParseInt(d, 10, 32)
	if err != nil {
		return 0, err
	}
	if result < 1 {
		return 0, errors.New("invalid dimension")
	}

	return result, nil
}
