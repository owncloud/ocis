package svc

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	merrors "github.com/asim/go-micro/v3/errors"
	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/response"
)

// ListUserGroups lists a users groups
func (o Ocs) ListUserGroups(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")
	var account *accounts.Account
	var err error

	// short circuit if there is a user already in the context
	if u, ok := revactx.ContextGetUser(r.Context()); ok {
		// we are not sure whether the current user in the context is the admin or the authenticated user.
		if u.Username == userid {
			// the OCS API is a REST API and it uses the username to look for groups. If the id from the user in the context
			// differs from that of the url we can assume we are an admin because we are past the selfOrAdmin middleware.
			if len(u.Groups) > 0 {
				mustNotFail(render.Render(w, r, response.DataRender(&data.Groups{Groups: u.Groups})))
				return
			}
		}
	}

	if isValidUUID(userid) {
		account, err = o.getAccountService().GetAccount(r.Context(), &accounts.GetAccountRequest{
			Id: userid,
		})
	} else {
		// despite the confusion, if we make it here we got ourselves a username
		account, err = o.fetchAccountByUsername(r.Context(), userid)
		if err != nil {
			merr := merrors.FromError(err)
			if merr.Code == http.StatusNotFound {
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found")))
			} else {
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
			}
			o.logger.Error().Err(err).Str("userid", userid).Msg("could not get list of user groups")
			return
		}
	}

	groups := make([]string, 0, len(account.MemberOf))
	for i := range account.MemberOf {
		if account.MemberOf[i].OnPremisesSamAccountName == "" {
			o.logger.Warn().Str("groupid", account.MemberOf[i].Id).Msg("group on_premises_sam_account_name is empty, trying to lookup by id")
			// we can try to look up the name
			group, err := o.getGroupsService().GetGroup(r.Context(), &accounts.GetGroupRequest{
				Id: account.MemberOf[i].Id,
			})

			if err != nil {
				o.logger.Error().Err(err).Str("groupid", account.MemberOf[i].Id).Msg("could not get group")
				continue
			}
			if group.OnPremisesSamAccountName == "" {
				o.logger.Error().Err(err).Str("groupid", account.MemberOf[i].Id).Msg("group on_premises_sam_account_name is empty")
				continue
			}
			groups = append(groups, group.OnPremisesSamAccountName)
		} else {
			groups = append(groups, account.MemberOf[i].OnPremisesSamAccountName)
		}
	}

	o.logger.Error().Err(err).Int("count", len(groups)).Str("userid", account.Id).Msg("listing groups for user")
	mustNotFail(render.Render(w, r, response.DataRender(&data.Groups{Groups: groups})))
}

// AddToGroup adds a user to a group
func (o Ocs) AddToGroup(w http.ResponseWriter, r *http.Request) {
	mustNotFail(r.ParseForm())
	userid := chi.URLParam(r, "userid")
	groupid := r.PostForm.Get("groupid")

	if groupid == "" {
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "empty group assignment: unspecified group")))
		return
	}
	account, err := o.fetchAccountByUsername(r.Context(), userid)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		return
	}

	// ocs only knows about names so we have to look up the internal id
	group, err := o.fetchGroupByName(r.Context(), groupid)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		return
	}

	_, err = o.getGroupsService().AddMember(r.Context(), &accounts.AddMemberRequest{
		AccountId: account.Id,
		GroupId:   group.Id,
	})

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("userid", account.Id).Str("groupid", group.Id).Msg("could not add user to group")
		return
	}

	o.logger.Debug().Str("userid", account.Id).Str("groupid", group.Id).Msg("added user to group")
	mustNotFail(render.Render(w, r, response.DataRender(struct{}{})))
}

// RemoveFromGroup removes a user from a group
func (o Ocs) RemoveFromGroup(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")

	var err error

	// Really? a DELETE with form encoded body?!?
	// but it is not encoded as mime, so we cannot just call r.ParseForm()
	// read it manually
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, err.Error())))
		return
	}
	if err = r.Body.Close(); err != nil {
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		return
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, err.Error())))
		return
	}

	groupid := values.Get("groupid")
	if groupid == "" {
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "no group id")))
		return
	}

	var account *accounts.Account

	if isValidUUID(userid) {
		account, _ = o.getAccountService().GetAccount(r.Context(), &accounts.GetAccountRequest{
			Id: userid,
		})
	} else {
		// despite the confusion, if we make it here we got ourselves a username
		account, err = o.fetchAccountByUsername(r.Context(), userid)
		if err != nil {
			merr := merrors.FromError(err)
			if merr.Code == http.StatusNotFound {
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "The requested user could not be found")))
			} else {
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
			}
			o.logger.Error().Err(err).Str("userid", userid).Msg("could not get list of user groups")
			return
		}
	}

	// ocs only knows about names so we have to look up the internal id
	group, err := o.fetchGroupByName(r.Context(), groupid)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		return
	}

	_, err = o.getGroupsService().RemoveMember(r.Context(), &accounts.RemoveMemberRequest{
		AccountId: account.Id,
		GroupId:   group.Id,
	})

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("userid", account.Id).Str("groupid", group.Id).Msg("could not remove user from group")
		return
	}

	o.logger.Debug().Str("userid", account.Id).Str("groupid", group.Id).Msg("removed user from group")
	mustNotFail(render.Render(w, r, response.DataRender(struct{}{})))
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
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not list users")))
		return
	}

	groups := make([]string, 0, len(res.Groups))
	for i := range res.Groups {
		groups = append(groups, res.Groups[i].OnPremisesSamAccountName)
	}

	mustNotFail(render.Render(w, r, response.DataRender(&data.Groups{Groups: groups})))
}

