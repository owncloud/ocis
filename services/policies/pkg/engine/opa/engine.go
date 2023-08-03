package opa

import (
	"context"
	"io"
	"os"
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
	options   []func(r *rego.Rego)
}

// NewOPA returns a ready to use opa engine.
func NewOPA(timeout time.Duration, logger log.Logger, conf config.Engine) (OPA, error) {
	var mtReader io.ReadCloser

	if conf.Mimes != "" {
		var err error
		mtReader, err = os.Open(conf.Mimes)
		if err != nil {
			return OPA{}, err
		}

		defer mtReader.Close()
	}

	rfMimetypeExtensions, err := RFMimetypeExtensions(mtReader)
	if err != nil {
		return OPA{}, err
	}

	return OPA{
		policies:  conf.Policies,
		timeout:   timeout,
		printHook: logPrinter{logger: logger},
		options: []func(r *rego.Rego){
			RFMimetypeDetect,
			RFResourceDownload,
			rfMimetypeExtensions,
		},
	}, nil
}

// Evaluate evaluates the opa policies and returns the result.
func (o OPA) Evaluate(ctx context.Context, qs string, env engine.Environment) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	q, err := rego.New(
		append([]func(r *rego.Rego){
			rego.Query(qs),
			rego.Load(o.policies, nil),
			rego.EnablePrintStatements(true),
			rego.PrintHook(o.printHook),
		}, o.options...)...,
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
