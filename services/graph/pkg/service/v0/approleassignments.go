package svc

import (
	"context"
	"fmt"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	libregraph "github.com/owncloud/libre-graph-api-go"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	settingsmsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/settings/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/utils"
	merrors "go-micro.dev/v4/errors"
)

const principalTypeUser = "User"

// ListAppRoleAssignments implements the Service interface.
func (g Graph) ListAppRoleAssignments(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling list appRoleAssignments")

	userID := chi.URLParam(r, "userID")

	lrar, err := g.roleService.ListRoleAssignments(r.Context(), &settingssvc.ListRoleAssignmentsRequest{
		AccountUuid: userID,
	})
	if err != nil {
		// TODO check the error type and return proper error code
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	values := make([]libregraph.AppRoleAssignment, 0, len(lrar.GetAssignments()))
	for _, assignment := range lrar.GetAssignments() {
		values = append(values, g.assignmentToAppRoleAssignment(assignment))
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &ListResponse{Value: values})
}

// CreateAppRoleAssignment implements the Service interface.
func (g Graph) CreateAppRoleAssignment(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("query", r.URL.Query()).Msg("calling create appRoleAssignment")

	appRoleAssignment := libregraph.NewAppRoleAssignmentWithDefaults()
	err := StrictJSONUnmarshal(r.Body, appRoleAssignment)
	if err != nil {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err.Error()))
		return
	}

	userID := chi.URLParam(r, "userID")

	// We can ignore the error, in the worst case the old role will be empty
	oldRoles, _ := g.roleService.ListRoleAssignments(r.Context(), &settingssvc.ListRoleAssignmentsRequest{
		AccountUuid: userID,
	})

	if appRoleAssignment.GetPrincipalId() != userID {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("user id %s does not match principal id %v", userID, appRoleAssignment.GetPrincipalId()))
		return
	}
	if appRoleAssignment.GetResourceId() != g.config.Application.ID {
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, fmt.Sprintf("resource id %s does not match expected application id %v", userID, g.config.Application.ID))
		return
	}

	artur, err := g.roleService.AssignRoleToUser(r.Context(), &settingssvc.AssignRoleToUserRequest{
		AccountUuid: userID,
		RoleId:      appRoleAssignment.GetAppRoleId(),
	})
	if err != nil {
		if merr, ok := merrors.As(err); ok && merr.Code == http.StatusForbidden {
			errorcode.NotAllowed.Render(w, r, http.StatusForbidden, err.Error())
			return
		}
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	var oldRole string
	roles := oldRoles.GetAssignments()
	if len(roles) > 0 {
		oldRole = roles[0].GetRoleId()
	}

	client, err := g.gatewaySelector.Next()
	if err != nil {
		errorcode.NotAllowed.Render(w, r, http.StatusForbidden, err.Error())
		return
	}

	canCreateDrives := g.checkPermission(r.Context(), "Drives.Create", userID, client)
	if canCreateDrives {
		err = shared.RestorePersonalSpace(r.Context(), client, userID)
		if err != nil {
			logger.Error().Any("userID", userID).Err(err).Msg("can't ensure the personal space")
			errorcode.RenderError(w, r, err)
			return
		}
	} else {
		err := shared.DisablePersonalSpace(r.Context(), client, userID)
		if err != nil {
			logger.Error().Any("userID", userID).Err(err).Msg("can't disable the personal space")
			errorcode.RenderError(w, r, err)
			return
		}
	}

	e := events.UserFeatureChanged{
		UserID: userID,
		Features: []events.UserFeature{
			{
				Name:     "roleChanged",
				Value:    appRoleAssignment.AppRoleId,
				OldValue: &oldRole,
			},
		},
		Timestamp: utils.TSNow(),
	}
	if currentUser, ok := revactx.ContextGetUser(r.Context()); ok {
		e.Executant = currentUser.GetId()
	}
	g.publishEvent(r.Context(), e)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, g.assignmentToAppRoleAssignment(artur.GetAssignment()))
}

// DeleteAppRoleAssignment implements the Service interface.
func (g Graph) DeleteAppRoleAssignment(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	logger.Info().Interface("body", r.Body).Msg("calling delete appRoleAssignment")

	userID := chi.URLParam(r, "userID")

	// check assignment belongs to the user
	lrar, err := g.roleService.ListRoleAssignments(r.Context(), &settingssvc.ListRoleAssignmentsRequest{
		AccountUuid: userID,
	})
	if err != nil {
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	appRoleAssignmentID := chi.URLParam(r, "appRoleAssignmentID")

	assignmentFound := false
	for _, roleAssignment := range lrar.GetAssignments() {
		if roleAssignment.Id == appRoleAssignmentID {
			assignmentFound = true
		}
	}
	if !assignmentFound {
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, fmt.Sprintf("appRoleAssignment %v not found for user %v", appRoleAssignmentID, userID))
		return
	}

	_, err = g.roleService.RemoveRoleFromUser(r.Context(), &settingssvc.RemoveRoleFromUserRequest{
		Id: appRoleAssignmentID,
	})
	if err != nil {
		if merr, ok := merrors.As(err); ok && merr.Code == http.StatusForbidden {
			errorcode.NotAllowed.Render(w, r, http.StatusForbidden, err.Error())
			return
		}
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.NoContent(w, r)
}

func (g Graph) assignmentToAppRoleAssignment(assignment *settingsmsg.UserRoleAssignment) libregraph.AppRoleAssignment {
	appRoleAssignment := libregraph.NewAppRoleAssignmentWithDefaults()
	appRoleAssignment.SetId(assignment.Id)
	appRoleAssignment.SetAppRoleId(assignment.RoleId)
	appRoleAssignment.SetPrincipalType(principalTypeUser) // currently always assigned to the user
	appRoleAssignment.SetResourceId(g.config.Application.ID)
	appRoleAssignment.SetResourceDisplayName(g.config.Application.DisplayName)
	appRoleAssignment.SetPrincipalId(assignment.AccountUuid)
	// appRoleAssignment.SetPrincipalDisplayName() // TODO fetch and cache
	return *appRoleAssignment
}

func (g Graph) checkPermission(ctx context.Context, perm string, userID string, gwc gateway.GatewayAPIClient) bool {
	u, err := utils.GetUserWithContext(ctx, &userv1beta1.UserId{OpaqueId: userID}, gwc)
	if err != nil {
		g.logger.Error().Err(err).Msg("could not get user")
		return false
	}

	if ok, err := utils.CheckPermission(revactx.ContextSetUser(context.Background(), u), perm, gwc); ok {
		return true
	} else if err != nil {
		g.logger.Error().Err(err).Msg("error checking permission")
	}
	return false
}
