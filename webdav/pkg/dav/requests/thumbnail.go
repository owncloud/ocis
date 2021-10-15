package requests

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

const (
	// DefaultWidth defines the default width of a thumbnail
	DefaultWidth = 32
	// DefaultHeight defines the default height of a thumbnail
	DefaultHeight = 32
)

// ThumbnailRequest combines all parameters provided when requesting a thumbnail
type ThumbnailRequest struct {
	// The file path of the source file
	Filepath        string
	// The file name of the source file including the extension
	Filename        string
	// The file extension
	Extension       string
	// The requested width of the thumbnail
	Width           int32
	// The requested height of the thumbnail
	Height          int32
	// In case of a public share the public link token.
	PublicLinkToken string
}

// ParseThumbnailRequest extracts all required parameters from a http request.
func ParseThumbnailRequest(r *http.Request) (*ThumbnailRequest, error) {
	fp, err := extractFilePath(r)
	if err != nil {
		return nil, err
	}
	q := r.URL.Query()

	width, height, err := parseDimensions(q)
	if err != nil {
		return nil, err
	}

	return &ThumbnailRequest{
		Filepath:        fp,
		Filename:        filepath.Base(fp),
		Extension:       filepath.Ext(fp),
		Width:           int32(width),
		Height:          int32(height),
		PublicLinkToken: chi.URLParam(r, "token"),
	}, nil
}

// the url looks as followed
//
// /remote.php/dav/files/<user>/<filepath>
//
// User and filepath are dynamic and filepath can contain slashes
// So using the URLParam function is not possible.
func extractFilePath(r *http.Request) (string, error) {
	user := chi.URLParam(r, "user")
	if user != "" {
		parts := strings.SplitN(r.URL.Path, user, 2)
		return parts[1], nil
	}
	token := chi.URLParam(r, "token")
	if token != "" {
		parts := strings.SplitN(r.URL.Path, token, 2)
		return parts[1], nil
	}
	return "", errors.New("could not extract file path")
}

func parseDimensions(q url.Values) (int64, int64, error) {
	width, err := parseDimension(q.Get("x"), "width", DefaultWidth)
	if err != nil {
		return 0, 0, err
	}
	height, err := parseDimension(q.Get("y"), "height", DefaultHeight)
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func parseDimension(d, name string, defaultValue int64) (int64, error) {
	if d == "" {
		return defaultValue, nil
	}
	result, err := strconv.ParseInt(d, 10, 32)
	if err != nil || result < 1 {
		// The error message doesn't fit but for OC10 API compatibility reasons we have to set this.
		return 0, fmt.Errorf("Cannot set %s of 0 or smaller!", name) //nolint:golint
	}
	return result, nil
}
