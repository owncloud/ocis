package svc

import (
	"net/http"

	"github.com/owncloud/ocis/ocis-pkg/log"
)

// NewLogging returns a service that logs messages.
func NewLogging(next Service, logger log.Logger) Service {
	return logging{
		next:   next,
		logger: logger,
	}
}

type logging struct {
	next   Service
	logger log.Logger
}

// ServeHTTP implements the Service interface.
func (l logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.next.ServeHTTP(w, r)
}

// GetMe implements the Service interface.
func (l logging) GetMe(w http.ResponseWriter, r *http.Request) {
	l.next.GetMe(w, r)
}

// GetUsers implements the Service interface.
func (l logging) GetUsers(w http.ResponseWriter, r *http.Request) {
	l.next.GetUsers(w, r)
}

// GetUser implements the Service interface.
func (l logging) GetUser(w http.ResponseWriter, r *http.Request) {
	l.next.GetUser(w, r)
}

// PostUser implements the Service interface.
func (l logging) PostUser(w http.ResponseWriter, r *http.Request) {
	l.next.PostUser(w, r)
}

// DeleteUser implements the Service interface.
func (l logging) DeleteUser(w http.ResponseWriter, r *http.Request) {
	l.next.DeleteUser(w, r)
}

// PatchUser implements the Service interface.
func (l logging) PatchUser(w http.ResponseWriter, r *http.Request) {
	l.next.PatchUser(w, r)
}

// GetGroups implements the Service interface.
func (l logging) GetGroups(w http.ResponseWriter, r *http.Request) {
	l.next.GetGroups(w, r)
}

// GetGroup implements the Service interface.
func (l logging) GetGroup(w http.ResponseWriter, r *http.Request) {
	l.next.GetGroup(w, r)
}

// PostGroup implements the Service interface.
func (l logging) PostGroup(w http.ResponseWriter, r *http.Request) {
	l.next.PostGroup(w, r)
}

// PatchGroup implements the Service interface.
func (l logging) PatchGroup(w http.ResponseWriter, r *http.Request) {
	l.next.PatchGroup(w, r)
}

// DeleteGroup implements the Service interface.
func (l logging) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	l.next.DeleteGroup(w, r)
}

// GetGroupMembers implements the Service interface.
func (l logging) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	l.next.GetGroupMembers(w, r)
}

// PostGroupMember implements the Service interface.
func (l logging) PostGroupMember(w http.ResponseWriter, r *http.Request) {
	l.next.PostGroupMember(w, r)
}

// DeleteGroupMember implements the Service interface.
func (l logging) DeleteGroupMember(w http.ResponseWriter, r *http.Request) {
	l.next.DeleteGroupMember(w, r)
}

// GetDrives implements the Service interface.
func (l logging) GetDrives(w http.ResponseWriter, r *http.Request) {
	l.next.GetDrives(w, r)
}
