package eventSVC

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/services/authz/pkg/authz"
	"github.com/owncloud/ocis/v2/services/authz/pkg/config"
	"os"
)

// Service defines the service handlers.

type Service struct {
	stream      events.Stream
	authorizers []authz.Authorizer
}

// New returns a service implementation for Service.
func New(cfg *config.Config, authorizers []authz.Authorizer) (Service, error) {
	evtsCfg := cfg.Events
	svc := Service{
		authorizers: authorizers,
	}

	var tlsConf *tls.Config
	if evtsCfg.EnableTLS {
		var rootCAPool *x509.CertPool
		if evtsCfg.TLSRootCACertificate != "" {
			rootCrtFile, err := os.Open(evtsCfg.TLSRootCACertificate)
			if err != nil {
				return svc, err
			}

			rootCAPool, err = crypto.NewCertPoolFromPEM(rootCrtFile)
			if err != nil {
				return svc, err
			}
			evtsCfg.TLSInsecure = false
		}

		tlsConf = &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: evtsCfg.TLSInsecure, //nolint:gosec
			RootCAs:            rootCAPool,
		}
	}
	eventStream, err := stream.Nats(
		natsjs.TLSConfig(tlsConf),
		natsjs.Address(evtsCfg.Endpoint),
		natsjs.ClusterID(evtsCfg.Cluster),
	)
	if err != nil {
		return svc, err
	}

	svc.stream = eventStream

	return svc, nil
}

func (s Service) Run() error {
	ch, err := events.Consume(s.stream, "authz", events.StartPostprocessingStep{})
	if err != nil {
		return err
	}

	for ce := range ch {
		ev := ce.(events.StartPostprocessingStep)
		if ev.StepToStart != "authz" {
			continue
		}

		env := authz.Environment{
			Name:       ev.Filename,
			URL:        ev.URL,
			Size:       ev.Filesize,
			User:       *ev.ExecutingUser,
			ResourceID: *ev.ResourceID,
			Stage:      authz.StagePP,
		}

		allowed, err := authz.Authorized(context.TODO(), env, s.authorizers...)
		if err != nil {
			return err
		}

		outcome := events.PPOutcomeContinue
		if !allowed {
			outcome = events.PPOutcomeDelete
		}

		if err := events.Publish(s.stream, events.PostprocessingStepFinished{
			Outcome:       outcome,
			UploadID:      ev.UploadID,
			ExecutingUser: ev.ExecutingUser,
			Filename:      ev.Filename,
			FinishedStep:  ev.StepToStart,
		}); err != nil {
			return err
		}
	}

	return nil
}
