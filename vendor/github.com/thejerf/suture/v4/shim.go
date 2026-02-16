package suture

import (
	"context"
)

type DeprecatedService interface {
	Serve()
	Stop()
}

// AsService converts old-style suture service to a new style suture service.
func AsService(service DeprecatedService) Service {
	return &serviceShim{service: service}
}

type serviceShim struct {
	service DeprecatedService
}

func (s *serviceShim) Serve(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.service.Serve()
		close(done)
	}()

	select {
	case <-done:
		// If the service stops by itself (done closes), return straight away, there is no error, and we don't need
		// to wait for the context.
		return nil
	case <-ctx.Done():
		// If the context is closed, stop the service, then wait for it's termination and return the error from the
		// context.
		s.service.Stop()
		<-done
		return ctx.Err()
	}
}
