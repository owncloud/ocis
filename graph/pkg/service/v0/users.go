package svc

import (
	"context"
	"net/http"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	revactx "github.com/cs3org/reva/pkg/ctx"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/graph/pkg/service/v0/errorcode"

	//msgraph "github.com/owncloud/open-graph-api-go" // FIXME needs OnPremisesSamAccountName, OnPremisesDomainName and AdditionalData
	msgraph "github.com/yaegashi/msgraph.go/v1.0"
)

// UserCtx middleware is used to load an User object from
// the URL parameters passed through as the request. In case
// the User could not be found, we stop here and return a 404.
// TODO use cs3 api to look up user
func (g Graph) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userID := chi.URLParam(r, "userID")
		if userID == "" {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing user id")
			return
		}

		client, err := g.GetClient()
		if err != nil {
			g.logger.Error().Err(err).Msg("could not get client")
			errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		res, err := client.GetUserByClaim(r.Context(), &cs3.GetUserByClaimRequest{
			Claim: "userid", // FIXME add consts to reva
			Value: userID,
		})

		switch {
		case err != nil:
			g.logger.Error().Err(err).Str("userid", userID).Msg("error sending get user by claim id grpc request")
			errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
			return
		case res.Status.Code != cs3rpc.Code_CODE_OK:
			if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
				errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message)
				return
			}
			g.logger.Error().Err(err).Str("userid", userID).Msg("error sending get user by claim id grpc request")
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, res.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetMe implements the Service interface.
func (g Graph) GetMe(w http.ResponseWriter, r *http.Request) {

	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		g.logger.Error().Msg("user not in context")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, "user not in context")
		return
	}

	g.logger.Info().Interface("user", u).Msg("User in /me")

	me := createUserModelFromCS3(u)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, me)
}

// GetUsers implements the Service interface.
// TODO use cs3 api to look up user
func (g Graph) GetUsers(w http.ResponseWriter, r *http.Request) {

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

	res, err := client.FindUsers(r.Context(), &cs3.FindUsersRequest{
		// FIXME presence match is currently not implemented, an empty search currently leads to
		// Unwilling To Perform": Search Error: error parsing filter: (&(objectclass=posixAccount)(|(cn=*)(displayname=*)(mail=*))), error: Present filter match for cn not implemented
		Filter: search,
	})
	switch {
	case err != nil:
		g.logger.Error().Err(err).Str("search", search).Msg("error sending find users grpc request")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	case res.Status.Code != cs3rpc.Code_CODE_OK:
		if res.Status.Code == cs3rpc.Code_CODE_NOT_FOUND {
			errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, res.Status.Message)
			return
		}
		g.logger.Error().Err(err).Str("search", search).Msg("error sending find users grpc request")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, res.Status.Message)
		return
	}

	users := make([]*msgraph.User, 0, len(res.Users))

	for _, user := range res.Users {
		users = append(users, createUserModelFromCS3(user))
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: users})
}

// GetUser implements the Service interface.
func (g Graph) GetUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(*cs3.User)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, createUserModelFromCS3(user))
}

func createUserModelFromCS3(u *cs3.User) *msgraph.User {
	if u.Id == nil {
		u.Id = &cs3.UserId{}
	}
	return &msgraph.User{
		DisplayName: &u.DisplayName,
		Mail:        &u.Mail,
		// TODO u.Groups are those ids or group names?
		OnPremisesSamAccountName: &u.Username,
		DirectoryObject: msgraph.DirectoryObject{
			Entity: msgraph.Entity{
				ID: &u.Id.OpaqueId,
				Object: msgraph.Object{
					AdditionalData: map[string]interface{}{
						"uidnumber": u.UidNumber,
						"gidnumber": u.GidNumber,
					},
				},
			},
		},
	}
}
