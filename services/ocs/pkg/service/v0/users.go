package svc

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"

	cs3user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/response"
	revactx "github.com/owncloud/reva/v2/pkg/ctx"
	"github.com/owncloud/reva/v2/pkg/utils"
	"go-micro.dev/v4/store"
)

// GetSigningKey returns the signing key for the current user. It will create it on the fly if it does not exist
// The signing key is part of the user settings and is used by the proxy to authenticate requests
// Currently, the username is used as the OC-Credential
func (o Ocs) GetSigningKey(w http.ResponseWriter, r *http.Request) {
	u, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		//o.logger.Error().Msg("missing user in context")
		o.mustRender(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "missing user in context"))
		return
	}

	// Refuse to disclose the signing key to public-share guests; the key belongs to
	// the share owner and would let the caller forge signed URLs as that user.
	if u.GetId().GetType() == cs3user.UserType_USER_TYPE_GUEST || utils.ExistsInOpaque(u.GetOpaque(), "public-share-role") {
		o.mustRender(w, r, response.ErrRender(data.MetaUnauthorized.StatusCode, "unauthorized"))
		return
	}

	// use the user's UUID
	userID := u.Id.OpaqueId

	// A request that itself currently carries MFA is issued the MFA-bound key
	// instead of the normal key. A request without MFA keeps getting the
	// normal key, unchanged. A non-MFA request must never receive the
	// MFA-bound key under any code path.
	storeKey := userID
	if revactx.HasMFA(r.Context()) {
		storeKey = mfaSigningKeyName(userID)
	}

	res, err := o.store.Read(storeKey)
	if err == nil && len(res) > 0 {
		o.mustRender(w, r, response.DataRender(&data.SigningKey{
			User:       userID,
			SigningKey: string(res[0].Value),
		}))
		return
	}
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// not found is ok, so we can continue and generate the key on the fly
		} else {
			o.logger.Error().Err(err).Msg("error reading from server")
			o.mustRender(w, r, response.ErrRender(data.MetaServerError.StatusCode, "error reading from store"))
			return
		}
	}

	// try creating it
	key := make([]byte, 64)
	_, err = rand.Read(key[:])
	if err != nil {
		o.mustRender(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not generate signing key"))
		return
	}
	signingKey := hex.EncodeToString(key)

	err = o.store.Write(&store.Record{Key: storeKey, Value: []byte(signingKey)})
	if err != nil {
		o.mustRender(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not persist signing key"))
		return
	}

	o.mustRender(w, r, response.DataRender(&data.SigningKey{
		User:       userID,
		SigningKey: signingKey,
	}))
}

// mfaSigningKeyName returns the store key for a user's MFA-bound signing key.
// This naming is a contract with the proxy's signed-URL verifier
// (services/proxy/pkg/middleware/signed_url_auth.go), which reads the
// MFA-bound key under this same name. Keep both sides in sync.
func mfaSigningKeyName(userID string) string {
	return userID + ":mfa"
}
