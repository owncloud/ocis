package autoprop

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"sync"
)

// Meta contains generic metadata that should be kept in the context
// It's stored in thread-safe, key-value map
type Meta struct {
	rwmutex  sync.RWMutex
	metadata map[string][]string
}

// NewMeta creates a new empty instance of Meta
func NewMeta() *Meta {
	return &Meta{
		// rwmutex is autoinitialized
		metadata: make(map[string][]string),
	}
}

// NewMetaFromJsonString creates a new instance of Meta from the provided
// string.
// This is just a convenient method for a NewMeta + FromJsonString. As such,
// if the json string can't be unmarshalled, an empty Meta instance will be
// returned.
func NewMetaFromJsonString(j string) *Meta {
	meta := NewMeta()
	meta.FromJsonString(j)
	return meta
}

// AppendMeta appends a new value for the key.
// The key will be converted to lowercase.
func (m *Meta) AppendMeta(key, value string) {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()

	internalKey := strings.ToLower(key)
	_, ok := m.metadata[internalKey]
	if ok {
		m.metadata[internalKey] = append(m.metadata[internalKey], value)
	} else {
		m.metadata[internalKey] = []string{value}
	}
}

// DeleteMeta deletes the provided key. All the values associated to the
// key will be deleted.
// The key will be converted to lowercase.
func (m *Meta) DeleteMeta(key string) {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()

	internalKey := strings.ToLower(key)
	delete(m.metadata, internalKey)
}

// GetMeta gets the values for the provided key, or a list containing just
// the default value (def) if the key doesn't exists.
// The key will be converted to lowercase.
func (m *Meta) GetMeta(key, def string) []string {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()

	internalKey := strings.ToLower(key)
	val, ok := m.metadata[internalKey]
	if !ok {
		return []string{def}
	}
	return val
}

// GetMetaWithExists gets the values of a key and "true" if the key exists, or
// an empty string list and "false" if the key doesn't exist.
// The key will be converted to lowercase.
func (m *Meta) GetMetaWithExists(key string) ([]string, bool) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()

	internalKey := strings.ToLower(key)
	val, ok := m.metadata[internalKey]
	return val, ok
}

// Len gets the number of elements keys in the metadata.
func (m *Meta) Len() int {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()

	return len(m.metadata)
}

// Create a new Meta instance. All the keys and values will be copied over,
// and the new instance will be completely separated (both the old and new
// instances can be modified separately)
func (m *Meta) CreateCopy() *Meta {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()

	data := make(map[string][]string, len(m.metadata))
	for k, v := range m.metadata {
		data[k] = slices.Clone(v)
	}

	newMeta := &Meta{
		// rwmutex is autoinitialized
		metadata: data,
	}
	return newMeta
}

// CreateCopyAsMap creates a copy of the data in this instance and returns
// it as a map[string][]string.
// The returned instance can be modified independently
func (m *Meta) CreateCopyAsMap(prefix string) map[string][]string {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()

	data := make(map[string][]string, len(m.metadata))
	for k, v := range m.metadata {
		data[prefix+k] = slices.Clone(v)
	}
	return data
}

// ToJsonString returns a json string of the metadata currently stored
// inside the instance.
// It might return an empty string if there are problems with the marshalling
func (m *Meta) ToJsonString() string {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()

	jsonBytes, _ := json.Marshal(m.metadata)
	return string(jsonBytes)
}

// FromJsonString uses the passed string as json and, if successful, it
// overwrites the stored metadata information of this instance.
// FromJsonString is expected to be used along with ToJsonString in order
// to serialize the information.
func (m *Meta) FromJsonString(j string) {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()

	var metaCopy map[string][]string
	if json.Unmarshal([]byte(j), &metaCopy) == nil {
		m.metadata = metaCopy
	}
}

type ctxMetaKey struct{}

// GetMetaFromContext gets the Meta from the context or nil if there is no
// Meta in the context.
// NOTE: a pointer to the instance is returned. This means that the instance
// can be dynamically modified in different parts of the code. This is
// intentional because we don't want to create new contexts whenever we update
// the metadata.
func GetMetaFromContext(ctx context.Context) *Meta {
	meta := ctx.Value(ctxMetaKey{})
	if meta == nil {
		return nil
	}
	return meta.(*Meta)
}

// SetMetaToContext creates a new context with the provided metadata.
// The old context can still access to the old metadata.
func SetMetaToContext(ctx context.Context, meta *Meta) context.Context {
	return context.WithValue(ctx, ctxMetaKey{}, meta)
}

// CopyMetaToContext copies the metadata from the old context to the newly
// created one (new context is derived from the old). The metadata from the
// old context and the metadata from the new one are completely different
// and can be modified independently.
// If the old context doesn't have associated metadata, a new Meta instance
// will be used.
//
// This can be used right before starting new goroutines: we don't want
// different goroutines to modify the metadata of eachother.
func CopyMetaToContext(ctx context.Context) context.Context {
	return CopyMetaFromContextToContext(ctx, ctx)
}

// CopyMetaFromContextToContext copies the meta from the "in" context to
// the "out" context. The meta from both contexts can be modified
// independently.
// The returned context will be derived from the "out" context.
// If the "in" context doesn't have a meta, a new meta instance will be
// used, so the returned context is guaranteed to have a meta, although
// it could be empty.
func CopyMetaFromContextToContext(in, out context.Context) context.Context {
	meta := GetMetaFromContext(in)
	if meta == nil {
		meta = NewMeta()
	}

	copiedMeta := meta.CreateCopy()
	return SetMetaToContext(out, copiedMeta)
}

// AppendMetaToContext will add a new key-value to the meta in the context.
// If the meta doesn't exist, a new one will be created and included in
// the returned context.
// Note that if the provided context already has a meta in it, no new
// context will be created, and the provided one will be returned (with
// the data appended in the meta).
// This method will use the meta.AppendMeta function, so if the key is
// already present, the value will be appended instead of replacing the
// existing value.
func AppendMetaToContext(ctx context.Context, key, value string) context.Context {
	meta := GetMetaFromContext(ctx)
	if meta == nil {
		meta = NewMeta()
		ctx = SetMetaToContext(ctx, meta)
	}
	meta.AppendMeta(key, value)
	return ctx
}
