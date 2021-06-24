package policy

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/errors"
	"github.com/open-policy-agent/opa/rego"
	"github.com/owncloud/ocis/ocis-pkg/oidc"
	"os"
)

type clientWrapper struct {
	client.Client
	policyPath string
}

func (c *clientWrapper) checkPolicy(ctx context.Context, req client.Request) error {
	if c.policyPath == "" {
		return nil
	}

	r := rego.New(
		//Todo: spec out more query rules
		rego.Query(`deny = data.ocis.deny`),
		rego.Load([]string{c.policyPath}, nil),
	)

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
	fmt.Println(input)
	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return err
	} else if len(results) == 0 {
		// continue, that's ok for now. Maybe decide in configuration how to handle policies without value
		return nil
	} else if _, ok := results[0].Bindings["deny"].(bool); !ok {
		return errors.New(req.Service(), "Unexpected policy result", 500)
	}

	if results[0].Bindings["deny"].(bool) {
		return errors.New(req.Service(), "Denied policy result", 403)
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
			policyPath: policyPath,
		}
	}
}
