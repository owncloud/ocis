package policy

import (
	"context"
	"encoding/json"
	"os"

	"github.com/asim/go-micro/v3/errors"

	"github.com/asim/go-micro/v3/client"
	"github.com/open-policy-agent/opa/rego"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
)

type clientWrapper struct {
	client.Client
	storage    IStorage
	policyPath string
}

var (
	maxAllowedUsers int64 = 5
)

func (c *clientWrapper) checkPolicy(ctx context.Context, req client.Request) error {
	if c.policyPath == "" {
		return nil
	}

	// let's operate with a single query gathering all data and see how far we can make it.
	r := rego.New(
		//Todo: spec out more query rules
		rego.Query(`users_count = input.external.users_count; deny = data.ocis.deny`),
		rego.Load([]string{c.policyPath}, nil),
	)

	// preparing queries in advance avoids parsing and compiling the policies on each query and improves performance considerably.
	// prepared queries are safe to share across multiple Go routines.
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return err
	}

	input := map[string]interface{}{
		"service":      req.Service(),
		"endpoint":     req.Endpoint(),
		"method":       req.Method(),
		"content_type": req.ContentType(),
		"stream":       req.Stream(),
		"external": map[string]interface{}{
			"users_count": c.storage.UsersCount(),
		},
	}

	if standardClaims := oidc.FromContext(ctx); standardClaims != nil {
		input["standard_claims"] = map[string]interface{}{
			"iss":     standardClaims.Iss,
			"sub":     standardClaims.Sub,
			"name":    standardClaims.Name,
			"email":   standardClaims.Email,
			"groups":  standardClaims.Groups,
			"ocis_id": standardClaims.OcisID,
		}
	}

	results, err := query.Eval(ctx, rego.EvalInput(input)) // provide input to correlate against the loaded data
	if err != nil {
		return err
	}

	// all conditions MUST evaluate to false in order to pass the policy checker.
	checkers := []func(set rego.ResultSet) bool{thumbnailerLicense}

	// all policies apply
	for _, f := range checkers {
		if f(results) {
			return errors.New(req.Service(), "denied policy result", 403)
		}
	}

	return nil
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if err := c.checkPolicy(ctx, req); err != nil {
		return err
	}

	return c.Client.Call(ctx, req, rsp, opts...)
}

func NewClientWrapper() client.Wrapper {
	// just for poc, needs to be part of conf options once accepted
	policyPath := os.Getenv("POC_POLICY_PATH")
	if _, err := os.Stat(policyPath); os.IsNotExist(err) {
		policyPath = ""
	}

	return func(c client.Client) client.Client {
		return &clientWrapper{
			Client:     c,
			storage:    NewStorage(), // defaults to localhost etcd store implementation. Enough for the POC. 💩 will hit the fan if no etcd instance is present, good enough for development.
			policyPath: policyPath,
		}
	}
}

// thumbnailerLicense permits requests if:
// - request is not denied (configured in the policy file)
// - user count under a threshold (default: 5)
func thumbnailerLicense(rs rego.ResultSet) bool {
	if current, err := rs[0].Bindings["users_count"].(json.Number).Int64(); err == nil {
		return current >= maxAllowedUsers && rs[0].Bindings["deny"].(bool)
	}
	return false
}
