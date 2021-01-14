package indexer

import (
	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	"sync"
)

// typeMap stores the indexer layout at runtime.
type fieldName = string
type tName = string
type typeMap struct {
	sync.Map
}

type typeMapping struct {
	PKFieldName    string
	IndicesByField map[fieldName][]index.Index
}

func (m *typeMap) allTypeMappings() map[tName]*typeMapping {
	var rv map[tName]*typeMapping

	m.Range(func(key, value interface{}) bool {
		rv[key.(string)] = value.(*typeMapping)
		return true
	})

	return rv
}

func (m *typeMap) getTypeMapping(typeName string) *typeMapping {
	if value, ok := m.Load(typeName); ok {
		return value.(*typeMapping)
	}

	return nil
}

func (m *typeMap) deleteTypeMapping(typeName string) {
	m.Delete(typeName)
}

func (m *typeMap) addIndex(typeName string, pkName string, idx index.Index) {
	if val, ok := m.Load(typeName); ok {
		rval := val.(*typeMapping)
		rval.IndicesByField[idx.IndexBy()] = append(rval.IndicesByField[idx.IndexBy()], idx)
		return
	}
	m.Store(typeName, &typeMapping{
		PKFieldName: pkName,
		IndicesByField: map[string][]index.Index{
			idx.IndexBy(): {idx},
		},
	})
}
