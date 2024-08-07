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
	types "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/auth/scope"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/config"
	"google.golang.org/grpc/metadata"
)

// AuthAppService defines the service interface.
type AuthAppService struct {
	log log.Logger
	cfg *config.Config
	gws pool.Selectable[gateway.GatewayAPIClient]
	m   *chi.Mux
}

// NewAuthAppService initializes a new AuthAppService.
func NewAuthAppService(opts ...Option) (*AuthAppService, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	a := &AuthAppService{
		log: o.Logger,
		cfg: o.Config,
		gws: o.GatewaySelector,
		m:   o.Mux,
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
	cid := buildClientID(q.Get("userID"), q.Get("userName"))
	if cid != "" {
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
		Expiration: &types.Timestamp{
			Seconds: uint64(time.Now().Add(time.Hour).Unix()),
		},
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
