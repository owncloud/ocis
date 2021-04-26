package svc

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asim/go-micro/plugins/client/grpc/v3"
	merrors "github.com/asim/go-micro/v3/errors"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	revauser "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/pkg/token"
	"github.com/cs3org/reva/pkg/token/manager/jwt"
	"github.com/cs3org/reva/pkg/user"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	accounts "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/ocs/pkg/service/v0/response"
	storepb "github.com/owncloud/ocis/store/pkg/proto/v0"
	"github.com/pkg/errors"
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

// GetSelf returns the currently logged in user
func (o Ocs) GetSelf(w http.ResponseWriter, r *http.Request) {
	var account *accounts.Account
	var err error
	u, ok := user.ContextGetUser(r.Context())
	if !ok || u.Id == nil || u.Id.OpaqueId == "" {
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "user is missing an id")))
		return
	}

	account, err = o.getAccountService().GetAccount(r.Context(), &accounts.GetAccountRequest{
		Id: u.Id.OpaqueId,
	})

	if err != nil {
		merr := merrors.FromError(err)
		// TODO(someone) this fix is in place because if the user backend (PROXY_ACCOUNT_BACKEND_TYPE) is set to, for instance,
		// cs3, we cannot count with the accounts service.
		if u != nil {
			uid, gid := o.extractUIDAndGID(u)
			d := &data.User{
				UserID:            u.Username,
				DisplayName:       u.DisplayName,
				LegacyDisplayName: u.DisplayName,
				Email:             u.Mail,
				UIDNumber:         uid,
				GIDNumber:         gid,
			}
			mustNotFail(render.Render(w, r, response.DataRender(d)))
			return
		}
		o.logger.Error().Err(merr).Interface("user", u).Msg("could not get account for user")
		return
	}

	// remove password from log if it is set
	if account.PasswordProfile != nil {
		account.PasswordProfile.Password = ""
	}
	o.logger.Debug().Interface("account", account).Msg("got user")

	d := &data.User{
		UserID:            account.OnPremisesSamAccountName,
		DisplayName:       account.DisplayName,
		LegacyDisplayName: account.DisplayName,
		Email:             account.Mail,
		UIDNumber:         account.UidNumber,
		GIDNumber:         account.GidNumber,
		// TODO hide enabled flag or it might get rendered as false
	}
	mustNotFail(render.Render(w, r, response.DataRender(d)))
}

// GetUser returns the user with the given userid
func (o Ocs) GetUser(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")
	var account *accounts.Account
	var err error

	switch {
	case userid == "":
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "missing user in context")))
	case o.config.AccountBackend == "accounts":
		account, err = o.fetchAccountByUsername(r.Context(), userid)
	case o.config.AccountBackend == "cs3":
		account, err = o.fetchAccountFromCS3Backend(r.Context(), userid)
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(merr).Str("userid", userid).Msg("could not get account for user")
		return
	}

	// remove password from log if it is set
	if account.PasswordProfile != nil {
		account.PasswordProfile.Password = ""
	}
	o.logger.Debug().Interface("account", account).Msg("got user")

	// mimic the oc10 bool as string for the user enabled property
	var enabled string
	if account.AccountEnabled {
		enabled = "true"
	} else {
		enabled = "false"
	}

	d := &data.User{
		UserID:            account.OnPremisesSamAccountName,
		DisplayName:       account.DisplayName,
		LegacyDisplayName: account.DisplayName,
		Email:             account.Mail,
		UIDNumber:         account.UidNumber,
		GIDNumber:         account.GidNumber,
		Enabled:           enabled, // TODO include in response only when admin?
		// TODO query storage registry for free space? of home storage, maybe...
		Quota: &data.Quota{
			Free:       2840756224000,
			Used:       5059416668,
			Total:      2845815640668,
			Relative:   0.18,
			Definition: "default",
		},
	}
	mustNotFail(render.Render(w, r, response.DataRender(d)))
}

