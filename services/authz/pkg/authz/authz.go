package authz

import (
	"context"
	"encoding/json"
	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	v0 "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/authz/v0"
)

type Authorizer interface {
	Allowed(context.Context, Environment) (bool, error)
}

type (
	// Stage defines the used auth stage
	Stage string
)

var (
	// StagePP defines the post-processing stage state
	StagePP Stage = "pp"

	// StageHTTP defines the http stage state
	StageHTTP Stage = "http"
)

// Environment contains every data that is needed to decide if the request should pass or not
type Environment struct {
	Stage      Stage               `json:"stage"`
	Method     string              `json:"method"`
	Name       string              `json:"name"`
	URL        string              `json:"url"`
	Size       uint64              `json:"size"`
	User       user.User           `json:"user"`
	ResourceID provider.ResourceId `json:"resource_id"`
}

func (e *Environment) UnmarshalJSON(b []byte) error {
	type TMP Environment
	var tmp TMP

	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}

	*e = Environment(tmp)

	switch string(e.Stage) {
	case v0.Stage_STAGE_HTTP.String():
		e.Stage = StageHTTP
	case v0.Stage_STAGE_PP.String():
		e.Stage = StagePP
	}

	return nil
}

func Authorized(ctx context.Context, env Environment, authorizers ...Authorizer) (bool, error) {
	for _, authorizer := range authorizers {
		if allowed, err := authorizer.Allowed(ctx, env); err != nil {
			return false, err
		} else if !allowed {
			return false, nil
		}
	}

	return true, nil
}
