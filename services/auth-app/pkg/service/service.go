package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	applications "github.com/cs3org/go-cs3apis/cs3/auth/applications/v1beta1"
	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/cs3org/reva/v2/pkg/appctx"
	"github.com/cs3org/reva/v2/pkg/auth/scope"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/config"
	settings "github.com/owncloud/ocis/v2/services/settings/pkg/service/v0"
	"google.golang.org/grpc/metadata"
)

// AuthAppService defines the service interface.
type AuthAppService struct {
	log log.Logger
	cfg *config.Config
	gws pool.Selectable[gateway.GatewayAPIClient]
	m   *chi.Mux
	r   *roles.Manager
}

// NewAuthAppService initializes a new AuthAppService.
func NewAuthAppService(opts ...Option) (*AuthAppService, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	r := roles.NewManager(
		// TODO: caching?
		roles.Logger(o.Logger),
		roles.RoleService(o.RoleClient),
	)

	a := &AuthAppService{
		log: o.Logger,
		cfg: o.Config,
		gws: o.GatewaySelector,
		m:   o.Mux,
		r:   &r,
	}

	a.m.Route("/auth-app/tokens", func(r chi.Router) {
		r.Get("/", a.HandleList)
		r.Post("/", a.HandleCreate)
		r.Delete("/", a.HandleDelete)
	})

	return a, nil
}

// ServeHTTP implements the http.Handler interface.
func (a *AuthAppService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.m.ServeHTTP(w, r)
}

