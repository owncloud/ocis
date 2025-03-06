package opa

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/owncloud/reva/v2/pkg/rhttp"
)

// RFResourceDownload extends the rego dictionary with the possibility to download oCis resources.
//
// Rego: `ocis.resource.download("ocis/path/0034892347349827")`
// Result: bytes
var RFResourceDownload = rego.Function1(
	&rego.Function{
		Name:             "ocis.resource.download",
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
