package handlers

import (
	"context"
	"fmt"
	"io"
	"maps"
	"net"
	"net/http"
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"golang.org/x/sync/errgroup"
)

// check is a function that performs a check.
type checker func(ctx context.Context) error

// CheckHandlerConfiguration defines the configuration for the CheckHandler.
type CheckHandlerConfiguration struct {
	checks        map[string]checker
	logger        log.Logger
	limit         int
	statusFailed  int
	statusSuccess int
}

// NewCheckHandlerConfiguration initializes a new CheckHandlerConfiguration.
func NewCheckHandlerConfiguration() CheckHandlerConfiguration {
	return CheckHandlerConfiguration{
		checks: make(map[string]checker),

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
func (c CheckHandlerConfiguration) WithCheck(name string, check checker) CheckHandlerConfiguration {
	if _, ok := c.checks[name]; ok {
		c.logger.Panic().Str("check", name).Msg("check already exists")
	}

	c.checks = maps.Clone(c.checks) // prevent propagated check duplication, maps are references;
	c.checks[name] = check

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

// FailSaveAddress replaces unspecified addresses with the outbound IP.
func FailSaveAddress(address string) (string, error) {
	host, port := SplitHostPort(address)

	hostIP := net.ParseIP(host)

	if host == "" || (hostIP != nil && hostIP.IsUnspecified()) {
		outboundIP, err := getOutBoundIP()
		if err != nil {
			return "", err
		}

		host = outboundIP.String()
	}

	if port != "" {
		if strings.Contains(host, ":") {
			host = "[" + host + "]"
		}
		return host + ":" + port, nil
	}

	return host, nil
}

// SplitHostPort returns host and port of the address.
// Contrary to the net.SplitHostPort the port is not mandatory.
func SplitHostPort(address string) (string, string) {
	columns := strings.Split(address, ":")
	brackets := strings.Split(address, "]")

	switch {
	case len(columns) == 1 && len(brackets) == 1: // 10.10.10.10
		return address, ""
	case len(columns) == 2 && len(brackets) == 1: // 10.10.10.10:80
		return columns[0], columns[1]
	case len(columns) > 2 && len(brackets) == 1: // 2a01::a
		return address, ""
	case len(brackets) == 2 && brackets[1] == "": // [2a01::a]
		return brackets[0][1:], ""
	case len(brackets) == 2: // [2a01::a]:10
		return brackets[0][1:], columns[len(columns)-1]
	}

	return address, ""
}

// getOutBoundIP returns the outbound IP address.
func getOutBoundIP() (net.IP, error) {
	interfacesAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var fallbackIpv6 net.IP
	for _, address := range interfacesAddresses {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP, nil
			}
			if ipNet.IP.To16() != nil && !ipNet.IP.IsLinkLocalUnicast() {
				fallbackIpv6 = ipNet.IP.To16()
			}
		}
	}

	if fallbackIpv6 != nil {
		return fallbackIpv6, nil
	}

	return nil, fmt.Errorf("no IP found")
}
