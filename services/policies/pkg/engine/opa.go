package engine

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/gabriel-vasile/mimetype"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config"
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

// Evaluate evaluates the opa policies and returns teh result.
func (o OPA) Evaluate(ctx context.Context, qs string, env Environment) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)
	defer cancel()

	q, err := rego.New(
		rego.Query(qs),
		rego.Load(o.policies, nil),
		getMimetype,
		getResource,
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

var getResource = rego.Function1(
	&rego.Function{
		Name:             "ocis_get_resource",
		Decl:             types.NewFunction(types.Args(types.S), types.A),
		Memoize:          true,
		Nondeterministic: true,
	},
	func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
		var url string

		if err := ast.As(a.Value, &url); err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		client := rhttp.GetHTTPClient(rhttp.Insecure(true))
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status code from Download %v", res.StatusCode)
		}

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(res.Body); err != nil {
			return nil, err
		}

		v, err := ast.InterfaceToValue(buf.Bytes())
		if err != nil {
			return nil, err
		}

		return ast.NewTerm(v), nil
	},
)

var getMimetype = rego.Function1(
	&rego.Function{
		Name:             "ocis_get_mimetype",
		Decl:             types.NewFunction(types.Args(types.A), types.S),
		Memoize:          true,
		Nondeterministic: true,
	},
	func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
		var body []byte

		if err := ast.As(a.Value, &body); err != nil {
			return nil, err
		}

		mimeInfo := mimetype.Detect(body).String()
		detectedMimetype := strings.Split(mimeInfo, ";")[0]
		v, err := ast.InterfaceToValue(detectedMimetype)
		if err != nil {
			return nil, err
		}

		return ast.NewTerm(v), nil
	},
)
