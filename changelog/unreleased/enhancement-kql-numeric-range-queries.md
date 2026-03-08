Enhancement: Support numeric range queries in KQL

The KQL parser now supports numeric range queries using comparison
operators (>=, <=, >, <) on numeric fields. Previously, range operators
only worked with DateTime values, causing queries like `size>=1048576`
or `photo.iso>=100` to silently fail by falling through to free-text
search.

Affected numeric fields: Size, photo.iso, photo.fNumber,
photo.focalLength, photo.orientation.

https://github.com/owncloud/ocis/pull/12094
https://github.com/owncloud/ocis/issues/12093
