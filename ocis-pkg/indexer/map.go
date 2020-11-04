package indexer

import "github.com/owncloud/ocis/ocis-pkg/indexer/index"

// typeMap stores the indexer layout at runtime.

type typeMap map[tName]typeMapping
type tName = string
type fieldName = string

type typeMapping struct {
	PKFieldName    string
	IndicesByField map[fieldName][]index.Index
}

func (m typeMap) addIndex(typeName string, pkName string, idx index.Index) {
	if val, ok := m[typeName]; ok {
		val.IndicesByField[idx.IndexBy()] = append(val.IndicesByField[idx.IndexBy()], idx)
		return
	}
	m[typeName] = typeMapping{
		PKFieldName: pkName,
		IndicesByField: map[string][]index.Index{
			idx.IndexBy(): {idx},
		},
	}
}
