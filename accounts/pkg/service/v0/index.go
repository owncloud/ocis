package service

import (
	"context"

	"github.com/owncloud/ocis/accounts/pkg/proto/v0"
)

func (s Service) RebuildIndex(ctx context.Context, request *proto.RebuildIndexRequest, response *proto.RebuildIndexResponse) error {
	if err := s.index.Reset(); err != nil {
		return err
	}
	response.Indices = []string{"foo", "bar"}
	return nil
}
