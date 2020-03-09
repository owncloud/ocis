package imgsource

import (
	"fmt"
	"image"
	"net/http"
	"net/url"
	"path"
)

// WebDav implements the Source interface for webdav services
type WebDav struct {
	Basepath string
}

const (
	// WebDavAuth is the parameter name for the autorization token
	WebDavAuth = "Authorization"
)

// Get downloads the file from a webdav service
func (s WebDav) Get(file string, ctx SourceContext) (image.Image, error) {
	u, _ := url.Parse(s.Basepath)
	u.Path = path.Join(u.Path, file)
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not get the file \"%s\" error: %s", file, err.Error())
	}

	auth := ctx.GetString(WebDavAuth)
	req.Header.Add("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not get the file \"%s\" error: %s", file, err.Error())
	}

	img, _, _ := image.Decode(resp.Body)
	return img, nil
}
