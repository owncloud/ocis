package svc

import (
	"context"
	"testing"

	"github.com/owncloud/ocis/ocis-pkg/middleware"
	"github.com/stretchr/testify/assert"
	"go-micro.dev/v4/metadata"
)

var (
	ctxWithUUID      = metadata.Set(context.Background(), middleware.AccountID, "61445573-4dbe-4d56-88dc-88ab47aceba7")
	ctxWithEmptyUUID = metadata.Set(context.Background(), middleware.AccountID, "")
	emptyCtx         = context.Background()

	scenarios = []struct {
		name        string
		accountUUID string
		ctx         context.Context
		expect      string
	}{
		{
			name:        "context with UUID; identifier = 'me'",
			ctx:         ctxWithUUID,
			accountUUID: "me",
			expect:      "61445573-4dbe-4d56-88dc-88ab47aceba7",
		},
		{
			name:        "context with empty UUID; identifier = 'me'",
			ctx:         ctxWithEmptyUUID,
			accountUUID: "me",
			expect:      "",
		},
		{
			name:        "context without UUID; identifier = 'me'",
			ctx:         emptyCtx,
			accountUUID: "me",
			expect:      "",
		},
		{
			name:        "context with UUID; identifier not 'me'",
			ctx:         ctxWithUUID,
			accountUUID: "",
			expect:      "",
		},
	}
)

func TestGetValidatedAccountUUID(t *testing.T) {
	for _, s := range scenarios {
		scenario := s
		t.Run(scenario.name, func(t *testing.T) {
			got := getValidatedAccountUUID(scenario.ctx, scenario.accountUUID)
			assert.NotPanics(t, func() {
				getValidatedAccountUUID(emptyCtx, scenario.accountUUID)
			})
			assert.Equal(t, scenario.expect, got)
		})
	}
}
