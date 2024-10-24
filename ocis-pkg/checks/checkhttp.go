package checks

import (
	"context"
	"fmt"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"net/http"
	"strings"
	"time"
)

// NewHTTPCheck checks the reachability of a http server.
func NewHTTPCheck(url string) func(context.Context) error {
	return func(_ context.Context) error {
		if strings.Contains(url, "0.0.0.0") {
			outboundIp, err := handlers.GetOutBoundIP()
			if err != nil {
				return err
			}
			url = strings.Replace(url, "0.0.0.0", outboundIp, 1)
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
