package svc

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	merrors "github.com/micro/go-micro/v2/errors"

	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/response"
)

// ListUserGroups lists a users groups
func (o Ocs) ListUserGroups(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")
	account, err := o.fetchAccountByUsername(r.Context(), userid)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found"))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(err).Str("userid", userid).Msg("could not get list of user groups")
		return
	}

	groups := []string{}
	for i := range account.MemberOf {
		groups = append(groups, account.MemberOf[i].Id)
	}

	o.logger.Error().Err(err).Int("count", len(groups)).Str("userid", userid).Msg("listing groups for user")
	render.Render(w, r, response.DataRender(&data.Groups{Groups: groups}))
}

// AddToGroup adds a user to a group
func (o Ocs) AddToGroup(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userid := chi.URLParam(r, "userid")
	groupid := r.PostForm.Get("groupid")

	if groupid == "" {
		render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "empty group assignment: unspecified group"))
		return
	}
	account, err := o.fetchAccountByUsername(r.Context(), userid)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found"))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		return
	}

	_, err = o.getGroupsService().AddMember(r.Context(), &accounts.AddMemberRequest{
		AccountId: account.Id,
		GroupId:   groupid,
	})

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found"))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(err).Str("userid", userid).Str("groupid", groupid).Msg("could not add user to group")
		return
	}

	o.logger.Debug().Str("userid", userid).Str("groupid", groupid).Msg("added user to group")
	render.Render(w, r, response.DataRender(struct{}{}))
}

// RemoveFromGroup removes a user from a group
func (o Ocs) RemoveFromGroup(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")
	groupid := r.URL.Query().Get("groupid")

	account, err := o.fetchAccountByUsername(r.Context(), userid)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found"))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		return
	}

	_, err = o.getGroupsService().RemoveMember(r.Context(), &accounts.RemoveMemberRequest{
		AccountId: account.Id,
		GroupId:   groupid,
	})

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found"))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(err).Str("userid", userid).Str("groupid", groupid).Msg("could not remove user from group")
		return
	}

	o.logger.Debug().Str("userid", userid).Str("groupid", groupid).Msg("removed user from group")
	render.Render(w, r, response.DataRender(struct{}{}))
}

// ListGroups lists all groups
func (o Ocs) ListGroups(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	query := ""
	if search != "" {
		query = fmt.Sprintf("id eq '%s' or on_premises_sam_account_name eq '%s'", escapeValue(search), escapeValue(search))
	}

	res, err := o.getGroupsService().ListGroups(r.Context(), &accounts.ListGroupsRequest{
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
	groupid := chi.URLParam(r, "groupid")

	_, err := o.getGroupsService().DeleteGroup(r.Context(), &accounts.DeleteGroupRequest{
		Id: groupid,
	})

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found"))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(err).Str("groupid", groupid).Msg("could not remove group")
		return
	}

	o.logger.Debug().Str("groupid", groupid).Msg("removed group")
	render.Render(w, r, response.DataRender(struct{}{}))
}

// GetGroupMembers lists all members of a group
func (o Ocs) GetGroupMembers(w http.ResponseWriter, r *http.Request) {

	groupid := chi.URLParam(r, "groupid")

	res, err := o.getGroupsService().ListMembers(r.Context(), &accounts.ListMembersRequest{Id: groupid})

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found"))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(err).Str("groupid", groupid).Msg("could not get list of members")
		return
	}

	members := []string{}
	for i := range res.Members {
		members = append(members, res.Members[i].Id)
	}

	o.logger.Error().Err(err).Int("count", len(members)).Str("groupid", groupid).Msg("listing group members")
	render.Render(w, r, response.DataRender(&data.Users{Users: members}))
}
