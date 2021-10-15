package svc

import (
	"context"
	"net/http"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/group/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	//msgraph "github.com/owncloud/open-graph-api-go" // FIXME add groups to open graph, needs OnPremisesSamAccountName and OnPremisesDomainName
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

// GroupCtx middleware is used to load an User object from
// the URL parameters passed through as the request. In case
// the User could not be found, we stop here and return a 404.
func (g Graph) GroupCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupID := chi.URLParam(r, "groupID")
		if groupID == "" {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing group id")
			return
		}

		client, err := g.GetClient()
		if err != nil {
			g.logger.Error().Err(err).Msg("could not get client")
			errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := client.GetGroupByClaim(r.Context(), &cs3.GetGroupByClaimRequest{
			Claim: "groupid", // FIXME add consts to reva
			Value: groupID,
		})

		switch {
		case err != nil:
			g.logger.Error().Err(err).Str("groupid", groupID).Msg("error sending get group by claim id grpc request")
			errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
			return
		case res.Status.Code != cs3rpc.Code_CODE_OK:
			if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
				errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message)
				return
			}
			g.logger.Error().Err(err).Str("groupid", groupID).Msg("error sending get group by claim id grpc request")
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
			return
		}

		ctx := context.WithValue(r.Context(), groupKey, res.Group)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetGroups implements the Service interface.
func (g Graph) GetGroups(w http.ResponseWriter, r *http.Request) {
	client, err := g.GetClient()
	if err != nil {
		g.logger.Error().Err(err).Msg("could not get client")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	search := r.URL.Query().Get("search")
	if search == "" {
		search = r.URL.Query().Get("$search")
	}

	res, err := client.FindGroups(r.Context(), &cs3.FindGroupsRequest{
		// FIXME presence match is currently not implemented, an empty search currently leads to
		// Unwilling To Perform": Search Error: error parsing filter: (&(objectclass=posixAccount)(|(cn=*)(displayname=*)(mail=*))), error: Present filter match for cn not implemented
		Filter: search,
	})
	switch {
	case err != nil:
		g.logger.Error().Err(err).Str("search", search).Msg("error sending find groups grpc request")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message)
			return
		}
		g.logger.Error().Err(err).Str("search", search).Msg("error sending find groups grpc request")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	groups := make([]*msgraph.Group, 0, len(res.Groups))

	for _, group := range res.Groups {
		groups = append(groups, createGroupModelFromCS3(group))
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: groups})
}

// GetGroup implements the Service interface.
func (g Graph) GetGroup(w http.ResponseWriter, r *http.Request) {
	group := r.Context().Value(groupKey).(*cs3.Group)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, createGroupModelFromCS3(group))
}

func createGroupModelFromCS3(g *cs3.Group) *msgraph.Group {
	if g.Id == nil {
		g.Id = &cs3.GroupId{}
	}
	return &msgraph.Group{
		DirectoryObject: msgraph.DirectoryObject{
			Entity: msgraph.Entity{
				ID: &g.Id.OpaqueId,
			},
		},
		OnPremisesDomainName:     &g.Id.Idp,
		OnPremisesSamAccountName: &g.GroupName,
		DisplayName:              &g.DisplayName,
		Mail:                     &g.Mail,
		// TODO when to fetch and expand memberof, usernames or ids?
	}
}
