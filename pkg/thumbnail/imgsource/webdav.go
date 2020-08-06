package imgsource

import (
	"context"
	"crypto/tls"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"path"

	"github.com/owncloud/ocis-thumbnails/pkg/config"
)

// NewWebDavSource creates a new webdav instance.
func NewWebDavSource(cfg config.WebDavSource) WebDav {
	return WebDav{
		baseURL: cfg.BaseURL,
	}
}

// WebDav implements the Source interface for webdav services
type WebDav struct {
	baseURL string
}

// Get downloads the file from a webdav service
func (s WebDav) Get(ctx context.Context, file string) (image.Image, error) {
	u, _ := url.Parse(s.baseURL)
	u.Path = path.Join(u.Path, file)
	fmt.Printf("url: %s", u.String())
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not get the image \"%s\" error: %s", file, err.Error())
	}

	// FIXME: make this configurable!!
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	auth := authorization(ctx)
	fmt.Printf("auth: %s", auth)
	if auth == "" {
		return nil, fmt.Errorf("could not get image \"%s\" error: authorization is missing", file)
	}
	req.Header.Add("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get the image \"%s\" error: %s", file, err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get the image \"%s\". Request returned with statuscode %d ", file, resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not decode the image \"%s\". error: %s", file, err.Error())
	}
	return img, nil
}
