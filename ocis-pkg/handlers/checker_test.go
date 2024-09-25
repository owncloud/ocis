package handlers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
)

func TestCheckHandler_AddCheck(t *testing.T) {
	c := handlers.NewCheckHandlerConfiguration().WithCheck("shared-check", func(ctx context.Context) error { return nil })

	t.Run("configured checks are unique once added", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("checks should be unique, got %v", r)
			}
		}()

		h1 := handlers.NewCheckHandler(c)
		h1.AddCheck("check-with-same-name", func(ctx context.Context) error { return nil })

		h2 := handlers.NewCheckHandler(c)
		h2.AddCheck("check-with-same-name", func(ctx context.Context) error { return nil })

		fmt.Print(1)
	})
}
