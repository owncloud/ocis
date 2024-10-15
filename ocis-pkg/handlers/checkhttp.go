package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// NewHttpCheck checks the reachability of a http server.
func NewHTTPCheck(url string) func(context.Context) error {
	return func(_ context.Context) error {
		c := http.Client{
			Timeout: 3 * time.Second,
		}
		_, err := c.Get(url)
		if err != nil {
			return fmt.Errorf("could not connect to http server: %v", err)
		}
		return nil
	}
}
