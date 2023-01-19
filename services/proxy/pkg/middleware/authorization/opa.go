package authorization

import (
	"context"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
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
