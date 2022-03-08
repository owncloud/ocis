package svc

// NewTracing returns a service that instruments traces.
func NewTracing(next Service) Service {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next Service
}

// ListenForEvents implements service interface
func (t tracing) ListenForEvents() {
	t.next.ListenForEvents()
}
