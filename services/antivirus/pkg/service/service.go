package service

import (
	"bytes"
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cs3org/reva/v2/pkg/bytesize"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/antivirus/pkg/config"
	"github.com/owncloud/ocis/v2/services/antivirus/pkg/scanners"
)

// Scanner is an abstraction for the actual virus scan
type Scanner interface {
	Scan(file io.Reader) (scanners.ScanResult, error)
}

// NewAntivirus returns a service implementation for Service.
func NewAntivirus(c *config.Config, l log.Logger) (Antivirus, error) {
	av := Antivirus{c: c, l: l, client: rhttp.GetHTTPClient(rhttp.Insecure(true))}

	var err error
	av.s, err = scanners.New(c.Scanner)
	if err != nil {
		return av, err
	}

	switch o := events.PostprocessingOutcome(c.InfectedFileHandling); o {
	case events.PPOutcomeContinue, events.PPOutcomeAbort, events.PPOutcomeDelete:
		av.o = o
	default:
		return av, fmt.Errorf("unknown infected file handling '%s'", o)
	}

	if c.MaxScanSize != "" {
		b, err := bytesize.Parse(c.MaxScanSize)
		if err != nil {
			return av, err
		}

		av.m = b.Bytes()
	}

	return av, nil
}

// Antivirus defines implements the business logic for Service.
type Antivirus struct {
	c *config.Config
	l log.Logger
	s Scanner
	o events.PostprocessingOutcome
	m uint64

	client *http.Client
}

// Run runs the service
func (av Antivirus) Run() error {
	evtsCfg := av.c.Events

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

	stream, err := stream.NatsFromConfig(stream.NatsConfig(av.c.Events))
	if err != nil {
		return err
	}

	ch, err := events.Consume(stream, "antivirus", events.StartPostprocessingStep{})
	if err != nil {
		return err
	}

	for e := range ch {
		ev := e.Event.(events.StartPostprocessingStep)
		if ev.StepToStart != events.PPStepAntivirus {
			continue
		}

		av.l.Info().Str("uploadid", ev.UploadID).Str("filename", ev.Filename).Msg("Starting virus scan.")
		var errmsg string
		res, err := av.process(ev)
		if err != nil {
			errmsg = err.Error()
		}

		var outcome events.PostprocessingOutcome
		switch {
		case res.Infected:
			outcome = av.o
		case !res.Infected && err == nil:
			outcome = events.PPOutcomeContinue
		default:
			outcome = events.PPOutcomeAbort
		}

		av.l.Info().Str("uploadid", ev.UploadID).Interface("resourceID", ev.ResourceID).Str("virus", res.Description).Str("outcome", string(outcome)).Str("filename", ev.Filename).Str("user", ev.ExecutingUser.GetId().GetOpaqueId()).Bool("infected", res.Infected).Msg("File scanned")
		if err := events.Publish(stream, events.PostprocessingStepFinished{
			FinishedStep:  events.PPStepAntivirus,
			Outcome:       outcome,
			UploadID:      ev.UploadID,
			ExecutingUser: ev.ExecutingUser,
			Filename:      ev.Filename,
			Result: events.VirusscanResult{
				Infected:    res.Infected,
				Description: res.Description,
				Scandate:    time.Now(),
				ResourceID:  ev.ResourceID,
				ErrorMsg:    errmsg,
			},
		}); err != nil {
			av.l.Fatal().Err(err).Str("uploadid", ev.UploadID).Interface("resourceID", ev.ResourceID).Msg("cannot publish events - exiting")
			return err
		}
	}

	return nil
}

// process the scan
func (av Antivirus) process(ev events.StartPostprocessingStep) (scanners.ScanResult, error) {
	if ev.Filesize == 0 || (0 < av.m && av.m < ev.Filesize) {
		av.l.Debug().Str("uploadid", ev.UploadID).Uint64("limit", av.m).Uint64("filesize", ev.Filesize).Msg("Skipping file because file size is higher than the limit.")
		return scanners.ScanResult{
			Scantime: time.Now(),
		}, nil
	}

	var err error
	var rrc io.ReadCloser

	switch ev.UploadID {
	default:
		rrc, err = av.downloadViaToken(ev.URL)
	case "":
		rrc, err = av.downloadViaReva(ev.URL, ev.Token, ev.RevaToken)
	}
	if err != nil {
		av.l.Error().Err(err).Str("uploadid", ev.UploadID).Msg("error downloading file")
		return scanners.ScanResult{}, err
	}
	defer rrc.Close()
	av.l.Debug().Str("uploadid", ev.UploadID).Msg("Downloaded file successfully.")

	res, err := av.s.Scan(rrc)
	if err != nil {
		av.l.Error().Err(err).Str("uploadid", ev.UploadID).Msg("error scanning file")
	}

	return res, err

}

// download will download the file
func (av Antivirus) downloadViaToken(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return av.doDownload(req)
}

// download will download the file
func (av Antivirus) downloadViaReva(url string, dltoken string, revatoken string) (io.ReadCloser, error) {
	ctx := ctxpkg.ContextSetToken(context.Background(), revatoken)

	req, err := rhttp.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Reva-Transfer", dltoken)

	return av.doDownload(req)
}

func (av Antivirus) doDownload(req *http.Request) (io.ReadCloser, error) {
	res, err := av.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		return nil, fmt.Errorf("unexpected status code from Download %v", res.StatusCode)
	}

	return res.Body, nil
}
