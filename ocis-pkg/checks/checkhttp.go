package checks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
)

// NewHTTPCheck checks the reachability of a http server.
func NewHTTPCheck(rawUrl string) func(context.Context) error {
	return func(_ context.Context) error {
		if !strings.HasPrefix(rawUrl, "http://") && !strings.HasPrefix(rawUrl, "https://") {
			rawUrl = "http://" + rawUrl
		}

		parsedUrl, err := url.Parse(rawUrl)
		if err != nil {
			return err
		}

		parsedUrl.Host, err = handlers.FailSaveAddress(parsedUrl.Host)
		if err != nil {
			return err
		}

		c := http.Client{
			Timeout: 3 * time.Second,
		}
		resp, err := c.Get(parsedUrl.String())
		if err != nil {
			return fmt.Errorf("could not connect to http server: %v", err)
		}
		_ = resp.Body.Close()
		return nil
	}
}
