package grpc

import (
	"context"

	permissions "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	svc "github.com/owncloud/ocis/extensions/settings/pkg/service/v0"
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/ocis-pkg/version"
	settingssvc "github.com/owncloud/ocis/protogen/gen/ocis/services/settings/v0"
	"go-micro.dev/v4/api"
	"go-micro.dev/v4/server"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Logger(options.Logger),
		grpc.Name(options.Name),
		grpc.Version(version.String),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Context(options.Context),
		grpc.Flags(options.Flags...),
	)

	handle := svc.NewService(options.Config, options.Logger)
	if err := settingssvc.RegisterBundleServiceHandler(service.Server(), handle); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register Bundle service handler")
	}
	if err := settingssvc.RegisterValueServiceHandler(service.Server(), handle); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register Value service handler")
	}
	if err := settingssvc.RegisterRoleServiceHandler(service.Server(), handle); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register Role service handler")
	}
	if err := settingssvc.RegisterPermissionServiceHandler(service.Server(), handle); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register Permission service handler")
	}

	if err := RegisterCS3PermissionsServiceHandler(service.Server(), handle); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register CS3 Permission service handler")
	}

	return service
}

func RegisterCS3PermissionsServiceHandler(s server.Server, hdlr permissions.PermissionsAPIServer, opts ...server.HandlerOption) error {
	type permissionsService interface {
		CheckPermission(context.Context, *permissions.CheckPermissionRequest, *permissions.CheckPermissionResponse) error
	}
	type PermissionsAPI struct {
		permissionsService
	}
	h := &permissionsServiceHandler{hdlr}
	opts = append(opts, api.WithEndpoint(&api.Endpoint{
		Name:    "PermissionsService.Checkpermission",
		Path:    []string{"/api/v0/permissions/check-permission"},
		Method:  []string{"POST"},
		Body:    "*",
		Handler: "rpc",
	}))
	return s.Handle(s.NewHandler(&PermissionsAPI{h}, opts...))
}

type permissionsServiceHandler struct {
	api permissions.PermissionsAPIServer
}

func (h *permissionsServiceHandler) CheckPermission(ctx context.Context, req *permissions.CheckPermissionRequest, res *permissions.CheckPermissionResponse) error {
	r, err := h.api.CheckPermission(ctx, req)
	if r != nil {
		*res = *r
	}
	return err
}
