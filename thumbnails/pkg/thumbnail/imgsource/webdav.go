package imgsource

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/pkg/errors"
	"image"
	_ "image/gif"  // Import the gif package so that image.Decode can understand gifs
	_ "image/jpeg" // Import the jpeg package so that image.Decode can understand jpegs
	_ "image/png"  // Import the png package so that image.Decode can understand pngs
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
func (s WebDav) Get(ctx context.Context, url string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, `could not get the image "%s"`, url)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: s.insecure} //nolint:gosec

	if auth, ok := ContextGetAuthorization(ctx); ok {
		req.Header.Add("Authorization", auth)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, `could not get the image "%s"`, url)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get the image \"%s\". Request returned with statuscode %d ", url, resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, `could not decode the image "%s"`, url)
	}
	return img, nil
}
