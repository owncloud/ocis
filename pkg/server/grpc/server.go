package grpc

import (
	"time"

	mclient "github.com/micro/go-micro/v2/client"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
	"github.com/owncloud/ocis-pkg/v2/roles"
	"github.com/owncloud/ocis-pkg/v2/service/grpc"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
)

// Server initializes a new go-micro service ready to run
func Server(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Name(options.Config.Server.Name),
		grpc.Context(options.Context),
		grpc.Address(options.Config.GRPC.Addr),
		grpc.Namespace(options.Config.GRPC.Namespace),
		grpc.Logger(options.Logger),
		grpc.Flags(options.Flags...),
	)

	var hdlr *svc.Service
	var err error

	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	rs := settings.NewRoleService("com.owncloud.api.settings", mclient.DefaultClient)
	roleManager := roles.NewManager(
		roles.CacheSize(1024),
		roles.CacheTTL(time.Hour*24*7),
		roles.Logger(options.Logger),
		roles.RoleService(rs),
	)

	if hdlr, err = svc.New(
		svc.Logger(options.Logger),
		svc.Config(options.Config),
		svc.RoleManager(&roleManager),
		svc.RoleService(rs),
	); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not initialize service handler")
	}
	if err = proto.RegisterAccountsServiceHandler(service.Server(), hdlr); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}
	if err = proto.RegisterGroupsServiceHandler(service.Server(), hdlr); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register groups handler")
	}

	service.Init()
	return service
}
