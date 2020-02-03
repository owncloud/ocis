package grpc

import (
	"context"

	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
	"github.com/owncloud/ocis-pkg/service/grpc"
)

// NewService creates a grpc service
func NewService(c context.Context) grpc.Service {
	service := grpc.NewService(
		grpc.Name("accounts"),
		grpc.Namespace("com.owncloud"),
		grpc.Address("localhost:9999"),
		grpc.Context(c),
	)

	hdlr := svc.New(config.New())
	proto.RegisterSettingsServiceHandler(service.Server(), hdlr)

	service.Init()
	return service
}
