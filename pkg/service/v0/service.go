package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
)

// New returns a new instance of Service
func New() Service {
	return Service{}
}

// Service implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
type Service struct{}

// Set implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) Set(c context.Context, req *proto.Record, res *proto.Record) error {
	res.Id = uuid.New().String()
	res.Theme = "dark"
	return nil
}

// Get implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) Get(c context.Context, req *proto.Query, res *proto.Record) error {
	return nil
}
