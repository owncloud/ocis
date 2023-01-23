package authorization

import (
	"context"
	"github.com/gabriel-vasile/mimetype"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/types"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"strings"
	"time"

	"github.com/open-policy-agent/opa/rego"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

// OPA utilizes open policy agent to enforce policies
type OPA struct {
	Logger log.Logger
	Config config.AuthorizationMiddlewareOPA
}

// NewOPAAuthorizer returns a ready to use OPA Authorizer.
func NewOPAAuthorizer(logger log.Logger, conf config.AuthorizationMiddlewareOPA) (OPA, error) {
	opaAuthorizer := OPA{
		Logger: logger,
		Config: conf,
	}

	return opaAuthorizer, nil
}

// Authorize implements the Authorizer interface to authorize requests via opa.
func (o OPA) Authorize(ctx context.Context, info Info) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(o.Config.Timeout))
	defer cancel()

	q, err := rego.New(
		rego.Query("data.ocis.authz.allow"),
		rego.Load(o.Config.Policies, nil),
		hasMimetype,
	).PrepareForEval(ctx)
	if err != nil {
		return false, err
	}

	result, err := q.Eval(ctx, rego.EvalInput(info))
	if err != nil {
		return false, err
	}

	allow := result.Allowed()

	return allow, nil
}

var hasMimetype = rego.Function2(
	&rego.Function{
		Name: "hasMimetype",
		Decl: types.NewFunction(types.Args(types.A, types.S), types.B),
	},
	func(_ rego.BuiltinContext, a, b *ast.Term) (*ast.Term, error) {
		var body []byte
		var expectedMimetype string

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
