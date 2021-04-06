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
func (s WebDav) Get(ctx context.Context, file string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, file, nil)
	if err != nil {
		return nil, errors.Wrapf(err, `could not get the image "%s"`, file)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: s.insecure} //nolint:gosec

	auth, ok := ContextGetAuthorization(ctx)
	if !ok {
		return nil, fmt.Errorf("could not get image \"%s\" error: authorization is missing", file)
	}
	req.Header.Add("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, `could not get the image "%s"`, file)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get the image \"%s\". Request returned with statuscode %d ", file, resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, `could not decode the image "%s"`, file)
	}
	return img, nil
}
