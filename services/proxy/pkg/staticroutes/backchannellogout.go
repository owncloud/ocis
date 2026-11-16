package staticroutes

import (
	"context"
	"fmt"
	"net/http"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/go-chi/render"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/utils"
	"github.com/pkg/errors"
	"github.com/shamaton/msgpack/v2"
	microstore "go-micro.dev/v4/store"
)

// mfaRecordKey returns the store key for a user's per-user MFA status
// record. This naming is a contract with the proxy MFA middleware's key()
// (services/proxy/pkg/middleware/mfa.go), which writes and reads the
// record under this same name. Keep both sides in sync.
func mfaRecordKey(userID string) string {
	return "mfa:" + userID
}

// mfaSigningKeyName returns the store key for a user's MFA-bound signing
// key. This naming is a contract with the ocs service's GetSigningKey
// handler (services/ocs/pkg/service/v0/users.go) and the proxy's
// signed-URL authenticator's mfaSigningKeyName()
// (services/proxy/pkg/middleware/signed_url_auth.go), which write/read the
// key under this same name. Keep all three in sync.
func mfaSigningKeyName(userID string) string {
	return userID + ":mfa"
}

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
		user, resolveErr := s.resolveLogoutUser(r.Context(), record)
		if resolveErr != nil {
			s.Logger.Warn().Err(resolveErr).Msg("could not resolve user for backchannel logout")
		} else {
			if err := s.publishBackchannelLogoutEvent(r.Context(), user, logoutToken); err != nil {
				s.Logger.Warn().Err(err).Msg("could not publish backchannel logout event")
			}
			s.deleteMFAState(logger, user.GetId().GetOpaqueId())
		}

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

// resolveLogoutUser resolves the user behind the logging-out session from a
// userinfo-cache record, as needed both to publish the backchannel logout
// event and to clean up that user's MFA-related state.
func (s StaticRouteHandler) resolveLogoutUser(ctx context.Context, record *microstore.Record) (*cs3.User, error) {
	urecords, err := s.UserInfoCache.Read(string(record.Value))
	if err != nil {
		return nil, fmt.Errorf("reading userinfo cache: %w", err)
	}
	if len(urecords) == 0 {
		return nil, fmt.Errorf("userinfo not found")
	}

	var claims map[string]any
	if err = msgpack.UnmarshalAsMap(urecords[0].Value, &claims); err != nil {
		return nil, fmt.Errorf("could not unmarshal userinfo: %w", err)
	}

	oidcClaim, ok := claims[s.Config.UserOIDCClaim].(string)
	if !ok {
		return nil, fmt.Errorf("could not get claim %s", s.Config.UserOIDCClaim)
	}

	user, _, err := s.UserProvider.GetUserByClaims(ctx, s.Config.UserCS3Claim, oidcClaim)
	if err != nil || user.GetId() == nil {
		return nil, fmt.Errorf("could not get user by claims: %w", err)
	}
	return user, nil
}

// publishBackchannelLogoutEvent publishes a backchannel logout event when the callback revived from the identity provider
func (s StaticRouteHandler) publishBackchannelLogoutEvent(ctx context.Context, user *cs3.User, logoutToken *oidc.LogoutToken) error {
	if s.EventsPublisher == nil {
		return fmt.Errorf("the events publisher is not set")
	}

	e := events.BackchannelLogout{
		Executant: user.GetId(),
		SessionId: logoutToken.SessionId,
		Timestamp: utils.TSNow(),
	}

	if err := events.Publish(ctx, s.EventsPublisher, e); err != nil {
		return fmt.Errorf("could not publish user created event %w", err)
	}
	return nil
}

// deleteMFAState deletes the resolved user's per-user MFA status record and
// MFA-bound signing key, closing the window during which other sessions of
// that user could use the stored MFA status or that key. Deletion is
// best-effort: any error is logged but never surfaces as a backchannel
// logout failure, so the existing userinfo-cache cleanup and spec-mandated
// response behavior are unaffected.
func (s StaticRouteHandler) deleteMFAState(logger log.Logger, userID string) {
	if userID == "" {
		return
	}

	if s.MFAStore != nil {
		if err := s.MFAStore.Delete(mfaRecordKey(userID)); err != nil && !errors.Is(err, microstore.ErrNotFound) {
			logger.Warn().Err(err).Str("userID", userID).Msg("could not delete MFA status record at backchannel logout")
		}
	}

	if s.SigningKeyStore != nil {
		if err := s.SigningKeyStore.Delete(mfaSigningKeyName(userID)); err != nil && !errors.Is(err, microstore.ErrNotFound) {
			logger.Warn().Err(err).Str("userID", userID).Msg("could not delete MFA-bound signing key at backchannel logout")
		}
	}
}
