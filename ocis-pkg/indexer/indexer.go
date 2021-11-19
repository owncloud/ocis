// Package indexer provides symlink-based indexer for on-disk document-directories.
package indexer

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/owncloud/ocis/ocis-pkg/sync"

	"github.com/CiscoM31/godata"
	"github.com/iancoleman/strcase"
	"github.com/owncloud/ocis/ocis-pkg/indexer/config"
	"github.com/owncloud/ocis/ocis-pkg/indexer/errors"
	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	_ "github.com/owncloud/ocis/ocis-pkg/indexer/index/cs3"  // to populate index
	_ "github.com/owncloud/ocis/ocis-pkg/indexer/index/disk" // to populate index
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
	"github.com/owncloud/ocis/ocis-pkg/indexer/registry"
)

// Indexer is a facade to configure and query over multiple indices.
type Indexer struct {
	config  *config.Config
	indices typeMap
	mu      sync.NamedRWMutex
}

// IdxAddResult represents the result of an Add call on an index
type IdxAddResult struct {
	Field, Value string
}

// CreateIndexer creates a new Indexer.
func CreateIndexer(cfg *config.Config) *Indexer {
	return &Indexer{
		config:  cfg,
		indices: typeMap{},
		mu:      sync.NewNamedRWMutex(),
	}
}

// Reset takes care of deleting all indices from storage and from the internal map of indices
func (i *Indexer) Reset() error {
	for j := range i.indices {
		for _, indices := range i.indices[j].IndicesByField {
			for _, idx := range indices {
				err := idx.Delete()
				if err != nil {
					return err
				}
			}
		}
		delete(i.indices, j)
	}

	return nil
}

// AddIndex adds a new index to the indexer receiver.
func (i *Indexer) AddIndex(t interface{}, indexBy, pkName, entityDirName, indexType string, bound *option.Bound, caseInsensitive bool) error {
	f := registry.IndexConstructorRegistry[i.config.Repo.Backend][indexType]
	var idx index.Index

	if i.config.Repo.Backend == "cs3" {
		idx = f(
			option.CaseInsensitive(caseInsensitive),
			option.WithEntity(t),
			option.WithBounds(bound),
			option.WithTypeName(getTypeFQN(t)),
			option.WithIndexBy(indexBy),
			option.WithDataURL(i.config.Repo.CS3.DataURL),
			option.WithDataPrefix(i.config.Repo.CS3.DataPrefix),
			option.WithJWTSecret(i.config.Repo.CS3.JWTSecret),
			option.WithProviderAddr(i.config.Repo.CS3.ProviderAddr),
			option.WithServiceUser(i.config.ServiceUser),
		)
	} else {
		idx = f(
			option.CaseInsensitive(caseInsensitive),
			option.WithEntity(t),
			option.WithBounds(bound),
			option.WithTypeName(getTypeFQN(t)),
			option.WithIndexBy(indexBy),
			option.WithFilesDir(path.Join(i.config.Repo.Disk.Path, entityDirName)),
			option.WithDataDir(i.config.Repo.Disk.Path),
		)
	}

	i.indices.addIndex(getTypeFQN(t), pkName, idx)
	return idx.Init()
}

// Add a new entry to the indexer
func (i *Indexer) Add(t interface{}) ([]IdxAddResult, error) {
	typeName := getTypeFQN(t)

	i.mu.Lock(typeName)
	defer i.mu.Unlock(typeName)

	var results []IdxAddResult
	if fields, ok := i.indices[typeName]; ok {
		for _, indices := range fields.IndicesByField {
			for _, idx := range indices {
				pkVal := valueOf(t, fields.PKFieldName)
				idxByVal := valueOf(t, idx.IndexBy())
				value, err := idx.Add(pkVal, idxByVal)
				if err != nil {
					return []IdxAddResult{}, err
				}
				if value == "" {
					continue
				}
				results = append(results, IdxAddResult{Field: idx.IndexBy(), Value: value})
			}
		}
	}

	return results, nil
}

