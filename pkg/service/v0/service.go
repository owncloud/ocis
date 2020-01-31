package service

import (
	"context"
	"encoding/json"
	"log"

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
	// uses a store manager to persist account information
	st := store.New()

	data, err := json.Marshal([]byte(`{"theme": "dark"}`))
	if err != nil {
		// deal with this accordingly and not panicking
		log.Panic(err)
	}
	record := mstore.Record{
		Key:   req.Id,
		Value: data,
	}

	return st.Write(&record)
}

// Get implements the SettingsServiceHandler interface generated on accounts.pb.micro.go
func (s Service) Get(c context.Context, req *proto.Query, res *proto.Record) error {
	return nil
}
