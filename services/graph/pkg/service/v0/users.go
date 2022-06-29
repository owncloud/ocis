package svc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/CiscoM31/godata"
	storageprovider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	settings "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	settingssvc "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
	"golang.org/x/exp/slices"
)

// GetMe implements the Service interface.
func (g Graph) GetMe(w http.ResponseWriter, r *http.Request) {

	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		g.logger.Error().Msg("user not in context")
		errorcode.ServiceNotAvailable.Render(w, r, http.StatusInternalServerError, "user not in context")
		return
	}

	g.logger.Info().Interface("user", u).Msg("User in /me")

	me := identity.CreateUserModelFromCS3(u)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, me)
}

// GetUsers implements the Service interface.
func (g Graph) GetUsers(w http.ResponseWriter, r *http.Request) {
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")
	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		g.logger.Err(err).Interface("query", r.URL.Query()).Msg("query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}
	users, err := g.identityBackend.GetUsers(r.Context(), r.URL.Query())
	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	users, err = sortUsers(odataReq, users)
	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &listResponse{Value: users})
}

func (g Graph) PostUser(w http.ResponseWriter, r *http.Request) {
	u := libregraph.NewUser()
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if _, ok := u.GetDisplayNameOk(); !ok {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing required Attribute: 'displayName'")
		return
	}
	if accountName, ok := u.GetOnPremisesSamAccountNameOk(); ok {
		if !isValidUsername(*accountName) {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest,
				fmt.Sprintf("username '%s' must be at least the local part of an email", *u.OnPremisesSamAccountName))
			return
		}
	} else {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing required Attribute: 'onPremisesSamAccountName'")
		return
	}

	if mail, ok := u.GetMailOk(); ok {
		if !isValidEmail(*mail) {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest,
				fmt.Sprintf("'%s' is not a valid email address", *u.Mail))
			return
		}
	} else {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing required Attribute: 'mail'")
		return
	}

	// Disallow user-supplied IDs. It's supposed to be readonly. We're either
	// generating them in the backend ourselves or rely on the Backend's
	// storage (e.g. LDAP) to provide a unique ID.
	if _, ok := u.GetIdOk(); ok {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "user id is a read-only attribute")
		return
	}

	if u, err = g.identityBackend.CreateUser(r.Context(), *u); err != nil {
		var ecErr errorcode.Error
		if errors.As(err, &ecErr) {
			ecErr.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// All users get the user role by default currently.
	// to all new users for now, as create Account request does not have any role field
	if g.roleService == nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "could not assign role to account: roleService not configured")
		return
	}
	if _, err = g.roleService.AssignRoleToUser(r.Context(), &settings.AssignRoleToUserRequest{
		AccountUuid: *u.Id,
		RoleId:      settingssvc.BundleUUIDRoleUser,
	}); err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, fmt.Sprintf("could not assign role to account %s", err.Error()))
		return
	}

	currentUser := ctxpkg.ContextMustGetUser(r.Context())
	g.publishEvent(events.UserCreated{Executant: currentUser.Id, UserID: *u.Id})

	render.Status(r, http.StatusOK)
	render.JSON(w, r, u)
}

