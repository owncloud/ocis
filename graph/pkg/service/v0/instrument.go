package svc

import (
	"net/http"

	"github.com/owncloud/ocis/graph/pkg/metrics"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next Service, metrics *metrics.Metrics) Service {
	return instrument{
		next:    next,
		metrics: metrics,
	}
}

type instrument struct {
	next    Service
	metrics *metrics.Metrics
}

// ServeHTTP implements the Service interface.
func (i instrument) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.next.ServeHTTP(w, r)
}

// GetMe implements the Service interface.
func (i instrument) GetMe(w http.ResponseWriter, r *http.Request) {
	i.next.GetMe(w, r)
}

// GetUsers implements the Service interface.
func (i instrument) GetUsers(w http.ResponseWriter, r *http.Request) {
	i.next.GetUsers(w, r)
}

// GetUser implements the Service interface.
func (i instrument) GetUser(w http.ResponseWriter, r *http.Request) {
	i.next.GetUser(w, r)
}

// PostUser implements the Service interface.
func (i instrument) PostUser(w http.ResponseWriter, r *http.Request) {
	i.next.PostUser(w, r)
}

// DeleteUser implements the Service interface.
func (i instrument) DeleteUser(w http.ResponseWriter, r *http.Request) {
	i.next.DeleteUser(w, r)
}

// PatchUser implements the Service interface.
func (i instrument) PatchUser(w http.ResponseWriter, r *http.Request) {
	i.next.PatchUser(w, r)
}

// GetGroups implements the Service interface.
func (i instrument) GetGroups(w http.ResponseWriter, r *http.Request) {
	i.next.GetGroups(w, r)
}

// GetGroup implements the Service interface.
func (i instrument) GetGroup(w http.ResponseWriter, r *http.Request) {
	i.next.GetGroup(w, r)
}

// PostGroup implements the Service interface.
func (i instrument) PostGroup(w http.ResponseWriter, r *http.Request) {
	i.next.PostGroup(w, r)
}

// DeleteGroup implements the Service interface.
func (i instrument) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	i.next.DeleteGroup(w, r)
}

// GetGroupMembers implements the Service interface.
func (i instrument) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	i.next.GetGroupMembers(w, r)
}

// PostGroupMember implements the Service interface.
func (i instrument) PostGroupMember(w http.ResponseWriter, r *http.Request) {
	i.next.PostGroupMember(w, r)
}

// GetDrives implements the Service interface.
func (i instrument) GetDrives(w http.ResponseWriter, r *http.Request) {
	i.next.GetDrives(w, r)
}
