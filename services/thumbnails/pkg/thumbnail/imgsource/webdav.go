package imgsource

import (
	"context"
	"crypto/tls"
	"fmt"
	_ "image/gif"  // Import the gif package so that image.Decode can understand gifs
	_ "image/jpeg" // Import the jpeg package so that image.Decode can understand jpegs
	_ "image/png"  // Import the png package so that image.Decode can understand pngs
	"io"
	"net/http"
	"strconv"

	"github.com/cs3org/reva/v2/pkg/bytesize"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/config"
	thumbnailerErrors "github.com/owncloud/ocis/v2/services/thumbnails/pkg/errors"
	"github.com/pkg/errors"
)

// NewWebDavSource creates a new webdav instance.
func NewWebDavSource(cfg config.Thumbnail, b bytesize.ByteSize) WebDav {
	return WebDav{
		insecure:         cfg.WebdavAllowInsecure,
		maxImageFileSize: b.Bytes(),
	}
}

// WebDav implements the Source interface for webdav services
type WebDav struct {
	insecure         bool
	maxImageFileSize uint64
}

// Get downloads the file from a webdav service
// The caller MUST make sure to close the returned ReadCloser
func (s WebDav) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, `could not get the image "%s"`, url)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: s.insecure, //nolint:gosec
	}

	if auth, ok := ContextGetAuthorization(ctx); ok {
		req.Header.Add("Authorization", auth)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, `could not get the image "%s"`, url)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get the image \"%s\". Request returned with statuscode %d ", url, resp.StatusCode)
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		// no size information - let's assume it is too big
		return nil, thumbnailerErrors.ErrImageTooLarge
	}
	c, err := strconv.ParseUint(contentLength, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(err, `could not parse content length of webdav response "%s"`, url)
	}
	if c > s.maxImageFileSize {
		return nil, thumbnailerErrors.ErrImageTooLarge
	}

	return resp.Body, nil
}
