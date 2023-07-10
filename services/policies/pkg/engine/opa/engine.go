package opa

import (
	"context"
	"time"

	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/topdown/print"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine"
)

// OPA wraps open policy agent makes it possible to ask if an action is granted.
type OPA struct {
	printHook print.Hook
	policies  []string
	timeout   time.Duration
}

// NewOPA returns a ready to use opa engine.
func NewOPA(timeout time.Duration, logger log.Logger, conf config.Engine) (OPA, error) {
	return OPA{
			policies:  conf.Policies,
			timeout:   timeout,
			printHook: logPrinter{logger: logger},
		},
		nil
}

// Evaluate evaluates the opa policies and returns the result.
func (o OPA) Evaluate(ctx context.Context, qs string, env engine.Environment) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	customFns := []func(r *rego.Rego){
		RFResourceDownload,
		RFMimetypeDetect,
		RFMimetypeExtensions,
	}

	q, err := rego.New(
		append([]func(r *rego.Rego){
			rego.Query(qs),
			rego.Load(o.policies, nil),
			rego.EnablePrintStatements(true),
			rego.PrintHook(o.printHook),
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
