package opa

import (
	"bufio"
	"io"
	"mime"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
)

// RFMimetypeExtensions extends the rego dictionary with the possibility of mapping mimetypes to file extensions.
// Be careful calling this multiple times with individual readers, the mime store is global,
// which results in one global store which holds all known mimetype mappings at once.
//
// Rego: `ocis.mimetype.extensions("application/pdf")`
// Result `[.pdf]`
func RFMimetypeExtensions(f io.Reader) (func(*rego.Rego), error) {
	if f != nil {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) <= 1 || fields[0][0] == '#' {
				continue
			}
			mimeType := fields[0]
			for _, ext := range fields[1:] {
				if ext[0] == '#' {
					break
				}
				if err := mime.AddExtensionType("."+ext, mimeType); err != nil {
					return nil, err
				}

			}
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return rego.Function1(
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
	), nil
}

// RFMimetypeDetect extends the rego dictionary with the possibility to detect mimetypes.
// Be careful, the list of known mimetypes is limited.
//
// Rego: `ocis.mimetype.extensions(".txt")`
// Result `text/plain`
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
