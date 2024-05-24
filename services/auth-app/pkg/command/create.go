package command

import (
	"context"
	"fmt"
	authpb "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/auth/scope"

	applicationsv1beta1 "github.com/cs3org/go-cs3apis/cs3/auth/applications/v1beta1"
	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	typesv1beta1 "github.com/cs3org/go-cs3apis/cs3/types/v1beta1"
	ctxpkg "github.com/cs3org/reva/v2/pkg/ctx"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/config"
	"github.com/owncloud/ocis/v2/services/auth-app/pkg/config/parser"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/metadata"
	"time"
)

// Create is the entrypoint for the app auth create command
func Create(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "create",
		Usage:    "create an app auth token for a user",
		Category: "maintenance",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user-name",
				Value:    "",
				Usage:    "the user name",
				Required: true,
			},
			&cli.IntFlag{
				Name:     "expiration",
				Value:    72,
				Usage:    "expiration of the app password in hours",
				Required: false,
			},
		},
		Before: func(_ *cli.Context) error {
			return configlog.ReturnError(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			gatewaySelector, err := pool.GatewaySelector(
				cfg.Reva.Address,
				append(
					cfg.Reva.GetRevaOptions(),
					pool.WithRegistry(registry.GetRegistry()),
					pool.WithTracerProvider(traceProvider),
				)...)
			if err != nil {
				return err
			}

			next, err := gatewaySelector.Next()
			if err != nil {
				return err
			}

			userID := c.String("user-name")
			ctx := context.Background()
			authRes, err := next.Authenticate(ctx, &gatewayv1beta1.AuthenticateRequest{
				Type:         "machine",
				ClientId:     "username:" + userID,
				ClientSecret: cfg.Commons.MachineAuthAPIKey,
			})
			if err != nil {
				return err
			}
			granteeCtx := ctxpkg.ContextSetUser(context.Background(), &userpb.User{Id: authRes.GetUser().GetId()})
			granteeCtx = metadata.AppendToOutgoingContext(granteeCtx, ctxpkg.TokenHeader, authRes.GetToken())

			exp := c.Int("expiration")
			expiryDuration := time.Duration(exp) * time.Hour
			scopes, err := scope.AddOwnerScope(map[string]*authpb.Scope{})
			if err != nil {
				return err
			}

			appPassword, err := next.GenerateAppPassword(granteeCtx, &applicationsv1beta1.GenerateAppPasswordRequest{
				TokenScope: scopes,
				Label:      "Generated via CLI",
				Expiration: &typesv1beta1.Timestamp{
					Seconds: uint64(time.Now().Add(expiryDuration).Unix()),
				},
			})
			if err != nil {
				return err
			}

			fmt.Printf("App password created for %s", authRes.GetUser().GetUsername())
			fmt.Println()
			fmt.Printf(" password: %s", appPassword.GetAppPassword().GetPassword())
			fmt.Println()

			return nil
		},
	}
}
