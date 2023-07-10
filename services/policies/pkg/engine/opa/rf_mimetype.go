package opa

import (
	"mime"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
)

var RFMimetypeExtensions = rego.Function1(
	&rego.Function{
		Name:             "ocis.mimetype.extensions",
		Decl:             types.NewFunction(types.Args(types.S), types.A),
		Memoize:          true,
		Nondeterministic: true,
	},
	func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
		var mt string

		if err := ast.As(a.Value, &mt); err != nil {
			return nil, err
		}

		detectedExtensions, err := mime.ExtensionsByType(mt)
		if err != nil {
			return nil, err
		}

		var mimeTerms []*ast.Term
		for _, extension := range detectedExtensions {
			mimeTerms = append(mimeTerms, ast.NewTerm(ast.String(extension)))
		}

		return ast.ArrayTerm(mimeTerms...), nil
	},
)

var RFMimetypeDetect = rego.Function1(
	&rego.Function{
		Name:             "ocis.mimetype.detect",
		Decl:             types.NewFunction(types.Args(types.A), types.S),
		Memoize:          true,
		Nondeterministic: true,
	},
	func(_ rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
		var body []byte

		if err := ast.As(a.Value, &body); err != nil {
			return nil, err
		}

		mimetype := mimetype.Detect(body).String()

		return ast.StringTerm(strings.Split(mimetype, ";")[0]), nil
	},
)
