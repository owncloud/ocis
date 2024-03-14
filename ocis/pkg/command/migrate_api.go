package command

import (
	"context"
	"encoding/json"
	"errors"
	rpcv1beta1 "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/parser"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"io"
	"net/http"
	"time"

	applicationsv1beta1 "github.com/cs3org/go-cs3apis/cs3/auth/applications/v1beta1"
	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	"github.com/cs3org/reva/v2/pkg/auth/scope"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/metadata"
)

type tokenExchangeRequest struct {
	EMail string `json:"email"`
}
type tokenExchangeResponse struct {
	AccessToken string `json:"accessToken"`
}
type httpHandler struct {
	SharedSecret      string
	Logger            *zerolog.Logger
	gatewaySelector   *pool.Selector[gatewayv1beta1.GatewayAPIClient]
	machineAuthAPIKey string
}

var lifetime = time.Duration(72) * time.Hour //nolint:mnd
var minSharedSecretLength = 32
var ErrNotFound = errors.New("user not found")

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if auth != "Bearer "+h.SharedSecret {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	request := tokenExchangeRequest{}
	err = json.Unmarshal(body, &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, err := h.createAppPassword(request.EMail, lifetime)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// set access token back
	response := tokenExchangeResponse{AccessToken: accessToken}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *httpHandler) createAppPassword(mail string, lifetime time.Duration) (string, error) {
	next, err := h.gatewaySelector.Next()
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	authRes, err := next.Authenticate(ctx, &gatewayv1beta1.AuthenticateRequest{
		Type:         "machine",
		ClientId:     "mail:" + mail,
		ClientSecret: h.machineAuthAPIKey,
	})
	if err != nil {
		return "", err
	}
	if authRes.GetStatus().GetCode() != rpcv1beta1.Code_CODE_OK {
		h.Logger.Error().Msg(authRes.GetStatus().GetMessage())
		return "", ErrNotFound
	}
	granteeCtx := ctxpkg.ContextSetUser(context.Background(), &userpb.User{Id: authRes.GetUser().GetId()})
	granteeCtx = metadata.AppendToOutgoingContext(granteeCtx, ctxpkg.TokenHeader, authRes.GetToken())

	scopes, err := scope.AddOwnerScope(map[string]*authpb.Scope{})
	if err != nil {
		return "", err
	}

	appPassword, err := next.GenerateAppPassword(granteeCtx, &applicationsv1beta1.GenerateAppPasswordRequest{
		TokenScope: scopes,
		Label:      "Generated via CLI",
		Expiration: &typesv1beta1.Timestamp{
			Seconds: uint64(time.Now().Add(lifetime).Unix()),
		},
	})
	if err != nil {
		return "", err
	}

	return appPassword.GetAppPassword().GetPassword(), nil
}

// StartMigrationAPI defines the entrypoint for the command to start the migration api service
func StartMigrationAPI(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "start-migration-api",
		Usage: "starts the migration api service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "shared-secret",
				Value:    "string",
				Usage:    "shared secret to gain api access",
				Required: true,
			},
		},
		Before: func(c *cli.Context) error {
			sharedSecret := c.String("shared-secret")
			if len(sharedSecret) < minSharedSecretLength {
				return errors.New("shared secret is too short. Minimum length is 32 characters")
			}

			err := configlog.ReturnError(parser.ParseConfig(cfg, false))
			if err != nil {
				return err
			}
			cfg.Reva = shared.DefaultRevaConfig()

			return nil
		},
		Action: func(c *cli.Context) error {
			log := logger()

			gatewaySelector, err := pool.GatewaySelector(
				cfg.Reva.Address,
				append(
					cfg.Reva.GetRevaOptions(),
					pool.WithRegistry(registry.GetRegistry()),
				)...)
			if err != nil {
				return err
			}

			sharedSecret := c.String("shared-secret")
			log.Info().Str("Secret", sharedSecret).Msg("User input")

			mux := http.NewServeMux()
			// API POST /satellites/tokenExchange with POST body {"email": "user1@example.com"}
			h := &httpHandler{
				SharedSecret:      sharedSecret,
				Logger:            log,
				machineAuthAPIKey: cfg.Commons.MachineAuthAPIKey,
				gatewaySelector:   gatewaySelector,
			}
			mux.Handle("/satellites/tokenExchange", h)

			return http.ListenAndServe(":9999", mux)
		},
	}
}
