package svc

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/owncloud/ocis/v2/ocis-pkg/l10n"
	"github.com/owncloud/ocis/v2/services/graph/pkg/errorcode"
	l10n_pkg "github.com/owncloud/ocis/v2/services/graph/pkg/l10n"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
)

// GetRoleDefinitions a list of permission roles than can be used when sharing with users or groups
func (g Graph) GetRoleDefinitions(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	roles := unifiedrole.GetRoles(unifiedrole.RoleFilterIDs(g.config.UnifiedRoles.AvailableRoles...))

	userID := revactx.ContextMustGetUser(r.Context()).GetId().GetOpaqueId()
	locale := l10n.MustGetUserLocale(r.Context(), userID, r.Header.Get(l10n.HeaderAcceptLanguage), g.valueService)
	if locale != "" && locale != "en" {
		err := l10n_pkg.TranslateEntity(locale, "en", roles,
			l10n.TranslateField("DisplayName"),
			l10n.TranslateField("Description"),
		)
		if err != nil {
			logger.Error().Err(err).Msg("could not translate role definitions")
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, roles)
}

// GetRoleDefinition a permission role than can be used when sharing with users or groups
func (g Graph) GetRoleDefinition(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	roleID, err := url.PathUnescape(chi.URLParam(r, "roleID"))
	if err != nil {
		logger.Debug().Err(err).Str("roleID", chi.URLParam(r, "roleID")).Msg("could not get roleID: unescaping is failed")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, "unescaping role id failed")
		return
	}
	role, err := unifiedrole.GetRole(unifiedrole.RoleFilterIDs(roleID))
	if err != nil {
		logger.Debug().Str("roleID", roleID).Msg("could not get role: not found")
		errorcode.ItemNotFound.Render(w, r, http.StatusNotFound, err.Error())
		return
	}

	userID := revactx.ContextMustGetUser(r.Context()).GetId().GetOpaqueId()
	locale := l10n.MustGetUserLocale(r.Context(), userID, r.Header.Get(l10n.HeaderAcceptLanguage), g.valueService)
	if locale != "" && locale != "en" {
		err := l10n_pkg.TranslateEntity(locale, "en", role,
			l10n.TranslateField("DisplayName"),
			l10n.TranslateField("Description"),
		)
		if err != nil {
			logger.Error().Str("roleID", roleID).Err(err).Msg("could not translate role definition")
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, role)
}
