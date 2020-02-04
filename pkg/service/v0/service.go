package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/owncloud/ocis-accounts/pkg/account"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"github.com/owncloud/ocis-accounts/pkg/proto/v0"
)

// New returns a new instance of Service
func New(cfg *config.Config) Service {
	fmt.Printf("config type: %T", account.Registry["filesystem"])
	return Service{
		Config:  cfg,
		Manager: account.Registry[cfg.Manager](cfg),
	}
}

// Service implements the SettingsServiceHandler interface
type Service struct {
	Config  *config.Config
	Manager account.Manager
}

// Set implements the SettingsServiceHandler interface
// This implementation replaces the existent data with the requested. It does not calculate diff
func (s Service) Set(c context.Context, req *proto.Record, res *proto.Record) error {
	settingsJSON, err := json.Marshal(req.Payload)
	if err != nil {
		return err
	}

	s.Manager.Write(&account.Record{
		Key:   req.Key,
		Value: settingsJSON,
	})

	return nil
}

// Get implements the SettingsServiceHandler interface
func (s Service) Get(c context.Context, req *proto.Query, res *proto.Record) error {
	contents := s.Manager.Read(req.Key)

	r := &proto.Payload{}
	json.Unmarshal(contents.Value, r)
	res.Payload = r

	return nil
}

// List implements the SettingsServiceHandler interface
func (s Service) List(ctx context.Context, in *empty.Empty, res *proto.Records) error {
	// r := &proto.Records{}
	// contents, err := registry.Store.List()
	// if err != nil {
	// 	return err
	// }

	// for _, v := range contents {
	// 	r.Records = append(r.Records, &proto.Record{
	// 		Key: v.Key,
	// 	})
	// }

	// res.Records = r.Records

	return nil
}
