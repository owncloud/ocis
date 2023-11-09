package svc

import (
	"github.com/CiscoM31/godata"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/v2/services/graph/pkg/service/v0/errorcode"
	"net/http"
	"strings"
)

// GetOwnLanguage returns the language of the current user.
func (g Graph) GetOwnLanguage(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	g.logger.Debug().Msg("Calling GetOwnLanguage")
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")

	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		logger.Debug().Msg("could not get user: user not in context")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "user not in context")
		return
	}

	me, err := g.identityBackend.GetUser(r.Context(), u.GetId().GetOpaqueId(), odataReq)
	if err != nil {
		logger.Debug().Err(err).Interface("user", u).Msg("could not get user from backend")
		errorcode.RenderError(w, r, err)
		return
	}

	// TODO: make sure that this actually returns the stored language
	lang, ok := me.GetPreferredLanguageOk()
	if !ok {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, nil)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, lang)
}

// SetOwnLanguage sets the language of the current user.
func (g Graph) SetOwnLanguage(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	g.logger.Debug().Msg("Calling SetOwnLanguage")
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")

	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		logger.Debug().Msg("could not get user: user not in context")
		errorcode.GeneralException.Render(w, r, http.StatusInternalServerError, "user not in context")
		return
	}

	me, err := g.identityBackend.GetUser(r.Context(), u.GetId().GetOpaqueId(), odataReq)
	if err != nil {
		logger.Debug().Err(err).Interface("user", u).Msg("could not get user from backend")
		errorcode.RenderError(w, r, err)
		return
	}

	lang := chi.URLParam(r, "language")
	me.SetPreferredLanguage(lang)
	// TODO: persist this change
	render.Status(r, http.StatusNoContent)
}

func (g Graph) SetUserLanguage(w http.ResponseWriter, r *http.Request) {
	logger := g.logger.SubloggerWithRequestID(r.Context())
	g.logger.Debug().Msg("Calling SetUserLanguage")
	sanitizedPath := strings.TrimPrefix(r.URL.Path, "/graph/v1.0/")

	odataReq, err := godata.ParseRequest(r.Context(), sanitizedPath, r.URL.Query())
	if err != nil {
		logger.Debug().Err(err).Interface("query", r.URL.Query()).Msg("could not get users: query error")
		errorcode.InvalidRequest.Render(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := g.identityBackend.GetUser(r.Context(), chi.URLParam(r, "userID"), odataReq)
	if err != nil {
		logger.Debug().Err(err).Interface("user", user.GetId()).Msg("could not get user from backend")
		errorcode.RenderError(w, r, err)
		return
	}

	lang := chi.URLParam(r, "language")
	user.SetPreferredLanguage(lang)
	// TODO: persist this change
	render.Status(r, http.StatusNoContent)
}
