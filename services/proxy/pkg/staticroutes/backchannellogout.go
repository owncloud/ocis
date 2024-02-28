package staticroutes

import (
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	microstore "go-micro.dev/v4/store"
	"net/http"
)

// handle backchannel logout requests as per https://openid.net/specs/openid-connect-backchannel-1_0.html#BCRequest
func (s *StaticRouteHandler) backchannelLogout(w http.ResponseWriter, r *http.Request) {
	// parse the application/x-www-form-urlencoded POST request
	logger := s.Logger.SubloggerWithRequestID(r.Context())
	if err := r.ParseForm(); err != nil {
		logger.Warn().Err(err).Msg("ParseForm failed")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, jse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	logoutToken, err := s.OidcClient.VerifyLogoutToken(r.Context(), r.PostFormValue("logout_token"))
	if err != nil {
		logger.Warn().Err(err).Msg("VerifyLogoutToken failed")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, jse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	records, err := s.UserInfoCache.Read(logoutToken.SessionId)
	if errors.Is(err, microstore.ErrNotFound) || len(records) == 0 {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, nil)
		return
	}

	if err != nil {
		logger.Error().Err(err).Msg("Error reading userinfo cache")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, jse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	for _, record := range records {
		err = s.UserInfoCache.Delete(string(record.Value))
		if err != nil && !errors.Is(err, microstore.ErrNotFound) {
			// Spec requires us to return a 400 BadRequest when the session could not be destroyed
			logger.Err(err).Msg("could not delete user info from cache")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jse{Error: "invalid_request", ErrorDescription: err.Error()})
			return
		}
		logger.Debug().Msg("Deleted userinfo from cache")
	}

	// we can ignore errors when cleaning up the lookup table
	err = s.UserInfoCache.Delete(logoutToken.SessionId)
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to cleanup sessionid lookup entry")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)
}
