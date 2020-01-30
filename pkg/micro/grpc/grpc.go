package grpc

// package grpc uses `ocis-pkg` to start a go-micro service
import (
	"context"

	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	svc "github.com/owncloud/ocis-accounts/pkg/service/v0"
	"github.com/owncloud/ocis-pkg/service/grpc"
)

// NewService initializes a new go-micro service ready to run
func NewService(c context.Context) grpc.Service {
	service := grpc.NewService(
		grpc.Name("accounts"),
		grpc.Namespace("com.owncloud"),
		grpc.Address("localhost:9999"),
		grpc.Context(c),
	)

	// add a handler to the service
	hdlr := svc.New()
	proto.RegisterSettingsServiceHandler(service.Server(), hdlr)

	service.Init()
	return service
}
