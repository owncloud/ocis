package opa

import (
	"log"
	"mime"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/rakyll/magicmime"
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
		if err := mime.AddExtensionType(".oform", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"); err != nil {
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

		if err := magicmime.Open(magicmime.MAGIC_MIME_TYPE | magicmime.MAGIC_SYMLINK | magicmime.MAGIC_ERROR); err != nil {
			log.Fatal(err)
		}
		defer magicmime.Close()
		mimetype, err := magicmime.TypeByBuffer(body)
		if err != nil {
			return nil, err
		}
		return ast.StringTerm(strings.Split(mimetype, ";")[0]), nil
	},
)
