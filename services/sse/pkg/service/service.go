package service

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/rhttp"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/sse/pkg/config"
)

// NewSSE returns a service implementation for Service.
func NewSSE(c *config.Config, l log.Logger) (SSE, error) {
	s := SSE{c: c, l: l, client: rhttp.GetHTTPClient(rhttp.Insecure(true))}

	return s, nil
}

// SSE defines implements the business logic for Service.
type SSE struct {
	c *config.Config
	l log.Logger
	m uint64

	client *http.Client
}

// Run runs the service
func (s SSE) Run() error {
	evtsCfg := s.c.Events

	var rootCAPool *x509.CertPool
	if evtsCfg.TLSRootCACertificate != "" {
		rootCrtFile, err := os.Open(evtsCfg.TLSRootCACertificate)
		if err != nil {
			return err
		}

		var certBytes bytes.Buffer
		if _, err := io.Copy(&certBytes, rootCrtFile); err != nil {
			return err
		}

		rootCAPool = x509.NewCertPool()
		rootCAPool.AppendCertsFromPEM(certBytes.Bytes())
		evtsCfg.TLSInsecure = false
	}

	natsStream, err := stream.NatsFromConfig(stream.NatsConfig(s.c.Events))
	if err != nil {
		return err
	}

	ch, err := events.Consume(natsStream, "sse", events.StartPostprocessingStep{})
	if err != nil {
		return err
	}

	for e := range ch {
		fmt.Println(e) // todo
	}

	return nil
}
