package svc

import (
	"net/http"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next Service) tracing {
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

// GetApplication implements the Service interface.
func (t tracing) GetApplication(w http.ResponseWriter, r *http.Request) {
	t.next.GetMe(w, r)
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

// ChangeOwnPassword implements the Service interface.
func (t tracing) ChangeOwnPassword(w http.ResponseWriter, r *http.Request) {
	t.next.ChangeOwnPassword(w, r)
}

// ListAppRoleAssignments implements the Service interface.
func (t tracing) ListAppRoleAssignments(w http.ResponseWriter, r *http.Request) {
	t.next.ListAppRoleAssignments(w, r)
}

// CreateAppRoleAssignment implements the Service interface.
func (t tracing) CreateAppRoleAssignment(w http.ResponseWriter, r *http.Request) {
	t.next.CreateAppRoleAssignment(w, r)
}

// DeleteAppRoleAssignment implements the Service interface.
func (t tracing) DeleteAppRoleAssignment(w http.ResponseWriter, r *http.Request) {
	t.next.DeleteAppRoleAssignment(w, r)
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

// GetEducationSchools implements the Service interface.
func (t tracing) GetEducationSchools(w http.ResponseWriter, r *http.Request) {
	t.next.GetEducationSchools(w, r)
}

// GetEducationSchool implements the Service interface.
func (t tracing) GetEducationSchool(w http.ResponseWriter, r *http.Request) {
	t.next.GetEducationSchool(w, r)
}

// PostEducationSchool implements the Service interface.
func (t tracing) PostEducationSchool(w http.ResponseWriter, r *http.Request) {
	t.next.PostEducationSchool(w, r)
}

// PatchEducationSchool implements the Service interface.
func (t tracing) PatchEducationSchool(w http.ResponseWriter, r *http.Request) {
	t.next.PatchEducationSchool(w, r)
}

// DeleteEducationSchool implements the Service interface.
func (t tracing) DeleteEducationSchool(w http.ResponseWriter, r *http.Request) {
	t.next.DeleteEducationSchool(w, r)
}

// GetEducationSchoolUsers implements the Service interface.
func (t tracing) GetEducationSchoolUsers(w http.ResponseWriter, r *http.Request) {
	t.next.GetEducationSchoolUsers(w, r)
}

// PostEducationSchoolUser implements the Service interface.
func (t tracing) PostEducationSchoolUser(w http.ResponseWriter, r *http.Request) {
	t.next.PostEducationSchoolUser(w, r)
}

// DeleteEducationSchoolUser implements the Service interface.
func (t tracing) DeleteEducationSchoolUser(w http.ResponseWriter, r *http.Request) {
	t.next.DeleteEducationSchoolUser(w, r)
}

// GetEducationUsers implements the Service interface.
func (t tracing) GetEducationUsers(w http.ResponseWriter, r *http.Request) {
	t.next.GetEducationUsers(w, r)
}

// GetEducationUser implements the Service interface.
func (t tracing) GetEducationUser(w http.ResponseWriter, r *http.Request) {
	t.next.GetEducationUser(w, r)
}

// PostEducationUser implements the Service interface.
func (t tracing) PostEducationUser(w http.ResponseWriter, r *http.Request) {
	t.next.PostEducationUser(w, r)
}

// DeleteEducationUser implements the Service interface.
func (t tracing) DeleteEducationUser(w http.ResponseWriter, r *http.Request) {
	t.next.DeleteEducationUser(w, r)
}

// PatchEducationUser implements the Service interface.
func (t tracing) PatchEducationUser(w http.ResponseWriter, r *http.Request) {
	t.next.PatchEducationUser(w, r)
}

// GetDrives implements the Service interface.
func (t tracing) GetDrives(w http.ResponseWriter, r *http.Request) {
	t.next.GetDrives(w, r)
}

// GetSingleDrive implements the Service interface.
func (t tracing) GetSingleDrive(w http.ResponseWriter, r *http.Request) {
	t.next.GetDrives(w, r)
}

// UpdateDrive implements the Service interface.
func (t tracing) UpdateDrive(w http.ResponseWriter, r *http.Request) {
	t.next.GetDrives(w, r)
}

// DeleteDrive implements the Service interface.
func (t tracing) DeleteDrive(w http.ResponseWriter, r *http.Request) {
	t.next.GetDrives(w, r)
}

// GetAllDrives implements the Service interface.
func (t tracing) GetAllDrives(w http.ResponseWriter, r *http.Request) {
	t.next.GetAllDrives(w, r)
}

// CreateDrive implements the Service interface.
func (t tracing) CreateDrive(w http.ResponseWriter, r *http.Request) {
	t.next.CreateDrive(w, r)
}

// GetRootDriveChildren implements the Service interface.
func (t tracing) GetRootDriveChildren(w http.ResponseWriter, r *http.Request) {
	t.next.GetRootDriveChildren(w, r)
}

// GetTags implements the Service interface.
func (t tracing) GetTags(w http.ResponseWriter, r *http.Request) {
	t.next.GetTags(w, r)
}

// AssignTags implements the Service interface.
func (t tracing) AssignTags(w http.ResponseWriter, r *http.Request) {
	t.next.AssignTags(w, r)
}

// UnassignTags implements the Service interface.
func (t tracing) UnassignTags(w http.ResponseWriter, r *http.Request) {
	t.next.UnassignTags(w, r)
}