// GetUser implements the Service interface.
func (g Graph) GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	userID, err := url.PathUnescape(userID)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping user id failed")
	}

	if userID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing user id")
		return
	}

	user, err := g.identityBackend.GetUser(r.Context(), userID, r.URL.Query())

	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
	}
	sel := strings.Split(r.URL.Query().Get("$select"), ",")
	exp := strings.Split(r.URL.Query().Get("$expand"), ",")
	if slices.Contains(sel, "drive") || slices.Contains(sel, "drives") || slices.Contains(exp, "drive") || slices.Contains(exp, "drives") {
		wdu, err := url.Parse(g.config.Spaces.WebDavBase + g.config.Spaces.WebDavPath)
		// TODO: filter this request for the user
		f := listStorageSpacesUserFilter(user.GetId())
		lspr, err := g.gatewayClient.ListStorageSpaces(r.Context(), &storageprovider.ListStorageSpacesRequest{
			Opaque:  nil,
			Filters: []*storageprovider.ListStorageSpacesRequest_Filter{f},
		})
		if err != nil {
			g.logger.Err(err).Interface("query", r.URL.Query()).Msg("error getting storages")
		}
		drives := []libregraph.Drive{}
		for _, sp := range lspr.GetStorageSpaces() {
			d, err := g.cs3StorageSpaceToDrive(r.Context(), wdu, sp)
			if err != nil {
				g.logger.Err(err).Interface("query", r.URL.Query()).Msg("error converting space to drive")
			}
			if slices.Contains(sel, "drive") || slices.Contains(exp, "drive") {
				if *d.DriveType == "personal" {
					drives = append(drives, *d)
				}
			} else {
				drives = append(drives, *d)
			}
		}
		user.Drives = drives
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

func (g Graph) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	userID, err := url.PathUnescape(userID)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping user id failed")
	}

	if userID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing user id")
		return
	}

	err = g.identityBackend.DeleteUser(r.Context(), userID)

	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
	}

	currentUser := ctxpkg.ContextMustGetUser(r.Context())
	g.publishEvent(events.UserDeleted{Executant: currentUser.Id, UserID: userID})

	render.Status(r, http.StatusNoContent)
	render.NoContent(w, r)
}

// PatchUser implements the Service Interface. Updates the specified attributes of an
// ExistingUser
func (g Graph) PatchUser(w http.ResponseWriter, r *http.Request) {
	nameOrID := chi.URLParam(r, "userID")
	nameOrID, err := url.PathUnescape(nameOrID)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping user id failed")
	}

	if nameOrID == "" {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "missing user id")
		return
	}
	changes := libregraph.NewUser()
	err = json.NewDecoder(r.Body).Decode(changes)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var features []events.UserFeature
	if mail, ok := changes.GetMailOk(); ok {
		if !isValidEmail(*mail) {
			errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest,
				fmt.Sprintf("'%s' is not a valid email address", *mail))
			return
		}
		features = append(features, events.UserFeature{Name: "email", Value: *mail})
	}

	if name, ok := changes.GetDisplayNameOk(); ok {
		features = append(features, events.UserFeature{Name: "displayname", Value: *name})
	}

	u, err := g.identityBackend.UpdateUser(r.Context(), nameOrID, *changes)
	if err != nil {
		var errcode errorcode.Error
		if errors.As(err, &errcode) {
			errcode.Render(w, r)
		} else {
			errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		}
	}

	currentUser := ctxpkg.ContextMustGetUser(r.Context())
	g.publishEvent(
		events.UserFeatureChanged{
			Executant: currentUser.Id,
			UserID:    nameOrID,
			Features:  features,
		},
	)
	render.Status(r, http.StatusOK)
	render.JSON(w, r, u)

}

// We want to allow email addresses as usernames so they show up when using them in ACLs on storages that allow integration with our glauth LDAP service
// so we are adding a few restrictions from https://stackoverflow.com/questions/6949667/what-are-the-real-rules-for-linux-usernames-on-centos-6-and-rhel-6
// names should not start with numbers
var usernameRegex = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]*(@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)*$")

func isValidUsername(e string) bool {
	if len(e) < 1 && len(e) > 254 {
		return false
	}
	return usernameRegex.MatchString(e)
}

// regex from https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html#valid-e-mail-address
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValidEmail(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func sortUsers(req *godata.GoDataRequest, users []*libregraph.User) ([]*libregraph.User, error) {
	var sorter sort.Interface
	if req.Query.OrderBy == nil || len(req.Query.OrderBy.OrderByItems) != 1 {
		return users, nil
	}
	switch req.Query.OrderBy.OrderByItems[0].Field.Value {
	case "displayName":
		sorter = usersByDisplayName{users}
	case "mail":
		sorter = usersByMail{users}
	case "onPremisesSamAccountName":
		sorter = usersByOnPremisesSamAccountName{users}
	default:
		return nil, fmt.Errorf("we do not support <%s> as a order parameter", req.Query.OrderBy.OrderByItems[0].Field.Value)
	}

	if req.Query.OrderBy.OrderByItems[0].Order == "desc" {
		sorter = sort.Reverse(sorter)
	}
	sort.Sort(sorter)
	return users, nil
}
