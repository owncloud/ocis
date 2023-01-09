package svc

import (
	"net/http"

	"github.com/owncloud/ocis/v2/services/graph/pkg/metrics"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next Service, metrics *metrics.Metrics) instrument {
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

// ChangeOwnPassword implements the Service interface.
func (i instrument) ChangeOwnPassword(w http.ResponseWriter, r *http.Request) {
	i.next.ChangeOwnPassword(w, r)
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

// PatchGroup implements the Service interface.
func (i instrument) PatchGroup(w http.ResponseWriter, r *http.Request) {
	i.next.PatchGroup(w, r)
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

// DeleteGroupMember implements the Service interface.
func (i instrument) DeleteGroupMember(w http.ResponseWriter, r *http.Request) {
	i.next.DeleteGroupMember(w, r)
}

// GetEducationSchools implements the Service interface.
func (i instrument) GetEducationSchools(w http.ResponseWriter, r *http.Request) {
	i.next.GetEducationSchools(w, r)
}

// GetEducationSchool implements the Service interface.
func (i instrument) GetEducationSchool(w http.ResponseWriter, r *http.Request) {
	i.next.GetEducationSchool(w, r)
}

// PostEducationSchool implements the Service interface.
func (i instrument) PostEducationSchool(w http.ResponseWriter, r *http.Request) {
	i.next.PostEducationSchool(w, r)
}

// PatchEducationSchool implements the Service interface.
func (i instrument) PatchEducationSchool(w http.ResponseWriter, r *http.Request) {
	i.next.PatchEducationSchool(w, r)
}

// DeleteEducationSchool implements the Service interface.
func (i instrument) DeleteEducationSchool(w http.ResponseWriter, r *http.Request) {
	i.next.DeleteEducationSchool(w, r)
}

// GetEducationSchoolUsers implements the Service interface.
func (i instrument) GetEducationSchoolUsers(w http.ResponseWriter, r *http.Request) {
	i.next.GetEducationSchoolUsers(w, r)
}

// PostEducationSchoolUser implements the Service interface.
func (i instrument) PostEducationSchoolUser(w http.ResponseWriter, r *http.Request) {
	i.next.PostEducationSchoolUser(w, r)
}

// DeleteEducationSchoolUser implements the Service interface.
func (i instrument) DeleteEducationSchoolUser(w http.ResponseWriter, r *http.Request) {
	i.next.DeleteEducationSchoolUser(w, r)
}

// GetEducationClasses implements the Service interface.
func (i instrument) GetEducationClasses(w http.ResponseWriter, r *http.Request) {
	i.next.GetEducationClasses(w, r)
}

// GetEducationClass implements the Service interface.
func (i instrument) GetEducationClass(w http.ResponseWriter, r *http.Request) {
	i.next.GetEducationClass(w, r)
}

// PostEducationClass implements the Service interface.
func (i instrument) PostEducationClass(w http.ResponseWriter, r *http.Request) {
	i.next.PostEducationClass(w, r)
}

// PatchEducationClass implements the Service interface.
func (i instrument) PatchEducationClass(w http.ResponseWriter, r *http.Request) {
	i.next.PatchEducationClass(w, r)
}

// DeleteEducationClass implements the Service interface.
func (i instrument) DeleteEducationClass(w http.ResponseWriter, r *http.Request) {
	i.next.DeleteEducationClass(w, r)
}

// GetEducationClassMembers implements the Service interface.
func (i instrument) GetEducationClassMembers(w http.ResponseWriter, r *http.Request) {
	i.next.GetEducationClassMembers(w, r)
}

// PostEducationClassMember implements the Service interface.
func (i instrument) PostEducationClassMember(w http.ResponseWriter, r *http.Request) {
	i.next.PostEducationClassMember(w, r)
}

// DeleteEducationClassMember implements the Service interface.
func (i instrument) DeleteEducationClassMember(w http.ResponseWriter, r *http.Request) {
	i.next.DeleteEducationClassMember(w, r)
}

// GetEducationUsers implements the Service interface.
func (i instrument) GetEducationUsers(w http.ResponseWriter, r *http.Request) {
	i.next.GetEducationUsers(w, r)
}

// GetEducationUser implements the Service interface.
func (i instrument) GetEducationUser(w http.ResponseWriter, r *http.Request) {
	i.next.GetEducationUser(w, r)
}

// PostEducationUser implements the Service interface.
func (i instrument) PostEducationUser(w http.ResponseWriter, r *http.Request) {
	i.next.PostEducationUser(w, r)
}

// DeleteEducationUser implements the Service interface.
func (i instrument) DeleteEducationUser(w http.ResponseWriter, r *http.Request) {
	i.next.DeleteEducationUser(w, r)
}

// PatchEducationUser implements the Service interface.
func (i instrument) PatchEducationUser(w http.ResponseWriter, r *http.Request) {
	i.next.PatchEducationUser(w, r)
}

// GetDrives implements the Service interface.
func (i instrument) GetDrives(w http.ResponseWriter, r *http.Request) {
	i.next.GetDrives(w, r)
}

// GetSingleDrive implements the Service interface.
func (i instrument) GetSingleDrive(w http.ResponseWriter, r *http.Request) {
	i.next.GetDrives(w, r)
}

// UpdateDrive implements the Service interface.
func (i instrument) UpdateDrive(w http.ResponseWriter, r *http.Request) {
	i.next.GetDrives(w, r)
}

// DeleteDrive implements the Service interface.
func (i instrument) DeleteDrive(w http.ResponseWriter, r *http.Request) {
	i.next.GetDrives(w, r)
}

// GetAllDrives implements the Service interface.
func (i instrument) GetAllDrives(w http.ResponseWriter, r *http.Request) {
	i.next.GetAllDrives(w, r)
}

// CreateDrive implements the Service interface.
func (i instrument) CreateDrive(w http.ResponseWriter, r *http.Request) {
	i.next.CreateDrive(w, r)
}

// GetRootDriveChildren implements the Service interface.
func (i instrument) GetRootDriveChildren(w http.ResponseWriter, r *http.Request) {
	i.next.GetRootDriveChildren(w, r)
}

// GetTags implements the Service interface.
func (i instrument) GetTags(w http.ResponseWriter, r *http.Request) {
	i.next.GetTags(w, r)
}

// AssignTags implements the Service interface.
func (i instrument) AssignTags(w http.ResponseWriter, r *http.Request) {
	i.next.AssignTags(w, r)
}

// UnassignTags implements the Service interface.
func (i instrument) UnassignTags(w http.ResponseWriter, r *http.Request) {
	i.next.UnassignTags(w, r)
}
