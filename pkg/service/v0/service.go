package service

import (
	"context"
	"encoding/json"

	mstore "github.com/micro/go-micro/v2/store"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
	store "github.com/owncloud/ocis-accounts/pkg/store/filesystem"
)

// New returns a new instance of Service
func New() Service {
	return Service{}
}

// Service implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
type Service struct{}

// Set implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) Set(c context.Context, req *proto.Record, res *proto.Record) error {
	// TODO this should be a globally initialized struct
	st := store.New()

	settingsJSON, err := json.Marshal(req.Payload)
	if err != nil {
		// TODO log the error
		return err
	}

	record := mstore.Record{
		Key:   req.Key,
		Value: settingsJSON,
	}

	return st.Write(&record)
}

// Get implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) Get(c context.Context, req *proto.Query, res *proto.Record) error {
	res.Payload = &proto.Payload{
		Phoenix: &proto.Phoenix{
			Theme: "light",
		},
	}
	return nil
}
