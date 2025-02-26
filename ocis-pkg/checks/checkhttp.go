package checks

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
)

// NewHTTPCheck checks the reachability of a http server.
func NewHTTPCheck(url string) func(context.Context) error {
	return func(_ context.Context) error {
		url, err := handlers.FailSaveAddress(url)
		if err != nil {
			return err
		}

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}

		c := http.Client{
			Timeout: 3 * time.Second,
		}
		resp, err := c.Get(url)
		if err != nil {
			return fmt.Errorf("could not connect to http server: %v", err)
		}
		_ = resp.Body.Close()
		return nil
	}
}
