package checks

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// NewHTTPCheck checks the reachability of a http server.
func NewHTTPCheck(url string) func(context.Context) error {
	return func(_ context.Context) error {
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
