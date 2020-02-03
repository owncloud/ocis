package service

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/empty"
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
// This implementation replaces the existent data with the requested. It does not calculate diff
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
	contents, err := registry.Store.Read(req.Key)
	if err != nil {
		return err
	}

	if len(contents) > 0 {
		r := &proto.Payload{}
		json.Unmarshal(contents[0].Value, r)
		res.Payload = r
	}

	return nil
}

// List implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) List(ctx context.Context, in *empty.Empty, res *proto.Records) error {
	records := &proto.Records{}
	contents, err := registry.Store.List()
	if err != nil {
		return err
	}

	for _, v := range contents {
		records.Records = append(records.Records, &proto.Record{
			Key: v.Key,
		})
	}

	res.Records = records.Records

	return nil
}
