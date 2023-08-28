package bleve

import (
	"regexp"
	"strings"

	bQuery "github.com/blevesearch/bleve/v2/search/query"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

// LegacyBuilder implements the legacy Builder interface.
type LegacyBuilder struct{}

// Build translates the ast to a valid bleve query.
func (b LegacyBuilder) Build(qs string) (*ast.Ast, error) {
	return &ast.Ast{
		Base: &ast.Base{
			Loc: &ast.Location{
				Start: ast.Position{
					Line:   0,
					Column: 0,
				},
				End: ast.Position{
					Line:   0,
					Column: len(qs),
				},
				Source: &qs,
			},
		},
	}, nil
}

// LegacyCompiler represents a default bleve query formatter.
type LegacyCompiler struct{}

// Compile implements the default bleve query formatter which converts the bleve likes query search string to the bleve query.
func (c LegacyCompiler) Compile(givenAst *ast.Ast) (bQuery.Query, error) {
	return &bQuery.QueryStringQuery{
		Query: c.formatQuery(*givenAst.Base.Loc.Source),
	}, nil
}

func (c LegacyCompiler) formatQuery(q string) string {
	cq := q
	fields := []string{"RootID", "Path", "ID", "Name", "Size", "Mtime", "MimeType", "Type"}
	for _, field := range fields {
		cq = strings.ReplaceAll(cq, strings.ToLower(field)+":", field+":")
	}

	fieldRe := regexp.MustCompile(`\w+:[^ ]+`)
	if fieldRe.MatchString(cq) {
		nameTagesRe := regexp.MustCompile(`\+?(Name|Tags)`) // detect "Name", "+Name, "Tags" and "+Tags"
		parts := strings.Split(cq, " ")

		cq = ""
		for _, part := range parts {
			fieldParts := strings.SplitN(part, ":", 2)
			if len(fieldParts) > 1 {
				key := fieldParts[0]
				value := fieldParts[1]
				if nameTagesRe.MatchString(key) {
					value = strings.ToLower(value) // do a lowercase query on the lowercased fields
				}
				cq += key + ":" + value + " "
			} else {
				cq += part + " "
			}
		}
		return cq // Sophisticated field based search
	}

	// this is a basic filename search
	cq = strings.ReplaceAll(cq, ":", `\:`)
	return "Name:*" + strings.ReplaceAll(strings.ToLower(cq), " ", `\ `) + "*"
}
