//go:build ceph_preview

package rados

// GetOmapValuesOrdered fetches a set of keys and their values from an omap and
// returns them as an ordered slice of OmapKeyValue. Unlike GetOmapValues, which
// returns an unordered map, the returned slice preserves the key ordering
// provided by RADOS. This makes it suitable for pagination: the Key of the last
// element can be passed as `startAfter` in a subsequent call.
//
// `startAfter`: retrieve only the keys after this specified one
// `filterPrefix`: retrieve only the keys beginning with this prefix
// `maxReturn`: retrieve no more than `maxReturn` key/value pairs
func (ioctx *IOContext) GetOmapValuesOrdered(
	oid string, startAfter string, filterPrefix string, maxReturn int64,
) ([]OmapKeyValue, error) {
	omap := make([]OmapKeyValue, 0)

	err := ioctx.ListOmapValues(
		oid, startAfter, filterPrefix, maxReturn,
		func(key string, value []byte) {
			omap = append(omap, OmapKeyValue{Key: key, Value: value})
		},
	)

	return omap, err
}

// GetAllOmapValuesOrdered fetches all the keys and their values from an omap and
// returns them as an ordered slice of OmapKeyValue. Unlike GetAllOmapValues,
// which returns an unordered map, the returned slice preserves the key ordering
// provided by RADOS.
//
// `startAfter`: retrieve only the keys after this specified one
// `filterPrefix`: retrieve only the keys beginning with this prefix
// `iteratorSize`: internal number of keys to fetch during a read operation
func (ioctx *IOContext) GetAllOmapValuesOrdered(
	oid string, startAfter string, filterPrefix string, iteratorSize int64,
) ([]OmapKeyValue, error) {
	omap := make([]OmapKeyValue, 0)

	for {
		omapSize := len(omap)

		err := ioctx.ListOmapValues(
			oid, startAfter, filterPrefix, iteratorSize,
			func(key string, value []byte) {
				omap = append(omap, OmapKeyValue{Key: key, Value: value})
				startAfter = key
			},
		)
		if err != nil {
			return omap, err
		}

		// End of omap: the last iteration returned no new key/value pairs.
		if len(omap) == omapSize {
			break
		}
	}

	return omap, nil
}
