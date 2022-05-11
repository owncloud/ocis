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
	"google.golang.org/grpc/metadata"

	cs3gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	cs3identity "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	cs3rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	cs3storage "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/go-chi/chi/v5"
	"github.com/go-micro/plugins/v4/client/grpc"
	"github.com/owncloud/ocis/v2/extensions/ocs/pkg/service/v0/data"
	"github.com/owncloud/ocis/v2/extensions/ocs/pkg/service/v0/response"
	ocstracing "github.com/owncloud/ocis/v2/extensions/ocs/pkg/tracing"
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

	currentUser, ok := revactx.ContextGetUser(r.Context())
	if !ok {
		response.ErrRender(data.MetaServerError.StatusCode, "missing user in context")
		return
	}

	var user *cs3identity.User
	switch {
	case userid == "":
		o.mustRender(w, r, response.ErrRender(data.MetaBadRequest.StatusCode, "missing username"))
	case userid == currentUser.Username:
		user = currentUser
	case o.config.AccountBackend == "cs3":
		user, err = o.fetchUserFromCS3Backend(r.Context(), userid)
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
		Quota:             &data.Quota{},
	}

	// lightweight and federated users don't have access to their storage space
	if currentUser.Id.Type != cs3identity.UserType_USER_TYPE_LIGHTWEIGHT && currentUser.Id.Type != cs3identity.UserType_USER_TYPE_FEDERATED {
		o.fillPersonalQuota(r.Context(), d, user)
	}

	_, span := ocstracing.TraceProvider.
		Tracer("ocs").
		Start(r.Context(), "GetUser")
	defer span.End()

	o.mustRender(w, r, response.DataRender(d))
}

func (o Ocs) fillPersonalQuota(ctx context.Context, d *data.User, u *cs3identity.User) {

	gc, err := pool.GetGatewayServiceClient(o.config.Reva.Address)
	if err != nil {
		o.logger.Error().Err(err).Msg("error getting gateway client")
		return
	}

	aRes, err := gc.Authenticate(ctx, &cs3gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "userid:" + u.Id.OpaqueId,
		ClientSecret: o.config.MachineAuthAPIKey,
	})

	switch {
	case err != nil:
		o.logger.Error().Err(err).Msg("could not fill personal quota")
		return
	case aRes.Status.Code != cs3rpc.Code_CODE_OK:
		o.logger.Error().Interface("status", aRes.Status).Msg("could not fill personal quota")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, revactx.TokenHeader, aRes.Token)

	res, err := gc.ListStorageSpaces(ctx, &cs3storage.ListStorageSpacesRequest{
		Filters: []*cs3storage.ListStorageSpacesRequest_Filter{
			{
				Type: cs3storage.ListStorageSpacesRequest_Filter_TYPE_OWNER,
				Term: &cs3storage.ListStorageSpacesRequest_Filter_Owner{
					Owner: aRes.User.Id,
				},
			},
			{
				Type: cs3storage.ListStorageSpacesRequest_Filter_TYPE_SPACE_TYPE,
				Term: &cs3storage.ListStorageSpacesRequest_Filter_SpaceType{
					SpaceType: "personal",
				},
			},
		},
	})
	if err != nil {
		o.logger.Error().Err(err).Msg("error calling ListStorageSpaces")
		return
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		o.logger.Debug().Interface("status", res.Status).Msg("ListStorageSpaces returned non OK result")
		return
	}

	if len(res.StorageSpaces) == 0 {
		o.logger.Debug().Err(err).Msg("list spaces returned empty list")
		return
	}

	getQuotaRes, err := gc.GetQuota(ctx, &cs3gateway.GetQuotaRequest{Ref: &cs3storage.Reference{
		ResourceId: res.StorageSpaces[0].Root,
		Path:       ".",
	}})
	if err != nil {
		o.logger.Error().Err(err).Msg("error calling GetQuota")
		return
	}
	if res.Status.Code != cs3rpc.Code_CODE_OK {
		o.logger.Debug().Interface("status", res.Status).Msg("GetQuota returned non OK result")
		return
	}

	total := getQuotaRes.TotalBytes
	used := getQuotaRes.UsedBytes

	d.Quota = &data.Quota{
		Used: int64(used),
		// TODO support negative values or flags for the quota to carry special meaning: -1 = uncalculated, -2 = unknown, -3 = unlimited
		// for now we can only report total and used
		Total: int64(total),
		// we cannot differentiate between `default` or a human readable `1 GB` defanation.
		// The web ui can create a human readable string from the actual total if it is sot. Otherwise it has to leave out relative and total anyway.
		// Definition: "default",
	}

	// only calculate free and relative when total is available
	if total > 0 {
		d.Quota.Free = int64(total - used)
		d.Quota.Relative = float32(float64(used) / float64(total))
	} else {
		d.Quota.Definition = "none" // this indicates no quota / unlimited to the ui
	}
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

func (o Ocs) fetchUserFromCS3Backend(ctx context.Context, name string) (*cs3identity.User, error) {
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
