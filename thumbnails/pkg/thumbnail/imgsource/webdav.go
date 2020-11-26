package imgsource

import (
	"context"
	"crypto/tls"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"path"

	"github.com/owncloud/ocis/thumbnails/pkg/config"
	"github.com/pkg/errors"
)

// NewWebDavSource creates a new webdav instance.
func NewWebDavSource(cfg config.WebDavSource) WebDav {
	return WebDav{
		baseURL:  cfg.BaseURL,
		insecure: cfg.Insecure,
	}
}

// WebDav implements the Source interface for webdav services
type WebDav struct {
	baseURL  string
	insecure bool
}

// Get downloads the file from a webdav service
func (s WebDav) Get(ctx context.Context, file string) (image.Image, error) {
	u, _ := url.Parse(s.baseURL)
	u.Path = path.Join(u.Path, file)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrapf(err, `could not get the image "%s"`, file)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: s.insecure}

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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get the image \"%s\". Request returned with statuscode %d ", file, resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, `could not decode the image "%s"`, file)
	}
	return img, nil
}
