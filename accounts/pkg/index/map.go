package index

// indexMap stores the index layout at runtime.
type indexMap map[tName]map[indexByKey][]Type

type tName = string
type indexByKey = string

func (m indexMap) addIndex(idx Type) {
	typeName, indexBy := idx.TypeName(), idx.IndexBy()
	if _, ok := m[typeName]; !ok {
		m[typeName] = map[indexByKey][]Type{}
	}

	m[typeName][indexBy] = append(m[typeName][indexBy], idx)
}
