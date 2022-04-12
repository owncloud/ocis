package svc

import (
	"net/http"
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
	t.next.ServeHTTP(w, r)
}

// GetMe implements the Service interface.
func (t tracing) GetMe(w http.ResponseWriter, r *http.Request) {
	t.next.GetMe(w, r)
}

// GetUsers implements the Service interface.
func (t tracing) GetUsers(w http.ResponseWriter, r *http.Request) {
	t.next.GetUsers(w, r)
}

// GetUser implements the Service interface.
func (t tracing) GetUser(w http.ResponseWriter, r *http.Request) {
	t.next.GetUser(w, r)
}

// PostUser implements the Service interface.
func (t tracing) PostUser(w http.ResponseWriter, r *http.Request) {
	t.next.PostUser(w, r)
}

// DeleteUser implements the Service interface.
func (t tracing) DeleteUser(w http.ResponseWriter, r *http.Request) {
	t.next.DeleteUser(w, r)
}

// PatchUser implements the Service interface.
func (t tracing) PatchUser(w http.ResponseWriter, r *http.Request) {
	t.next.PatchUser(w, r)
}

// GetGroups implements the Service interface.
func (t tracing) GetGroups(w http.ResponseWriter, r *http.Request) {
	t.next.GetGroups(w, r)
}

// GetGroup implements the Service interface.
func (t tracing) GetGroup(w http.ResponseWriter, r *http.Request) {
	t.next.GetGroup(w, r)
}

// PostGroup implements the Service interface.
func (t tracing) PostGroup(w http.ResponseWriter, r *http.Request) {
	t.next.PostGroup(w, r)
}

// PatchGroup implements the Service interface.
func (t tracing) PatchGroup(w http.ResponseWriter, r *http.Request) {
	t.next.PatchGroup(w, r)
}

// DeleteGroup implements the Service interface.
func (t tracing) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	t.next.DeleteGroup(w, r)
}

// GetGroupMembers implements the Service interface.
func (t tracing) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	t.next.GetGroupMembers(w, r)
}

// PostGroupMember implements the Service interface.
func (t tracing) PostGroupMember(w http.ResponseWriter, r *http.Request) {
	t.next.PostGroupMember(w, r)
}

// DeleteGroupMember implements the Service interface.
func (t tracing) DeleteGroupMember(w http.ResponseWriter, r *http.Request) {
	t.next.DeleteGroupMember(w, r)
}

// GetDrives implements the Service interface.
func (t tracing) GetDrives(w http.ResponseWriter, r *http.Request) {
	t.next.GetDrives(w, r)
}
