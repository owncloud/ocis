package service

import (
	"context"

	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
)

// New returns a new instance of Service
func New() Service {
	return Service{}
}

// Service implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
type Service struct{}

// Set implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) Set(c context.Context, req *proto.SettingsRequest, res *proto.SettingsResponse) error {
	res.Response = &proto.AccountSettings{
		Name: req.Request.Name,
	}
	return nil
}

// Get implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) Get(c context.Context, req *proto.AccountQueryRequest, res *proto.SettingsResponse) error {
	res.Response = &proto.AccountSettings{
		Name: "hej",
	}
	return nil
}
