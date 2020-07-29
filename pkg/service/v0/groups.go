package svc

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	accounts "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis-ocs/pkg/service/v0/response"
)

// ListUserGroups lists a users groups
func (o Ocs) ListUserGroups(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, response.DataRender(&data.Groups{Groups: []string{}}))
}

// AddToGroup adds a user to a group
func (o Ocs) AddToGroup(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, response.ErrRender(data.MetaUnknownError.StatusCode, "not implemented"))
}

// RemoveFromGroup removes a user from a group
func (o Ocs) RemoveFromGroup(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, response.ErrRender(data.MetaUnknownError.StatusCode, "not implemented"))
}

// ListGroups lists all groups
func (o Ocs) ListGroups(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	query := ""
	if search != "" {
		query = fmt.Sprintf("id eq '%s' or on_premises_sam_account_name eq '%s'", escapeValue(search), escapeValue(search))
	}
	accSvc := o.getGroupsService()
	res, err := accSvc.ListGroups(r.Context(), &accounts.ListGroupsRequest{
		Query: query,
	})
	if err != nil {
		o.logger.Err(err).Msg("could not list users")
		render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not list users"))
		return
	}
	groups := []string{}
	for i := range res.Groups {
		groups = append(groups, res.Groups[i].Id)
	}

	render.Render(w, r, response.DataRender(&data.Groups{Groups: groups}))
}

// AddGroup adds a group
func (o Ocs) AddGroup(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, response.ErrRender(data.MetaUnknownError.StatusCode, "not implemented"))
}

// DeleteGroup deletes a group
func (o Ocs) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, response.ErrRender(data.MetaUnknownError.StatusCode, "not implemented"))
}

// GetGroupMembers lists all members of a group
func (o Ocs) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, response.ErrRender(data.MetaUnknownError.StatusCode, "not implemented"))
}
