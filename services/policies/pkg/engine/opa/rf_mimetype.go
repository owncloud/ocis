package opa

import (
	"github.com/gabriel-vasile/mimetype"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"mime"
	"strings"
)

var RFMimetypeExtension = rego.Function1(
	&rego.Function{
		Name:             "ocis.mimetype.extension_for_mimetype",
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

		v, err := ast.InterfaceToValue(detectedExtensions)
		if err != nil {
			return nil, err
		}

		return ast.NewTerm(v), nil
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

		mimeInfo := mimetype.Detect(body).String()
		detectedMimetype := strings.Split(mimeInfo, ";")[0]
		v, err := ast.InterfaceToValue(detectedMimetype)
		if err != nil {
			return nil, err
		}

		return ast.NewTerm(v), nil
	},
)
