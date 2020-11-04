package registry

import (
	"github.com/owncloud/ocis/ocis-pkg/indexer/index"
	"github.com/owncloud/ocis/ocis-pkg/indexer/option"
)

// IndexConstructor is a constructor function for creating index.Index.
type IndexConstructor func(o ...option.Option) index.Index

// IndexConstructorRegistry undocumented.
var IndexConstructorRegistry = map[string]map[string]IndexConstructor{
	"disk": {},
	"cs3":  {},
}
