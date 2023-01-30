package authz

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cs3org/reva/v2/pkg/rhttp"
	"github.com/gabriel-vasile/mimetype"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/owncloud/ocis/v2/services/authz/pkg/config"
)

type OPA struct {
	config *config.Config
}

func NewOPA(conf *config.Config) (OPA, error) {
	return OPA{
			config: conf,
		},
		nil
}

func (o OPA) Allowed(ctx context.Context, env Environment) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, o.config.OPA.Timeout)
	defer cancel()

	q, err := rego.New(
		rego.Query("data.ocis.authz.allow"),
		rego.Load(o.config.OPA.Policies, nil),
		loadResource,
		hasMimetype,
		convertBtoS,
	).PrepareForEval(ctx)
	if err != nil {
		return false, err
	}

	result, err := q.Eval(ctx, rego.EvalInput(env))
	if err != nil {
		return false, err
	}

	allow := result.Allowed()

	return allow, nil
}

var loadResource = rego.Function1(
	&rego.Function{
		Name: "loadResource",
		Decl: types.NewFunction(
			types.Args(
				types.Named("url", types.S).Description("download url"),
			),
			types.Named("bytes", types.A).Description("resource bytes"),
		),
		Memoize:          true,
		Nondeterministic: true,
	},
	func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
		var (
			url string
		)

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

var hasMimetype = rego.Function2(
	&rego.Function{
		Name: "hasMimetype",
		Decl: types.NewFunction(
			types.Args(
				types.Named("bytes", types.A).Description("lookup bytes"),
				types.Named("mimetype", types.S).Description("the mimetype to check for"),
			),
			types.Named("result", types.B).Description("result of the suffix check"),
		),
		Memoize:          true,
		Nondeterministic: true,
	},
	func(_ rego.BuiltinContext, a, b *ast.Term) (*ast.Term, error) {
		var (
			body             []byte
			expectedMimetype string
		)

		if err := ast.As(a.Value, &body); err != nil {
			return nil, err
		} else if err := ast.As(b.Value, &expectedMimetype); err != nil {
			return nil, err
		}

		mimeInfo := mimetype.Detect(body).String()
		detectedMimetype := strings.Split(mimeInfo, ";")[0]
		same := detectedMimetype == expectedMimetype
		v, err := ast.InterfaceToValue(same)
		if err != nil {
			return nil, err
		}

		return ast.NewTerm(v), nil
	},
)

var convertBtoS = rego.Function1(
	&rego.Function{
		Name: "convertBtoS",
		Decl: types.NewFunction(
			types.Args(
				types.Named("bytes", types.A).Description("input bytes"),
			),
			types.Named("result", types.S).Description("output string"),
		),
		Memoize:          true,
		Nondeterministic: true,
	},
	func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
		var (
			body []byte
		)

		if err := ast.As(a.Value, &body); err != nil {
			return nil, err
		}

		v, err := ast.InterfaceToValue(string(body))
		if err != nil {
			return nil, err
		}

		return ast.NewTerm(v), nil
	},
)
