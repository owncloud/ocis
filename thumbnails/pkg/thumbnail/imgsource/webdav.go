package imgsource

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/pkg/errors"
	_ "image/gif"  // Import the gif package so that image.Decode can understand gifs
	_ "image/jpeg" // Import the jpeg package so that image.Decode can understand jpegs
	_ "image/png"  // Import the png package so that image.Decode can understand pngs
	"io"
	"net/http"
)

// NewWebDavSource creates a new webdav instance.
func NewWebDavSource(cfg config.Thumbnail) WebDav {
	return WebDav{
		insecure: cfg.WebdavAllowInsecure,
	}
}

// WebDav implements the Source interface for webdav services
type WebDav struct {
	insecure bool
}

// Get downloads the file from a webdav service
// The caller MUST make sure to close the returned ReadCloser
func (s WebDav) Get(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, `could not get the image "%s"`, url)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: s.insecure} //nolint:gosec

	if auth, ok := ContextGetAuthorization(ctx); ok {
		req.Header.Add("Authorization", auth)
	}

	client := &http.Client{}
	resp, err := client.Do(req) // nolint:bodyclose
	if err != nil {
		return nil, errors.Wrapf(err, `could not get the image "%s"`, url)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get the image \"%s\". Request returned with statuscode %d ", url, resp.StatusCode)
	}

	return resp.Body, nil
}
