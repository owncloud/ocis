package handlers

import (
	"context"
	"fmt"
	"io"
	"maps"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// check is a function that performs a check.
type check func(ctx context.Context) error

// CheckHandlerConfiguration defines the configuration for the CheckHandler.
type CheckHandlerConfiguration struct {
	logger        log.Logger
	checks        map[string]check
	limit         int
	statusFailed  int
	statusSuccess int
}

// NewCheckHandlerConfiguration initializes a new CheckHandlerConfiguration.
func NewCheckHandlerConfiguration() CheckHandlerConfiguration {
	return CheckHandlerConfiguration{
		checks:        make(map[string]check),
		limit:         -1,
		statusFailed:  http.StatusInternalServerError,
		statusSuccess: http.StatusOK,
	}
}

// WithLogger sets the logger for the CheckHandlerConfiguration.
func (c CheckHandlerConfiguration) WithLogger(l log.Logger) CheckHandlerConfiguration {
	c.logger = l
	return c
}

// WithCheck sets a check for the CheckHandlerConfiguration.
func (c CheckHandlerConfiguration) WithCheck(name string, f check) CheckHandlerConfiguration {
	if _, ok := c.checks[name]; ok {
		c.logger.Panic().Str("check", name).Msg("check already exists")
	}

	c.checks[name] = f
	return c
}

// WithLimit limits the number of active goroutines for the checks to at most n
func (c CheckHandlerConfiguration) WithLimit(n int) CheckHandlerConfiguration {
	c.limit = n
	return c
}

// WithStatusFailed sets the status code for the failed checks.
func (c CheckHandlerConfiguration) WithStatusFailed(status int) CheckHandlerConfiguration {
	c.statusFailed = status
	return c
}

// WithStatusSuccess sets the status code for the successful checks.
func (c CheckHandlerConfiguration) WithStatusSuccess(status int) CheckHandlerConfiguration {
	c.statusSuccess = status
	return c
}

// CheckHandler is a http Handler that performs different checks.
type CheckHandler struct {
	conf CheckHandlerConfiguration
}

// NewCheckHandler initializes a new CheckHandler.
func NewCheckHandler(c CheckHandlerConfiguration) *CheckHandler {
	c.checks = maps.Clone(c.checks) // prevent check duplication after initialization
	return &CheckHandler{
		conf: c,
	}
}

// AddCheck adds a check to the CheckHandler.
func (h *CheckHandler) AddCheck(name string, c check) {
	h.conf.WithCheck(name, c)
}

func (h *CheckHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g, ctx := errgroup.WithContext(r.Context())
	g.SetLimit(h.conf.limit)

	for name, check := range h.conf.checks {
		checker := check
		checkerName := name
		g.Go(func() error { // https://go.dev/blog/loopvar-preview per iteration scope since go 1.22
			if err := checker(ctx); err != nil { // since go 1.22 for loops have a per-iteration scope instead of per-loop scope, no need to pin the check...
				return fmt.Errorf("'%s': %w", checkerName, err)
			}

			return nil
		})
	}

	status := h.conf.statusSuccess
	if err := g.Wait(); err != nil {
		status = h.conf.statusFailed
		h.conf.logger.Error().Err(err).Msg("check failed")
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)

	if _, err := io.WriteString(w, http.StatusText(status)); err != nil { // io.WriteString should not fail, but if it does, we want to know.
		h.conf.logger.Panic().Err(err).Msg("failed to write response")
	}
}
