package engine

import (
	"context"
	"encoding/json"

	user "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	v0 "github.com/owncloud/ocis/v2/protogen/gen/ocis/messages/policies/v0"
)

// Engine defines the granted handlers.
type Engine interface {
	Evaluate(ctx context.Context, query string, env Environment) (bool, error)
}

type (
	// Stage defines the used auth stage
	Stage string
)

var (
	// StagePP defines the post-processing stage
	StagePP Stage = "pp"

	// StageHTTP defines the http stage
	StageHTTP Stage = "http"
)

// Resource contains resource information and is used as part of the evaluated environment.
type Resource struct {
	ID   provider.ResourceId `json:"resource_id"`
	Name string              `json:"name"`
	URL  string              `json:"url"`
	Size uint64              `json:"size"`
}

// Request contains request information and is used as part of the evaluated environment.
type Request struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

// Environment contains every data that is needed to decide if the request should pass or not
type Environment struct {
	Stage    Stage     `json:"stage"`
	User     user.User `json:"user"`
	Request  Request   `json:"request"`
	Resource Resource  `json:"resource"`
}

// UnmarshalJSON satisfies the encoding/json Unmarshaler interface.
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