// AddGroup adds a group
func (o Ocs) AddGroup(w http.ResponseWriter, r *http.Request) {
	groupid := r.PostFormValue("groupid")
	displayname := r.PostFormValue("displayname")
	gid := r.PostFormValue("gidnumber")

	var gidNumber int64
	var err error

	if gid != "" {
		gidNumber, err = strconv.ParseInt(gid, 10, 64)
		if err != nil {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "Cannot use the gidnumber provided")))
			o.logger.Error().Err(err).Str("gid", gid).Str("groupid", groupid).Msg("Cannot use the gidnumber provided")
			return
		}
	}

	if displayname == "" {
		displayname = groupid
	}

	newGroup := &accounts.Group{
		Id:                       groupid,
		DisplayName:              displayname,
		OnPremisesSamAccountName: groupid,
		GidNumber:                gidNumber,
	}
	group, err := o.getGroupsService().CreateGroup(r.Context(), &accounts.CreateGroupRequest{
		Group: newGroup,
	})
	if err != nil {
		merr := merrors.FromError(err)
		switch merr.Code {
		case http.StatusBadRequest:
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, merr.Detail)))
		case http.StatusConflict:
			if response.APIVersion(r.Context()) == "2" {
				// it seems the application framework sets the ocs status code to the httpstatus code, which affects the provisioning api
				// see https://github.com/owncloud/core/blob/b9ff4c93e051c94adfb301545098ae627e52ef76/lib/public/AppFramework/OCSController.php#L142-L150
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, merr.Detail)))
			} else {
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaInvalidInput.StatusCode, merr.Detail)))
			}
		default:
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("groupid", groupid).Msg("could not add group")
		// TODO check error if group already existed
		return
	}
	o.logger.Debug().Interface("group", group).Msg("added group")

	mustNotFail(render.Render(w, r, response.DataRender(struct{}{})))
}

// DeleteGroup deletes a group
func (o Ocs) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	groupid := chi.URLParam(r, "groupid")

	// ocs only knows about names so we have to look up the internal id
	group, err := o.fetchGroupByName(r.Context(), groupid)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		return
	}

	_, err = o.getGroupsService().DeleteGroup(r.Context(), &accounts.DeleteGroupRequest{
		Id: group.Id,
	})

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("groupid", group.Id).Msg("could not remove group")
		return
	}

	o.logger.Debug().Str("groupid", group.Id).Msg("removed group")
	mustNotFail(render.Render(w, r, response.DataRender(struct{}{})))
}

// GetGroupMembers lists all members of a group
func (o Ocs) GetGroupMembers(w http.ResponseWriter, r *http.Request) {

	groupid := chi.URLParam(r, "groupid")

	// ocs only knows about names so we have to look up the internal id
	group, err := o.fetchGroupByName(r.Context(), groupid)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		return
	}

	res, err := o.getGroupsService().ListMembers(r.Context(), &accounts.ListMembersRequest{Id: group.Id})

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested group could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("groupid", group.Id).Msg("could not get list of members")
		return
	}

	members := make([]string, 0, len(res.Members))
	for i := range res.Members {
		members = append(members, res.Members[i].OnPremisesSamAccountName)
	}

	o.logger.Error().Err(err).Int("count", len(members)).Str("groupid", groupid).Msg("listing group members")
	mustNotFail(render.Render(w, r, response.DataRender(&data.Users{Users: members})))
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func (o Ocs) fetchGroupByName(ctx context.Context, name string) (*accounts.Group, error) {
	var res *accounts.ListGroupsResponse
	res, err := o.getGroupsService().ListGroups(ctx, &accounts.ListGroupsRequest{
		Query: fmt.Sprintf("on_premises_sam_account_name eq '%v'", escapeValue(name)),
	})
	if err != nil {
		return nil, err
	}
	if res != nil && len(res.Groups) == 1 {
		return res.Groups[0], nil
	}
	return nil, merrors.NotFound("", "The requested group could not be found")
}

func mustNotFail(err error) {
	if err != nil {
		panic(err)
	}
}
