package grpc

import (
	"github.com/owncloud/ocis/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/ocis-pkg/service/http"
	"github.com/owncloud/ocis/settings/pkg/proto/v0"
	svc "github.com/owncloud/ocis/settings/pkg/service/v0"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Logger(options.Logger),
		grpc.Name(options.Name),
		grpc.Version(options.Config.Service.Version),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Context(options.Context),
	)

	handle := svc.NewService(options.Config, options.Logger)
	if err := proto.RegisterBundleServiceHandler(service.Server(), handle); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register Bundle service handler")
	}
	if err := proto.RegisterValueServiceHandler(service.Server(), handle); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register Value service handler")
	}
	if err := proto.RegisterRoleServiceHandler(service.Server(), handle); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register Role service handler")
	}
	if err := proto.RegisterPermissionServiceHandler(service.Server(), handle); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register Permission service handler")
	}

	service.Init()
	http.M.Unlock()
	return service
}
