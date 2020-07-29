package svc

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/cs3org/reva/pkg/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/micro/go-micro/v2/client/grpc"
	merrors "github.com/micro/go-micro/v2/errors"
	accounts "github.com/owncloud/ocis-accounts/pkg/proto/v0"
	"github.com/owncloud/ocis-ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis-ocs/pkg/service/v0/response"
	storepb "github.com/owncloud/ocis-store/pkg/proto/v0"
)

// GetUser returns the currently logged in user
func (o Ocs) GetUser(w http.ResponseWriter, r *http.Request) {
	// TODO this endpoint needs authentication
	userid := chi.URLParam(r, "userid")

	if userid == "" {
		u, ok := user.ContextGetUser(r.Context())
		if !ok || u.Id == nil || u.Id.OpaqueId == "" {
			render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "missing user in context"))
			return
		}

		userid = u.Id.OpaqueId
	}

	accSvc := o.getAccountService()
	account, err := accSvc.GetAccount(r.Context(), &accounts.GetAccountRequest{
		Id: userid,
	})
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found"))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(err).Str("userid", userid).Msg("could not get user")
		return
	}

	// remove password from log if it is set
	if account.PasswordProfile != nil {
		account.PasswordProfile.Password = ""
	}
	o.logger.Debug().Interface("account", account).Msg("got user")

	render.Render(w, r, response.DataRender(&data.User{
		UserID:      account.Id, // TODO userid vs username! implications for clients if we return the userid here? -> implement graph ASAP?
		Username:    account.PreferredName,
		DisplayName: account.DisplayName,
		Email:       account.Mail,
		Enabled:     account.AccountEnabled,
	}))
}

// AddUser creates a new user account
func (o Ocs) AddUser(w http.ResponseWriter, r *http.Request) {
	// TODO this endpoint needs authentication
	userid := r.PostFormValue("userid")
	password := r.PostFormValue("password")
	username := r.PostFormValue("username")
	displayname := r.PostFormValue("displayname")
	email := r.PostFormValue("email")

	accSvc := o.getAccountService()
	account, err := accSvc.CreateAccount(r.Context(), &accounts.CreateAccountRequest{
		Account: &accounts.Account{
			DisplayName:              displayname,
			PreferredName:            username,
			OnPremisesSamAccountName: username,
			PasswordProfile: &accounts.PasswordProfile{
				Password: password,
			},
			Id:             userid,
			Mail:           email,
			AccountEnabled: true,
		},
	})
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusBadRequest {
			render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, merr.Detail))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(err).Str("userid", userid).Msg("could not add user")
		// TODO check error if account already existed
		return
	}

	// remove password from log if it is set
	if account.PasswordProfile != nil {
		account.PasswordProfile.Password = ""
	}
	o.logger.Debug().Interface("account", account).Msg("added user")

	render.Render(w, r, response.DataRender(&data.User{
		UserID:      account.Id,
		Username:    account.PreferredName,
		DisplayName: account.DisplayName,
		Email:       account.Mail,
		Enabled:     account.AccountEnabled,
	}))
}

// EditUser creates a new user account
func (o Ocs) EditUser(w http.ResponseWriter, r *http.Request) {
	// TODO this endpoint needs authentication
	req := accounts.UpdateAccountRequest{
		Account: &accounts.Account{
			Id: chi.URLParam(r, "userid"),
		},
	}
	key := r.PostFormValue("key")
	value := r.PostFormValue("value")

	switch key {
	case "email":
		req.Account.Mail = value
		req.UpdateMask = &fieldmaskpb.FieldMask{Paths: []string{"Mail"}}
	case "username":
		req.Account.PreferredName = value
		req.Account.OnPremisesSamAccountName = value
		req.UpdateMask = &fieldmaskpb.FieldMask{Paths: []string{"PreferredName", "OnPremisesSamAccountName"}}
	case "password":
		req.Account.PasswordProfile = &accounts.PasswordProfile{
			Password: value,
		}
		req.UpdateMask = &fieldmaskpb.FieldMask{Paths: []string{"PasswordProfile.Password"}}
	case "displayname", "display":
		req.Account.DisplayName = value
		req.UpdateMask = &fieldmaskpb.FieldMask{Paths: []string{"DisplayName"}}
	default:
		// https://github.com/owncloud/core/blob/24b7fa1d2604a208582055309a5638dbd9bda1d1/apps/provisioning_api/lib/Users.php#L321
		render.Render(w, r, response.ErrRender(103, "unknown key '"+key+"'"))
		return
	}

	accSvc := o.getAccountService()
	account, err := accSvc.UpdateAccount(r.Context(), &req)
	if err != nil {
		merr := merrors.FromError(err)
		switch merr.Code {
		case http.StatusNotFound:
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found"))
		case http.StatusBadRequest:
			render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, merr.Detail))
		default:
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(err).Str("userid", req.Account.Id).Msg("could not edit user")
		return
	}

	// remove password from log if it is set
	if account.PasswordProfile != nil {
		account.PasswordProfile.Password = ""
	}
	o.logger.Debug().Interface("account", account).Msg("updated user")
	render.Render(w, r, response.DataRender(struct{}{}))
}

// DeleteUser deletes a user
func (o Ocs) DeleteUser(w http.ResponseWriter, r *http.Request) {
	req := accounts.DeleteAccountRequest{
		Id: chi.URLParam(r, "userid"),
	}
	accSvc := o.getAccountService()
	_, err := accSvc.DeleteAccount(r.Context(), &req)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found"))
		} else {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error()))
		}
		o.logger.Error().Err(err).Str("userid", req.Id).Msg("could not delete user")
		return
	}
	o.logger.Debug().Str("userid", req.Id).Msg("deleted user")
	render.Render(w, r, response.DataRender(struct{}{}))
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
	search := r.URL.Query().Get("search")
	query := ""
	if search != "" {
		query = fmt.Sprintf("id eq '%s' or on_premises_sam_account_name eq '%s'", escapeValue(search), escapeValue(search))
	}
	accSvc := o.getAccountService()
	res, err := accSvc.ListAccounts(r.Context(), &accounts.ListAccountsRequest{
		Query: query,
	})
	if err != nil {
		o.logger.Err(err).Msg("could not list users")
		render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not list users"))
		return
	}
	users := []string{}
	for i := range res.Accounts {
		users = append(users, res.Accounts[i].Id)
	}

	render.Render(w, r, response.DataRender(&data.Users{Users: users}))
}

// escapeValue escapes all special characters in the value
func escapeValue(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