// AddUser creates a new user account
func (o Ocs) AddUser(w http.ResponseWriter, r *http.Request) {
	userid := r.PostFormValue("userid")
	password := r.PostFormValue("password")
	displayname := r.PostFormValue("displayname")
	email := r.PostFormValue("email")
	uid := r.PostFormValue("uidnumber")
	gid := r.PostFormValue("gidnumber")

	var uidNumber, gidNumber int64
	var err error

	if uid != "" {
		uidNumber, err = strconv.ParseInt(uid, 10, 64)
		if err != nil {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "Cannot use the uidnumber provided")))
			o.logger.Error().Err(err).Str("userid", userid).Msg("Cannot use the uidnumber provided")
			return
		}
	}
	if gid != "" {
		gidNumber, err = strconv.ParseInt(gid, 10, 64)
		if err != nil {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "Cannot use the gidnumber provided")))
			o.logger.Error().Err(err).Str("userid", userid).Msg("Cannot use the gidnumber provided")
			return
		}
	}

	// fallbacks
	/* TODO decide if we want to make these fallbacks. Keep in mind:
	- oCIS requires a preferred_name and email
	*/
	if displayname == "" {
		displayname = userid
	}

	newAccount := &accounts.Account{
		Id:                       userid,
		DisplayName:              displayname,
		PreferredName:            userid,
		OnPremisesSamAccountName: userid,
		PasswordProfile: &accounts.PasswordProfile{
			Password: password,
		},
		Mail:           email,
		AccountEnabled: true,
	}

	if uidNumber != 0 {
		newAccount.UidNumber = uidNumber
	}

	if gidNumber != 0 {
		newAccount.GidNumber = gidNumber
	}

	var account *accounts.Account

	switch o.config.AccountBackend {
	case "accounts":
		account, err = o.getAccountService().CreateAccount(r.Context(), &accounts.CreateAccountRequest{
			Account: newAccount,
		})
	case "cs3":
		o.logger.Fatal().Msg("cs3 backend doesn't support adding users")
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}

	if err != nil {
		merr := merrors.FromError(err)
		switch merr.Code {
		case http.StatusBadRequest:
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, merr.Detail)))
		case http.StatusConflict:
			if response.APIVersion(r.Context()) == "2" {
				// it seems the application framework sets the ocs status code to the httpstatus code, which affects the provisioning api
				// see https://github.com/owncloud/core/blob/b9ff4c93e051c94adfb301545098ae627e52ef76/lib/public/AppFramework/OCSController.php#L142-L150
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, merr.Detail)))
			} else {
				mustNotFail(render.Render(w, r, response.ErrRender(data.MetaInvalidInput.StatusCode, merr.Detail)))
			}
		default:
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
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

	// mimic the oc10 bool as string for the user enabled property
	var enabled string
	if account.AccountEnabled {
		enabled = "true"
	} else {
		enabled = "false"
	}
	mustNotFail(render.Render(w, r, response.DataRender(&data.User{
		UserID:            account.OnPremisesSamAccountName,
		DisplayName:       account.DisplayName,
		LegacyDisplayName: account.DisplayName,
		Email:             account.Mail,
		UIDNumber:         account.UidNumber,
		GIDNumber:         account.GidNumber,
		Enabled:           enabled,
	})))
}

// EditUser creates a new user account
func (o Ocs) EditUser(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")

	var account *accounts.Account
	var err error
	switch o.config.AccountBackend {
	case "accounts":
		account, err = o.fetchAccountByUsername(r.Context(), userid)
	case "cs3":
		o.logger.Fatal().Msg("cs3 backend doesn't support editing users")
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("userid", userid).Msg("could not edit user")
		return
	}

	req := accounts.UpdateAccountRequest{
		Account: &accounts.Account{
			Id: account.Id,
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
		mustNotFail(render.Render(w, r, response.ErrRender(103, "unknown key '"+key+"'")))
		return
	}

	account, err = o.getAccountService().UpdateAccount(r.Context(), &req)
	if err != nil {
		merr := merrors.FromError(err)
		switch merr.Code {
		case http.StatusBadRequest:
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, merr.Detail)))
		default:
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("account_id", req.Account.Id).Str("user_id", userid).Msg("could not edit user")
		return
	}

	// remove password from log if it is set
	if account.PasswordProfile != nil {
		account.PasswordProfile.Password = ""
	}

	o.logger.Debug().Interface("account", account).Msg("updated user")
	mustNotFail(render.Render(w, r, response.DataRender(struct{}{})))
}

