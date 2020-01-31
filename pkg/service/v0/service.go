package service

import (
	"context"
	"encoding/json"

	mstore "github.com/micro/go-micro/v2/store"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"

	"github.com/owncloud/ocis-accounts/pkg/registry"
)

// New returns a new instance of Service
func New() Service {
	return Service{}
}

// Service implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
type Service struct{}

// Set implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) Set(c context.Context, req *proto.Record, res *proto.Record) error {
	settingsJSON, err := json.Marshal(req.Payload)
	if err != nil {
		return err
	}

	record := mstore.Record{
		Key:   req.Key,
		Value: settingsJSON,
	}

	return registry.Store.Write(&record)
}

// Get implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) Get(c context.Context, req *proto.Query, res *proto.Record) error {
	contents, _ := registry.Store.Read(req.Key)

	if len(contents) > 0 {
		r := &proto.Payload{}
		json.Unmarshal(contents[0].Value, r)
		res.Payload = r
	}

	return nil
}