// HandleCreate handles the creation of app tokens
func (a *AuthAppService) HandleCreate(w http.ResponseWriter, r *http.Request) {
	gwc, err := a.gws.Next()
	if err != nil {
		http.Error(w, "error getting gateway client", http.StatusInternalServerError)
		return
	}

	ctx := getContext(r)

	q := r.URL.Query()
	expiry, err := time.ParseDuration(q.Get("expiry"))
	if err != nil {
		a.log.Info().Err(err).Msg("error parsing expiry")
		http.Error(w, "error parsing expiry. Use e.g. 30m or 72h", http.StatusBadRequest)
		return
	}

	cid := buildClientID(q.Get("userID"), q.Get("userName"))
	if cid != "" {
		if !a.cfg.AllowImpersonation {
			a.log.Error().Msg("impersonation is not allowed")
			http.Error(w, "impersonation is not allowed", http.StatusForbidden)
			return
		}
		ok, err := isAdmin(ctx, a.r)
		if err != nil {
			a.log.Error().Err(err).Msg("error checking if user is admin")
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if !ok {
			a.log.Error().Msg("user is not admin")
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		ctx, err = a.authenticateUser(cid, gwc)
		if err != nil {
			a.log.Error().Err(err).Msg("error authenticating user")
			http.Error(w, "error authenticating user", http.StatusInternalServerError)
			return
		}
	}

	scopes, err := scope.AddOwnerScope(map[string]*authpb.Scope{})
	if err != nil {
		a.log.Error().Err(err).Msg("error adding owner scope")
		http.Error(w, "error adding owner scope", http.StatusInternalServerError)
		return
	}

	res, err := gwc.GenerateAppPassword(ctx, &applications.GenerateAppPasswordRequest{
		TokenScope: scopes,
		Label:      "Generated via API",
		Expiration: utils.TimeToTS(time.Now().Add(expiry)),
	})
	if err != nil {
		a.log.Error().Err(err).Msg("error generating app password")
		http.Error(w, "error generating app password", http.StatusInternalServerError)
		return
	}
	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		a.log.Error().Str("status", res.GetStatus().GetCode().String()).Msg("error generating app password")
		http.Error(w, "error generating app password: "+res.GetStatus().GetMessage(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(res.GetAppPassword())
	if err != nil {
		a.log.Error().Err(err).Msg("error marshaling app password")
		http.Error(w, "error marshaling app password", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		a.log.Error().Err(err).Msg("error writing response")
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleList handles listing of app tokens
func (a *AuthAppService) HandleList(w http.ResponseWriter, r *http.Request) {
	gwc, err := a.gws.Next()
	if err != nil {
		a.log.Error().Err(err).Msg("error getting gateway client")
		http.Error(w, "error getting gateway client", http.StatusInternalServerError)
		return
	}

	ctx := getContext(r)

	res, err := gwc.ListAppPasswords(ctx, &applications.ListAppPasswordsRequest{})
	if err != nil {
		a.log.Error().Err(err).Msg("error listing app passwords")
		http.Error(w, "error listing app passwords", http.StatusInternalServerError)
		return
	}

	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		a.log.Error().Str("status", res.GetStatus().GetCode().String()).Msg("error listing app passwords")
		http.Error(w, "error listing app passwords: "+res.GetStatus().GetMessage(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(res.GetAppPasswords())
	if err != nil {
		a.log.Error().Err(err).Msg("error marshaling app passwords")
		http.Error(w, "error marshaling app passwords", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(b); err != nil {
		a.log.Error().Err(err).Msg("error writing response")
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleDelete handles deletion of app tokens
func (a *AuthAppService) HandleDelete(w http.ResponseWriter, r *http.Request) {
	gwc, err := a.gws.Next()
	if err != nil {
		a.log.Error().Err(err).Msg("error getting gateway client")
		http.Error(w, "error getting gateway client", http.StatusInternalServerError)
		return
	}

	ctx := getContext(r)

	pw := r.URL.Query().Get("token")
	if pw == "" {
		a.log.Info().Msg("missing token")
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	res, err := gwc.InvalidateAppPassword(ctx, &applications.InvalidateAppPasswordRequest{Password: pw})
	if err != nil {
		a.log.Error().Err(err).Msg("error invalidating app password")
		http.Error(w, "error invalidating app password", http.StatusInternalServerError)
		return
	}

	if res.GetStatus().GetCode() != rpc.Code_CODE_OK {
		a.log.Error().Str("status", res.GetStatus().GetCode().String()).Msg("error invalidating app password")
		http.Error(w, "error invalidating app password: "+res.GetStatus().GetMessage(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *AuthAppService) authenticateUser(clientID string, gwc gateway.GatewayAPIClient) (context.Context, error) {
	ctx := context.Background()
	authRes, err := gwc.Authenticate(ctx, &gateway.AuthenticateRequest{
		Type:         "machine",
		ClientId:     clientID,
		ClientSecret: a.cfg.MachineAuthAPIKey,
	})
	if err != nil {
		return nil, err
	}

	if authRes.GetStatus().GetCode() != rpc.Code_CODE_OK {
		return nil, errors.New("error authenticating user: " + authRes.GetStatus().GetMessage())
	}

	ctx = ctxpkg.ContextSetUser(ctx, &userpb.User{Id: authRes.GetUser().GetId()})
	return metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, authRes.GetToken()), nil
}

func getContext(r *http.Request) context.Context {
	ctx := r.Context()
	return metadata.AppendToOutgoingContext(ctx, ctxpkg.TokenHeader, r.Header.Get("X-Access-Token"))
}

func buildClientID(userID, userName string) string {
	switch {
	default:
		return ""
	case userID != "":
		return "userid:" + userID
	case userName != "":
		return "username:" + userName
	}
}

// isAdmin determines if the user in the context is an admin / has account management permissions
func isAdmin(ctx context.Context, rm *roles.Manager) (bool, error) {
	logger := appctx.GetLogger(ctx)

	u, ok := ctxpkg.ContextGetUser(ctx)
	uid := u.GetId().GetOpaqueId()
	if !ok || uid == "" {
		logger.Error().Str("userid", uid).Msg("user not in context")
		return false, errors.New("no user in context")
	}
	// get roles from context
	roleIDs, ok := roles.ReadRoleIDsFromContext(ctx)
	if !ok {
		logger.Debug().Str("userid", uid).Msg("No roles in context, contacting settings service")
		var err error
		roleIDs, err = rm.FindRoleIDsForUser(ctx, uid)
		if err != nil {
			logger.Err(err).Str("userid", uid).Msg("failed to get roles for user")
			return false, err
		}

		if len(roleIDs) == 0 {
			logger.Err(err).Str("userid", uid).Msg("user has no roles")
			return false, errors.New("user has no roles")
		}
	}

	// check if permission is present in roles of the authenticated account
	return rm.FindPermissionByID(ctx, roleIDs, settings.AccountManagementPermissionID) != nil, nil
}