// DeleteUser deletes a user
func (o Ocs) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")

	var account *accounts.Account
	var err error
	switch o.config.AccountBackend {
	case "accounts":
		account, err = o.fetchAccountByUsername(r.Context(), userid)
	case "cs3":
		o.logger.Fatal().Msg("cs3 backend doesn't support deleting users")
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("userid", userid).Msg("could not delete user")
		return
	}

	if o.config.RevaAddress != "" && os.Getenv("STORAGE_USERS_DRIVER") != "owncloud" {
		t, err := o.mintTokenForUser(r.Context(), account)
		if err != nil {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, errors.Wrap(err, "could not mint token").Error()))
			return
		}

		ctx := metadata.AppendToOutgoingContext(r.Context(), token.TokenHeader, t)

		gwc, err := pool.GetGatewayServiceClient(o.config.RevaAddress)
		if err != nil {
			o.logger.Error().Err(err).Msg("error securing a connection to the reva gateway, could not delete user home")
		}

		homeResp, err := gwc.GetHome(ctx, &provider.GetHomeRequest{})
		if err != nil {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, errors.Wrap(err, "could not get home").Error()))
			return
		}

		if homeResp.Status.Code != rpcv1beta1.Code_CODE_OK {
			o.logger.Error().
				Str("stat_status_code", homeResp.Status.Code.String()).
				Str("stat_message", homeResp.Status.Message).
				Msg("DeleteUser: could not get user home: get failed")
			return
		}

		statResp, err := gwc.Stat(ctx, &provider.StatRequest{
			Ref: &provider.Reference{
				Spec: &provider.Reference_Path{
					Path: homeResp.Path,
				},
			},
		})

		if err != nil {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, errors.Wrap(err, "could not stat home").Error()))
			return
		}

		if statResp.Status.Code != rpcv1beta1.Code_CODE_OK {
			o.logger.Error().
				Str("stat_status_code", statResp.Status.Code.String()).
				Str("stat_message", statResp.Status.Message).
				Msg("DeleteUser: could not delete user home: stat failed")
			return
		}

		delReq := &provider.DeleteRequest{
			Ref: &provider.Reference{
				Spec: &provider.Reference_Id{
					Id: statResp.Info.Id,
				},
			},
		}

		delResp, err := gwc.Delete(ctx, delReq)
		if err != nil {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, errors.Wrap(err, "could not delete home").Error()))
			return
		}

		if delResp.Status.Code != rpcv1beta1.Code_CODE_OK {
			o.logger.Error().
				Str("stat_status_code", statResp.Status.Code.String()).
				Str("stat_message", statResp.Status.Message).
				Msg("DeleteUser: could not delete user home: delete failed")
			return
		}

		// delete trash is a combination of ListRecycle + PurgeRecycle (iterative)
		listRecycle := &gateway.ListRecycleRequest{
			Ref: &provider.Reference{
				Spec: &provider.Reference_Path{
					Path: homeResp.Path,
				},
			},
		}

		listRecycleResponse, err := gwc.ListRecycle(ctx, listRecycle)
		if err != nil {
			render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, errors.Wrap(err, "could not delete trash").Error()))
			return
		}

		if listRecycleResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
			o.logger.Error().
				Str("stat_status_code", statResp.Status.Code.String()).
				Str("stat_message", statResp.Status.Message).
				Msg("DeleteUser: could not delete user trash: delete failed")
			return
		}

		// now that we've got the items, we iterate, create references and fire PurgeRecycleRequests to the Reva gateway.
		//for i := range listRecycleResponse.RecycleItems {
		// craft purge request
		req := &gateway.PurgeRecycleRequest{
			Ref: &provider.Reference{
				Spec: &provider.Reference_Path{
					Path: homeResp.Path,
				},
			},
		}

		// do request
		purgeRecycleResponse, err := gwc.PurgeRecycle(ctx, req)
		if err != nil {

		}

		if purgeRecycleResponse.Status.Code != rpcv1beta1.Code_CODE_OK {
			o.logger.Error().
				Str("stat_status_code", statResp.Status.Code.String()).
				Str("stat_message", statResp.Status.Message).
				Msg("DeleteUser: could not delete user trash: delete failed")
			return
		}
	}

	req := accounts.DeleteAccountRequest{
		Id: account.Id,
	}

	_, err = o.getAccountService().DeleteAccount(r.Context(), &req)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("userid", req.Id).Msg("could not delete user")
		return
	}

	o.logger.Debug().Str("userid", req.Id).Msg("deleted user")
	mustNotFail(render.Render(w, r, response.DataRender(struct{}{})))
}

