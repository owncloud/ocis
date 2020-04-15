package svc

// NewTracing returns a service that instruments traces.
func NewTracing(next Service) Service {
	return Service{}
}

type tracing struct {
	next Service
}
