package svc

import (
	"net/http"

	graphTracing "github.com/owncloud/ocis/v2/services/graph/pkg/tracing"
	"go.opentelemetry.io/otel/propagation"
)

const tracer = "graph"

var propagator = propagation.NewCompositeTextMapPropagator(
	propagation.Baggage{},
	propagation.TraceContext{},
)

// NewTracing returns a service that instruments traces.
func NewTracing(next Service) Service {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next Service
}

// ServeHTTP implements the Service interface.
func (t tracing) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
	propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))

	//pkgmiddleware.TraceContext(t.next).ServeHTTP(w, r.WithContext(r.Context()))

	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(ctx, "serve http")
	defer span.End()
	t.next.ServeHTTP(w, r.WithContext(ctx))
}

// GetMe implements the Service interface.
func (t tracing) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "GetMe")
	defer span.End()
	t.next.GetMe(w, r.WithContext(ctx))
}

// GetUsers implements the Service interface.
func (t tracing) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "GetUsers")
	defer span.End()
	t.next.GetUsers(w, r.WithContext(ctx))
}

// GetUser implements the Service interface.
func (t tracing) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "GetUser")
	defer span.End()
	t.next.GetUser(w, r.WithContext(ctx))
}

// PostUser implements the Service interface.
func (t tracing) PostUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "PostUser")
	defer span.End()
	t.next.PostUser(w, r.WithContext(ctx))
}

// DeleteUser implements the Service interface.
func (t tracing) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "DeleteUser")
	defer span.End()
	t.next.DeleteUser(w, r.WithContext(ctx))
}

// PatchUser implements the Service interface.
func (t tracing) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "PatchUser")
	defer span.End()
	t.next.PatchUser(w, r.WithContext(ctx))
}

// ChangeOwnPassword implements the Service interface.
func (t tracing) ChangeOwnPassword(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "ChangeOwnPassword")
	defer span.End()
	t.next.ChangeOwnPassword(w, r.WithContext(ctx))
}

// GetGroups implements the Service interface.
func (t tracing) GetGroups(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "GetGroups")
	defer span.End()
	t.next.GetGroups(w, r.WithContext(ctx))
}

// GetGroup implements the Service interface.
func (t tracing) GetGroup(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "GetGroup")
	defer span.End()
	t.next.GetGroup(w, r.WithContext(ctx))
}

// PostGroup implements the Service interface.
func (t tracing) PostGroup(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "PostGroup")
	defer span.End()
	t.next.PostGroup(w, r.WithContext(ctx))
}

// PatchGroup implements the Service interface.
func (t tracing) PatchGroup(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "PatchGroup")
	defer span.End()
	t.next.PatchGroup(w, r.WithContext(ctx))
}

// DeleteGroup implements the Service interface.
func (t tracing) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "DeleteGroup")
	defer span.End()
	t.next.DeleteGroup(w, r.WithContext(ctx))
}

// GetGroupMembers implements the Service interface.
func (t tracing) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "GetGroupMembers")
	defer span.End()
	t.next.GetGroupMembers(w, r.WithContext(ctx))
}

// PostGroupMember implements the Service interface.
func (t tracing) PostGroupMember(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "PostGroupMember")
	defer span.End()
	t.next.PostGroupMember(w, r.WithContext(ctx))
}

// DeleteGroupMember implements the Service interface.
func (t tracing) DeleteGroupMember(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "DeleteGroupMember")
	defer span.End()
	t.next.DeleteGroupMember(w, r.WithContext(ctx))
}

// GetDrives implements the Service interface.
func (t tracing) GetDrives(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "GetDrives")
	defer span.End()
	t.next.GetDrives(w, r.WithContext(ctx))
}

// GetAllDrives implements the Service interface.
func (t tracing) GetAllDrives(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "GetAllDrives")
	defer span.End()
	t.next.GetAllDrives(w, r.WithContext(ctx))
}

// CreateDrive implements the Service interface.
func (t tracing) CreateDrive(w http.ResponseWriter, r *http.Request) {
	ctx, span := graphTracing.TraceProvider.Tracer(tracer).Start(r.Context(), "CreateDrive")
	defer span.End()
	t.next.CreateDrive(w, r.WithContext(ctx))
}