// TODO(refs) this to ocis-pkg ... we are minting tokens all over the place ... or use a service? ... like reva?
func (o Ocs) mintTokenForUser(ctx context.Context, account *accounts.Account) (string, error) {
	tm, _ := jwt.New(map[string]interface{}{
		"secret":  o.config.TokenManager.JWTSecret,
		"expires": int64(60),
	})

	u := &revauser.User{
		Id: &revauser.UserId{
			OpaqueId: account.Id,
			Idp:      o.config.IdentityManagement.Address,
		},
		Groups: []string{},
		Opaque: &types.Opaque{
			Map: map[string]*types.OpaqueEntry{
				"uid": {
					Decoder: "plain",
					Value:   []byte(strconv.FormatInt(account.UidNumber, 10)),
				},
				"gid": {
					Decoder: "plain",
					Value:   []byte(strconv.FormatInt(account.GidNumber, 10)),
				},
			},
		},
	}
	return tm.MintToken(ctx, u)
}

// EnableUser enables a user
func (o Ocs) EnableUser(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")

	var account *accounts.Account
	var err error
	switch o.config.AccountBackend {
	case "accounts":
		account, err = o.fetchAccountByUsername(r.Context(), userid)
	case "cs3":
		o.logger.Fatal().Msg("cs3 backend doesn't support enabling users")
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("userid", userid).Msg("could not enable user")
		return
	}

	account.AccountEnabled = true

	req := accounts.UpdateAccountRequest{
		Account: account,
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"AccountEnabled"},
		},
	}

	_, err = o.getAccountService().UpdateAccount(r.Context(), &req)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested account could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("account_id", account.Id).Msg("could not enable account")
		return
	}

	o.logger.Debug().Str("account_id", account.Id).Msg("enabled user")
	mustNotFail(render.Render(w, r, response.DataRender(struct{}{})))
}

// DisableUser disables a user
func (o Ocs) DisableUser(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userid")

	var account *accounts.Account
	var err error
	switch o.config.AccountBackend {
	case "accounts":
		account, err = o.fetchAccountByUsername(r.Context(), userid)
	case "cs3":
		o.logger.Fatal().Msg("cs3 backend doesn't support disabling users")
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}

	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested user could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("userid", userid).Msg("could not disable user")
		return
	}

	account.AccountEnabled = false

	req := accounts.UpdateAccountRequest{
		Account: account,
		UpdateMask: &field_mask.FieldMask{
			Paths: []string{"AccountEnabled"},
		},
	}

	_, err = o.getAccountService().UpdateAccount(r.Context(), &req)
	if err != nil {
		merr := merrors.FromError(err)
		if merr.Code == http.StatusNotFound {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaNotFound.StatusCode, "The requested account could not be found")))
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, err.Error())))
		}
		o.logger.Error().Err(err).Str("account_id", account.Id).Msg("could not disable account")
		return
	}

	o.logger.Debug().Str("account_id", account.Id).Msg("disabled user")
	mustNotFail(render.Render(w, r, response.DataRender(struct{}{})))
}

