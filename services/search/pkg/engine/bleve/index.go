package bleve

import (
	"errors"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

// IndexGetter is an interface that provides a way to get an index.
// Implementations might differ in how the index is created and how the
// index is gotten (reused, created on the fly, etc).
//
// Some implementations might require the index to be kept opened, meaning
// the index should be closed only when the application is shutting down. In
// this case, IndexCanBeClosed should return false. If the index can be
// closed and reopened safely at any time, IndexCanBeClosed should
// return true.
type IndexGetter interface {
	GetIndex(opts ...GetIndexOption) (bleve.Index, error)
	IndexCanBeClosed() bool
}

type IndexGetterMemory struct {
	mapping mapping.IndexMapping
	index   bleve.Index
}

// NewIndexGetterMemory creates a new IndexGetterMemory. This implementation
// creates a new in-memory index every time GetIndex is called. As such, the
// index must be kept opened. Closing the index will result in wiping the
// data.
func NewIndexGetterMemory(mapping mapping.IndexMapping) *IndexGetterMemory {
	return &IndexGetterMemory{
		mapping: mapping,
	}
}

// GetIndex creates a new in-memory index every time it is called.
// The options are ignored in this implementation.
func (i *IndexGetterMemory) GetIndex(opts ...GetIndexOption) (bleve.Index, error) {
	if i.index != nil {
		return i.index, nil
	}

	index, err := bleve.NewMemOnly(i.mapping)
	if err != nil {
		return nil, err
	}

	i.index = index
	return i.index, nil
}

// IndexCanBeClosed returns false, meaning the index must be kept opened.
func (i *IndexGetterMemory) IndexCanBeClosed() bool {
	return false
}

type IndexGetterPersistent struct {
	rootDir string
	mapping mapping.IndexMapping
	index   bleve.Index
}

// NewIndexGetterPersistent creates a new IndexGetterPersistent. The index
// will be persisted on the FS. If the index does not exist, it will be
// created. If the index exists, it will be opened.
//
// The index will be cached and reused every time GetIndex is called. You
// should not close the index unless you are shutting down the application.
func NewIndexGetterPersistent(rootDir string, mapping mapping.IndexMapping) *IndexGetterPersistent {
	return &IndexGetterPersistent{
		rootDir: rootDir,
		mapping: mapping,
	}
}

// GetIndex returns the cached index. The options are ignored in this
// implementation.
func (i *IndexGetterPersistent) GetIndex(opts ...GetIndexOption) (bleve.Index, error) {
	if i.index != nil {
		return i.index, nil
	}

	destination := filepath.Join(i.rootDir, "bleve")
	index, err := bleve.Open(destination)
	if errors.Is(bleve.ErrorIndexPathDoesNotExist, err) {
		index, err = bleve.New(destination, i.mapping)
		if err != nil {
			return nil, err
		}
	}

	i.index = index
	return i.index, nil
}

// IndexCanBeClosed returns false, meaning the index must be kept opened.
func (i *IndexGetterPersistent) IndexCanBeClosed() bool {
	return false
}

type IndexGetterPersistentScale struct {
	rootDir string
	mapping mapping.IndexMapping
}

// NewIndexGetterPersistentScale creates a new IndexGetterPersistentScale.
// The index will be persisted on the FS. If the index does not exist, it will
// be created. If the index exists, it will be opened.
// The GetIndex method will create a new connection to the index every time
// it is called. That connection must be closed after use.
func NewIndexGetterPersistentScale(rootDir string, mapping mapping.IndexMapping) *IndexGetterPersistentScale {
	return &IndexGetterPersistentScale{
		rootDir: rootDir,
		mapping: mapping,
	}
}

// GetIndex creates a new connection to the index every time it is called.
// You can use the ReadOnly option to open the index in read-only mode. This
// allow read-only operations to be performed in parallel.
// In order to avoid blocking write operations, you should close the index
// as soon as you are done with it.
func (i *IndexGetterPersistentScale) GetIndex(opts ...GetIndexOption) (bleve.Index, error) {
	options := newGetIndexOptions(opts...)
	destination := filepath.Join(i.rootDir, "bleve")
	params := map[string]interface{}{
		"read_only": options.ReadOnly,
	}
	index, err := bleve.OpenUsing(destination, params)
	if errors.Is(bleve.ErrorIndexPathDoesNotExist, err) {
		index, err = bleve.New(destination, i.mapping)
		if err != nil {
			return nil, err
		}

		return index, nil
	}

	return index, err
}

// IndexCanBeClosed returns true, meaning the index can be closed and
// reopened. You should close the index as soon as you are done with it.
func (i *IndexGetterPersistentScale) IndexCanBeClosed() bool {
	return true
}
