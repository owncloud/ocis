package grpc

import (
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
	"github.com/owncloud/ocis-pkg/v2/service/grpc"
)

// NewService initializes a new go-micro service ready to run
func NewService(opts ...Option) grpc.Service {
	options := newOptions(opts...)

	service := grpc.NewService(
		grpc.Name(options.Name),
		grpc.Context(options.Context),
		grpc.Address(options.Address),
		grpc.Namespace(options.Namespace),
		grpc.Logger(options.Logger),
	)

	hdlr := svc.New(options.Config)
	if err := proto.RegisterAccountsServiceHandler(service.Server(), hdlr); err != nil {
		options.Logger.Fatal().Err(err).Msg("could not register service handler")
	}

	service.Init()
	return service
}
