package svc

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/cs3org/reva/pkg/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/micro/go-micro/v2/client/grpc"
	merrors "github.com/micro/go-micro/v2/errors"
	accounts "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis-ocs/pkg/service/v0/response"
	storepb "github.com/owncloud/ocis-store/pkg/proto/v0"
)

// GetUser returns the currently logged in user
func (o Ocs) GetUser(w http.ResponseWriter, r *http.Request) {

	userid := chi.URLParam(r, "userid")

	if userid == "" {
		u, ok := user.ContextGetUser(r.Context())
		if !ok || u.Id == nil || u.Id.OpaqueId == "" {
			render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "missing user in context"))
			return
		}

		userid = u.Id.OpaqueId
	}

	accSvc := accounts.NewAccountsService("com.owncloud.api.accounts", grpc.NewClient())
	account, err := accSvc.GetAccount(r.Context(), &accounts.GetAccountRequest{
		Id: userid,
	})
	if err != nil {
		render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		return
	}

	render.Render(w, r, response.DataRender(&data.User{
		UserID:      account.Id, // TODO userid vs username! implications for clients if we return the userid here? -> implement graph ASAP?
		Username:    account.PreferredName,
		DisplayName: account.DisplayName,
		Email:       account.Mail,
		Enabled:     account.AccountEnabled,
	}))
}

// GetSigningKey returns the signing key for the current user. It will create it on the fly if it does not exist
// The signing key is part of the user settings and is used by the proxy to authenticate requests
// Currently, the username is used as the OC-Credential
func (o Ocs) GetSigningKey(w http.ResponseWriter, r *http.Request) {
	u, ok := user.ContextGetUser(r.Context())
	if !ok {
		//o.logger.Error().Msg("missing user in context")
		render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "missing user in context"))
		return
	}
	c := storepb.NewStoreService("com.owncloud.api.store", grpc.NewClient())
	res, err := c.Read(r.Context(), &storepb.ReadRequest{
		Options: &storepb.ReadOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Key: u.Username,
	})
	if err == nil && len(res.Records) > 0 {
		render.Render(w, r, response.DataRender(&data.SigningKey{
			User:       u.Username,
			SigningKey: string(res.Records[0].Value),
		}))
		return
	}
	if err != nil {
		e := merrors.Parse(err.Error())
		if e.Code == http.StatusNotFound {
			//o.logger.Debug().Str("username", u.Username).Msg("signing key not found")
			// not found is ok, so we can continue and generate the key on the fly
		} else {
			//o.logger.Err(err).Msg("error reading from store")
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "error reading from store"))
			return
		}
	}

	// try creating it
	key := make([]byte, 64)
	_, err = rand.Read(key[:])
	if err != nil {
		//o.logger.Error().Err(err).Msg("could not generate signing key")
		render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not generate signing key"))
		return
	}
	signingKey := hex.EncodeToString(key)

	_, err = c.Write(r.Context(), &storepb.WriteRequest{
		Options: &storepb.WriteOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Record: &storepb.Record{
			Key:   u.Username,
			Value: []byte(signingKey),
			// TODO Expiry?
		},
	})
	if err != nil {
		//o.logger.Error().Err(err).Msg("error writing key")
		render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not persist signing key"))
		return
	}

	render.Render(w, r, response.DataRender(&data.SigningKey{
		User:       u.Username,
		SigningKey: signingKey,
	}))
}

// ListUsers lists the users
func (o Ocs) ListUsers(w http.ResponseWriter, r *http.Request) {

	render.Render(w, r, response.ErrRender(data.MetaUnknownError.StatusCode, "please check the syntax. API specifications are here: http://www.freedesktop.org/wiki/Specifications/open-collaboration-services"))
}
