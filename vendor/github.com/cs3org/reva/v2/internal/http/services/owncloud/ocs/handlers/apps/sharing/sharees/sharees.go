// Copyright 2018-2021 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

package sharees

import (
	"net/http"

	grouppb "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"

	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/config"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/conversions"
	"github.com/cs3org/reva/v2/internal/http/services/owncloud/ocs/response"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storage/utils/templates"
)

// Handler implements the ownCloud sharing API
type Handler struct {
	gatewayAddr             string
	additionalInfoAttribute string
}

// Init initializes this and any contained handlers
func (h *Handler) Init(c *config.Config) {
	h.gatewayAddr = c.GatewaySvc
	h.additionalInfoAttribute = c.AdditionalInfoAttribute
}

// FindSharees implements the /apps/files_sharing/api/v1/sharees endpoint
func (h *Handler) FindSharees(w http.ResponseWriter, r *http.Request) {
	log := appctx.GetLogger(r.Context())
	term := r.URL.Query().Get("search")

	if term == "" {
		response.WriteOCSError(w, r, response.MetaBadRequest.StatusCode, "search must not be empty", nil)
		return
	}

	selector, err := pool.GatewaySelector(h.gatewayAddr)
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error getting gateway selector", err)
		return
	}
	gwc, err := selector.Next()
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error selecting next client", err)
		return
	}
	usersRes, err := gwc.FindUsers(r.Context(), &userpb.FindUsersRequest{Filter: term, SkipFetchingUserGroups: true})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error searching users", err)
		return
	}
	log.Debug().Int("count", len(usersRes.GetUsers())).Str("search", term).Msg("users found")

	userMatches := make([]*conversions.MatchData, 0, len(usersRes.GetUsers()))
	for _, user := range usersRes.GetUsers() {
		match := h.userAsMatch(user)
		log.Debug().Interface("user", user).Interface("match", match).Msg("mapped")
		userMatches = append(userMatches, match)
	}

	groupsRes, err := gwc.FindGroups(r.Context(), &grouppb.FindGroupsRequest{Filter: term, SkipFetchingMembers: true})
	if err != nil {
		response.WriteOCSError(w, r, response.MetaServerError.StatusCode, "error searching groups", err)
		return
	}
	log.Debug().Int("count", len(groupsRes.GetGroups())).Str("search", term).Msg("groups found")

	groupMatches := make([]*conversions.MatchData, 0, len(groupsRes.GetGroups()))
	for _, g := range groupsRes.GetGroups() {
		match := h.groupAsMatch(g)
		log.Debug().Interface("group", g).Interface("match", match).Msg("mapped")
		groupMatches = append(groupMatches, match)
	}

	response.WriteOCSSuccess(w, r, &conversions.ShareeData{
		Exact: &conversions.ExactMatchesData{
			Users:   []*conversions.MatchData{},
			Groups:  []*conversions.MatchData{},
			Remotes: []*conversions.MatchData{},
		},
		Users:   userMatches,
		Groups:  groupMatches,
		Remotes: []*conversions.MatchData{},
	})
}

func (h *Handler) userAsMatch(u *userpb.User) *conversions.MatchData {
	var ocsUserType int
	if u.Id.Type == userpb.UserType_USER_TYPE_GUEST || u.Id.Type == userpb.UserType_USER_TYPE_LIGHTWEIGHT {
		ocsUserType = 1
	}

	return &conversions.MatchData{
		Label: u.DisplayName,
		Value: &conversions.MatchValueData{
			ShareType: int(conversions.ShareTypeUser),
			// api compatibility with oc10: mark guest users in share invite dialogue
			UserType: ocsUserType,
			// api compatibility with oc10: always use the username
			ShareWith:               u.Username,
			ShareWithAdditionalInfo: h.getAdditionalInfoAttribute(u),
		},
	}
}

func (h *Handler) groupAsMatch(g *grouppb.Group) *conversions.MatchData {
	return &conversions.MatchData{
		Label: g.DisplayName,
		Value: &conversions.MatchValueData{
			ShareType:               int(conversions.ShareTypeGroup),
			ShareWith:               g.GroupName,
			ShareWithAdditionalInfo: g.Mail,
		},
	}
}

func (h *Handler) getAdditionalInfoAttribute(u *userpb.User) string {
	return templates.WithUser(u, h.additionalInfoAttribute)
}
