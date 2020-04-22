package svc

// NewTracing returns a service that instruments traces.
func NewTracing(next Service) Service {
	return Service{
		manager: next.manager,
		config:  next.config,
	}
}

type tracing struct {
	next Service
}
