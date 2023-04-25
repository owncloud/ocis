package opa

import (
	"context"
	"time"

	"github.com/open-policy-agent/opa/rego"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine"
)

// OPA wraps open policy agent makes it possible to ask if an action is granted.
type OPA struct {
	policies []string
	timeout  time.Duration
}

// NewOPA returns a ready to use opa engine.
func NewOPA(timeout time.Duration, conf config.Engine) (OPA, error) {
	return OPA{
			policies: conf.Policies,
			timeout:  timeout,
		},
		nil
}

// Evaluate evaluates the opa policies and returns the result.
func (o OPA) Evaluate(ctx context.Context, qs string, env engine.Environment) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	customFns := []func(r *rego.Rego){
		RFResourceDownload,
		RFMimetypeExtension,
		RFMimetypeDetect,
	}

	q, err := rego.New(
		append([]func(r *rego.Rego){
			rego.Query(qs),
			rego.Load(o.policies, nil),
		}, customFns...)...,
	).PrepareForEval(ctx)
	if err != nil {
		return false, err
	}

	result, err := q.Eval(ctx, rego.EvalInput(env))
	if err != nil {
		return false, err
	}

	return result.Allowed(), nil
}
