package svc

import (
	"context"
	"testing"

	"github.com/owncloud/ocis-pkg/v2/middleware"
	"github.com/owncloud/ocis-settings/pkg/proto/v0"
	"github.com/stretchr/testify/assert"
)

var (
	ctxWithUUID      = context.WithValue(context.Background(), middleware.UUIDKey, "61445573-4dbe-4d56-88dc-88ab47aceba7")
	ctxWithEmptyUUID = context.WithValue(context.Background(), middleware.UUIDKey, "")
	emptyCtx         = context.Background()

	scenarios = []struct {
		name       string
		identifier *proto.Identifier
		ctx        context.Context
		expect     *proto.Identifier
	}{
		{
			name: "context with UUID; identifier = 'me'",
			ctx:  ctxWithUUID,
			identifier: &proto.Identifier{
				AccountUuid: "me",
			},
			expect: &proto.Identifier{
				AccountUuid: ctxWithUUID.Value(middleware.UUIDKey).(string),
			},
		},
		{
			name: "context without UUID; identifier = 'me'",
			ctx:  ctxWithEmptyUUID,
			identifier: &proto.Identifier{
				AccountUuid: "me",
			},
			expect: &proto.Identifier{
				AccountUuid: "",
			},
		},
		{
			name:       "context with UUID; identifier not 'me'",
			ctx:        ctxWithUUID,
			identifier: &proto.Identifier{},
			expect: &proto.Identifier{
				AccountUuid: "",
			},
		},
	}
)

func TestGetFailsafeIdentifier(t *testing.T) {
	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.name, func(t *testing.T) {
			got := getFailsafeIdentifier(scenario.ctx, scenario.identifier)
			assert.NotPanics(t, func() {
				getFailsafeIdentifier(emptyCtx, scenario.identifier)
			})
			assert.Equal(t, scenario.expect, got)
		})
	}
}