// GetSigningKey returns the signing key for the current user. It will create it on the fly if it does not exist
// The signing key is part of the user settings and is used by the proxy to authenticate requests
// Currently, the username is used as the OC-Credential
func (o Ocs) GetSigningKey(w http.ResponseWriter, r *http.Request) {
	u, ok := user.ContextGetUser(r.Context())
	if !ok {
		//o.logger.Error().Msg("missing user in context")
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "missing user in context")))
		return
	}

	// use the user's UUID
	userID := u.Id.OpaqueId

	c := storepb.NewStoreService("com.owncloud.api.store", grpc.NewClient())
	res, err := c.Read(r.Context(), &storepb.ReadRequest{
		Options: &storepb.ReadOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Key: userID,
	})
	if err == nil && len(res.Records) > 0 {
		mustNotFail(render.Render(w, r, response.DataRender(&data.SigningKey{
			User:       userID,
			SigningKey: string(res.Records[0].Value),
		})))
		return
	}
	if err != nil {
		e := merrors.Parse(err.Error())
		if e.Code == http.StatusNotFound {
			// not found is ok, so we can continue and generate the key on the fly
		} else {
			mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "error reading from store")))
			return
		}
	}

	// try creating it
	key := make([]byte, 64)
	_, err = rand.Read(key[:])
	if err != nil {
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not generate signing key")))
		return
	}
	signingKey := hex.EncodeToString(key)

	_, err = c.Write(r.Context(), &storepb.WriteRequest{
		Options: &storepb.WriteOptions{
			Database: "proxy",
			Table:    "signing-keys",
		},
		Record: &storepb.Record{
			Key:   userID,
			Value: []byte(signingKey),
			// TODO Expiry?
		},
	})

	if err != nil {
		//o.logger.Error().Err(err).Msg("error writing key")
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not persist signing key")))
		return
	}

	mustNotFail(render.Render(w, r, response.DataRender(&data.SigningKey{
		User:       userID,
		SigningKey: signingKey,
	})))
}

// ListUsers lists the users
func (o Ocs) ListUsers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	query := ""
	if search != "" {
		query = fmt.Sprintf("on_premises_sam_account_name eq '%s'", escapeValue(search))
	}

	var res *accounts.ListAccountsResponse
	var err error
	switch o.config.AccountBackend {
	case "accounts":
		res, err = o.getAccountService().ListAccounts(r.Context(), &accounts.ListAccountsRequest{
			Query: query,
		})
	case "cs3":
		// TODO
		o.logger.Fatal().Msg("cs3 backend doesn't support listing users")
	default:
		o.logger.Fatal().Msgf("Invalid accounts backend type '%s'", o.config.AccountBackend)
	}

	if err != nil {
		o.logger.Err(err).Msg("could not list users")
		mustNotFail(render.Render(w, r, response.ErrRender(data.MetaServerError.StatusCode, "could not list users")))
		return
	}

	users := make([]string, 0, len(res.Accounts))
	for i := range res.Accounts {
		users = append(users, res.Accounts[i].OnPremisesSamAccountName)
	}

	mustNotFail(render.Render(w, r, response.DataRender(&data.Users{Users: users})))
}

// escapeValue escapes all special characters in the value
func escapeValue(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func (o Ocs) fetchAccountByUsername(ctx context.Context, name string) (*accounts.Account, error) {
	var res *accounts.ListAccountsResponse
	res, err := o.getAccountService().ListAccounts(ctx, &accounts.ListAccountsRequest{
		Query: fmt.Sprintf("on_premises_sam_account_name eq '%v'", escapeValue(name)),
	})
	if err != nil {
		return nil, err
	}
	if res != nil && len(res.Accounts) == 1 {
		return res.Accounts[0], nil
	}
	return nil, merrors.NotFound("", "The requested user could not be found")
}

func (o Ocs) fetchAccountFromCS3Backend(ctx context.Context, name string) (*accounts.Account, error) {
	backend := o.getCS3Backend()
	u, err := backend.GetUserByClaims(ctx, "username", name, false)
	if err != nil {
		return nil, err
	}
	uid, gid := o.extractUIDAndGID(u)
	return &accounts.Account{
		OnPremisesSamAccountName: u.Username,
		DisplayName:              u.DisplayName,
		Mail:                     u.Mail,
		UidNumber:                uid,
		GidNumber:                gid,
	}, nil
}

func (o Ocs) extractUIDAndGID(u *cs3.User) (int64, int64) {
	var uid, gid int64
	var err error
	if u.Opaque != nil && u.Opaque.Map != nil {
		if uidObj, ok := u.Opaque.Map["uid"]; ok {
			if uidObj.Decoder == "plain" {
				uid, err = strconv.ParseInt(string(uidObj.Value), 10, 64)
				if err != nil {
					o.logger.Error().Err(err).Interface("user", u).Msg("could not extract uid for user")
				}
			}
		}
		if gidObj, ok := u.Opaque.Map["gid"]; ok {
			if gidObj.Decoder == "plain" {
				gid, err = strconv.ParseInt(string(gidObj.Value), 10, 64)
				if err != nil {
					o.logger.Error().Err(err).Interface("user", u).Msg("could not extract gid for user")
				}
			}
		}
	}
	return uid, gid
}