// FindBy finds a value on an index by field and value.
func (i *Indexer) FindBy(t interface{}, field string, val string) ([]string, error) {
	typeName := getTypeFQN(t)

	i.mu.RLock(typeName)
	defer i.mu.RUnlock(typeName)

	resultPaths := make([]string, 0)
	if fields, ok := i.indices[typeName]; ok {
		for _, idx := range fields.IndicesByField[strcase.ToCamel(field)] {
			idxVal := val
			res, err := idx.Lookup(idxVal)
			if err != nil {
				if errors.IsNotFoundErr(err) {
					continue
				}

				if err != nil {
					return nil, err
				}
			}

			resultPaths = append(resultPaths, res...)

		}
	}

	result := make([]string, 0, len(resultPaths))
	for _, v := range resultPaths {
		result = append(result, path.Base(v))
	}

	return result, nil
}

// Delete deletes all indexed fields of a given type t on the Indexer.
func (i *Indexer) Delete(t interface{}) error {
	typeName := getTypeFQN(t)

	i.mu.Lock(typeName)
	defer i.mu.Unlock(typeName)

	if fields, ok := i.indices[typeName]; ok {
		for _, indices := range fields.IndicesByField {
			for _, idx := range indices {
				pkVal := valueOf(t, fields.PKFieldName)
				idxByVal := valueOf(t, idx.IndexBy())
				if err := idx.Remove(pkVal, idxByVal); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// FindByPartial allows for glob search across all indexes.
func (i *Indexer) FindByPartial(t interface{}, field string, pattern string) ([]string, error) {
	typeName := getTypeFQN(t)

	i.mu.RLock(typeName)
	defer i.mu.RUnlock(typeName)

	resultPaths := make([]string, 0)
	if fields, ok := i.indices[typeName]; ok {
		for _, idx := range fields.IndicesByField[strcase.ToCamel(field)] {
			res, err := idx.Search(pattern)
			if err != nil {
				if errors.IsNotFoundErr(err) {
					continue
				}

				if err != nil {
					return nil, err
				}
			}

			resultPaths = append(resultPaths, res...)

		}
	}

	result := make([]string, 0, len(resultPaths))
	for _, v := range resultPaths {
		result = append(result, path.Base(v))
	}

	return result, nil

}

// Update updates all indexes on a value <from> to a value <to>.
func (i *Indexer) Update(from, to interface{}) error {
	typeNameFrom := getTypeFQN(from)
	typeNameTo := getTypeFQN(to)

	i.mu.Lock(typeNameFrom)
	defer i.mu.Unlock(typeNameFrom)

	if typeNameFrom != typeNameTo {
		return fmt.Errorf("update types do not match: from %v to %v", typeNameFrom, typeNameTo)
	}

	if fields, ok := i.indices[typeNameFrom]; ok {
		for fName, indices := range fields.IndicesByField {
			oldV := valueOf(from, fName)
			newV := valueOf(to, fName)
			pkVal := valueOf(from, fields.PKFieldName)
			for _, idx := range indices {
				if oldV == newV {
					continue
				}
				if oldV == "" {
					if _, err := idx.Add(pkVal, newV); err != nil {
						return err
					}
					continue
				}
				if newV == "" {
					if err := idx.Remove(pkVal, oldV); err != nil {
						return err
					}
					continue
				}
				if err := idx.Update(pkVal, oldV, newV); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Query parses an OData query into something our indexer.Index understands and resolves it.
func (i *Indexer) Query(ctx context.Context, t interface{}, q string) ([]string, error) {
	query, err := godata.ParseFilterString(ctx, q)
	if err != nil {
		return nil, err
	}

	tree := newQueryTree()
	if err := buildTreeFromOdataQuery(query.Tree, &tree); err != nil {
		return nil, err
	}

	results := make([]string, 0)
	if err := i.resolveTree(t, &tree, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// t is used to infer the indexed field names. When building an index search query, field names have to respect Golang
// conventions and be in PascalCase. For a better overview on this contemplate reading the reflection package under the
// indexer directory. Traversal of the tree happens in a pre-order fashion.
// TODO implement logic for `and` operators.
func (i *Indexer) resolveTree(t interface{}, tree *queryTree, partials *[]string) error {
	if partials == nil {
		return fmt.Errorf("return value cannot be nil: partials")
	}

	if tree.left != nil {
		_ = i.resolveTree(t, tree.left, partials)
	}

	if tree.right != nil {
		_ = i.resolveTree(t, tree.right, partials)
	}

	// by the time we're here we reached a leaf node.
	if tree.token != nil {
		switch tree.token.filterType {
		case "FindBy":
			operand, err := sanitizeInput(tree.token.operands)
			if err != nil {
				return err
			}

			r, err := i.FindBy(t, operand.field, operand.value)
			if err != nil {
				return err
			}

			*partials = append(*partials, r...)
		case "FindByPartial":
			operand, err := sanitizeInput(tree.token.operands)
			if err != nil {
				return err
			}

			r, err := i.FindByPartial(t, operand.field, fmt.Sprintf("%v*", operand.value))
			if err != nil {
				return err
			}

			*partials = append(*partials, r...)
		default:
			return fmt.Errorf("unsupported filter: %v", tree.token.filterType)
		}
	}

	*partials = dedup(*partials)
	return nil
}

type indexerTuple struct {
	field, value string
}

// sanitizeInput returns a tuple of fieldName + value to be applied on indexer.Index filters.
func sanitizeInput(operands []string) (*indexerTuple, error) {
	if len(operands) != 2 {
		return nil, fmt.Errorf("invalid number of operands for filter function: got %v expected 2", len(operands))
	}

	// field names are Go public types and by design they are in PascalCase, therefore we need to adhere to this rules.
	// for further information on this have a look at the reflection package.
	f := strcase.ToCamel(operands[0])

	// remove single quotes from value.
	v := strings.ReplaceAll(operands[1], "'", "")
	return &indexerTuple{
		field: f,
		value: v,
	}, nil
}

// buildTreeFromOdataQuery builds an indexer.queryTree out of a GOData ParseNode. The purpose of this intermediate tree
// is to transform godata operators and functions into supported operations on our index. At the time of this writing
// we only support `FindBy` and `FindByPartial` queries as these are the only implemented filters on indexer.Index(es).
func buildTreeFromOdataQuery(root *godata.ParseNode, tree *queryTree) error {
	if root.Token.Type == godata.ExpressionTokenFunc { // i.e "startswith", "contains"
		switch root.Token.Value {
		case "startswith":
			token := token{
				operator:   root.Token.Value,
				filterType: "FindByPartial",
				// TODO sanitize the number of operands it the expected one.
				operands: []string{
					root.Children[0].Token.Value, // field name, i.e: Name
					root.Children[1].Token.Value, // field value, i.e: Jac
				},
			}

			tree.insert(&token)
		default:
			return fmt.Errorf("operation not supported")
		}
	}

	if root.Token.Type == godata.ExpressionTokenLogical {
		switch root.Token.Value {
		case "or":
			tree.insert(&token{operator: root.Token.Value})
			for _, child := range root.Children {
				if err := buildTreeFromOdataQuery(child, tree.left); err != nil {
					return err
				}
			}
		case "eq":
			tree.insert(&token{
				operator:   root.Token.Value,
				filterType: "FindBy",
				operands: []string{
					root.Children[0].Token.Value,
					root.Children[1].Token.Value,
				},
			})
			for _, child := range root.Children {
				if err := buildTreeFromOdataQuery(child, tree.left); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("operator not supported")
		}
	}
	return nil
}
