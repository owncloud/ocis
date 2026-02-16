package service

import (
	"context"

	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/thejerf/suture/v4"
)

// SutureService allows for the settings command to be embedded and supervised by a suture supervisor tree.
type SutureService struct {
	exec func(ctx context.Context) error
}

// NewSutureServiceBuilder creates a new suture service
func NewSutureServiceBuilder(f func(context.Context, *ociscfg.Config) error) func(*ociscfg.Config) suture.Service {
	return func(cfg *ociscfg.Config) suture.Service {
		return SutureService{
			exec: func(ctx context.Context) error {
				return f(ctx, cfg)
			},
		}
	}
}

// Serve to fullfil Server interface
func (s SutureService) Serve(ctx context.Context) error {
	return s.exec(ctx)
}
