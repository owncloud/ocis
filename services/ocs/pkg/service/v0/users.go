package svc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"

	storemsg "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/store/v0"
	storesvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/store/v0"

	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-micro/plugins/v4/client/grpc"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/v2/services/ocs/pkg/service/v0/response"
	ocstracing "github.com/owncloud/ocis/v2/services/ocs/pkg/tracing"
	merrors "go-micro.dev/v4/errors"
)

// GetSelf returns the currently logged in user
func (o Ocs) GetSelf(w http.ResponseWriter, r *http.Request) {
	u, ok := revactx.ContextGetUser(r.Context())
	if !ok || u.Id == nil || u.Id.OpaqueId == "" {
		o.mustRender(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "user is missing an id"))
		return
	}
	d := &data.User{
		UserID:            u.Username,
		DisplayName:       u.DisplayName,
		LegacyDisplayName: u.DisplayName,
		Email:             u.Mail,
		UIDNumber:         u.UidNumber,
		GIDNumber:         u.GidNumber,
	}
	o.mustRender(w, r, response.DataRender(d))
	return
}

// GetUser returns the user with the given userid
func (o Ocs) GetUser(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")
	userid, err := url.PathUnescape(userid)
	if err != nil {
		o.mustRender(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
	}

	var user *cs3.User
	switch {
	case userid == "":
		o.mustRender(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "missing user in context"))
	case o.config.AccountBackend == "cs3":
		user, err = o.fetchAccountFromCS3Backend(r.Context(), userid)
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			o.mustRender(w, r, response.ErrRender(data.MetaNotFound.StatusCode, data.MessageUserNotFound))
		} else {
			o.mustRender(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(merr).Str("userid", userid).Msg("could not get account for user")
		return
	}

	o.logger.Debug().Interface("user", user).Msg("got user")

	d := &data.User{
		UserID:            user.Username,
		DisplayName:       user.DisplayName,
		LegacyDisplayName: user.DisplayName,
		Email:             user.Mail,
		UIDNumber:         user.UidNumber,
		GIDNumber:         user.GidNumber,
		Enabled:           "true", // TODO include in response only when admin?
		// TODO query storage registry for free space? of home storage, maybe...
		Quota: &data.Quota{
			Free:       2840756224000,
			Used:       5059416668,
			Total:      2845815640668,
			Relative:   0.18,
			Definition: "default",
		},
	}

	_, span := ocstracing.TraceProvider.
		Tracer("ocs").
		Start(r.Context(), "GetUser")
	defer span.End()

	o.mustRender(w, r, response.DataRender(d))
}

// AddUser creates a new user account
func (o Ocs) AddUser(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

// EditUser creates a new user account
func (o Ocs) EditUser(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

// DeleteUser deletes a user
func (o Ocs) DeleteUser(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

// EnableUser enables a user
func (o Ocs) EnableUser(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

// DisableUser disables a user
func (o Ocs) DisableUser(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

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

	// use the user's UUID
	userID := u.Id.OpaqueId

	c := storesvc.NewStoreService("com.owncloud.api.store", grpc.NewClient())
	res, err := c.Read(r.Context(), &storesvc.ReadRequest{
		Options: &storemsg.ReadOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Key: userID,
	})
	if err == nil && len(res.Records) > 0 {
		o.mustRender(w, r, response.DataRender(&data.SigningKey{
			User:       userID,
			SigningKey: string(res.Records[0].Value),
		}))
		return
	}
	if err != nil {
		e := merrors.Parse(err.Error())
		if e.Code == http.StatusNotFound {
			// not found is ok, so we can continue and generate the key on the fly
		} else {
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

	_, err = c.Write(r.Context(), &storesvc.WriteRequest{
		Options: &storemsg.WriteOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Record: &storemsg.Record{
			Key:   userID,
			Value: []byte(signingKey),
			// TODO Expiry?
		},
	})

	if err != nil {
		//o.logger.Error().Err(err).Msg("error writing key")
		o.mustRender(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not persist signing key"))
		return
	}

	o.mustRender(w, r, response.DataRender(&data.SigningKey{
		User:       userID,
		SigningKey: signingKey,
	}))
}

// ListUsers lists the users
func (o Ocs) ListUsers(w http.ResponseWriter, r *http.Request) {
	switch o.config.AccountBackend {
	case "cs3":
		// TODO
		o.cs3WriteNotSupported(w, r)
		return
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}
}

// escapeValue escapes all special characters in the value
func escapeValue(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func (o Ocs) fetchAccountFromCS3Backend(ctx context.Context, name string) (*cs3.User, error) {
	backend := o.getCS3Backend()
	u, _, err := backend.GetUserByClaims(ctx, "username", name, false)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (o Ocs) cs3WriteNotSupported(w http.ResponseWriter, r *http.Request) {
	o.logger.Warn().Msg("the CS3 backend does not support adding or updating users")
	o.NotImplementedStub(w, r)
	return
}
