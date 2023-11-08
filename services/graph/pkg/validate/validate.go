package validate

import (
	"context"
	"sync/atomic"

	"github.com/go-playground/validator/v10"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

var defaultValidator atomic.Value
var structMapValidations = map[any]map[string]string{
	&libregraph.DriveItemInvite{}: {
		"Recipients":         "min=1",
		"Roles":              "len=1", // currently it is not possible to set more than one role
		"ExpirationDateTime": "omitnil,gt",
	},
}

func init() {
	v := validator.New()

	for s, rules := range structMapValidations {
		v.RegisterStructValidationMapRules(rules, s)
	}

	defaultValidator.Store(v)
}

// Default returns the default validator.
func Default() *validator.Validate { return defaultValidator.Load().(*validator.Validate) }

// StructCtx validates a struct and returns the error.
func StructCtx(ctx context.Context, s interface{}) error {
	return Default().StructCtx(ctx, s)
}
