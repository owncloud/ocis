package rados

// #cgo LDFLAGS: -lrados
// #include <rados/librados.h>
//
import "C"

// Iter supports iterating over objects in the ioctx.
type Iter struct {
	ctx       C.rados_list_ctx_t
	err       error
	entry     string
	namespace string
}

// IterToken supports reporting on and seeking to different positions.
type IterToken uint32

// Iter returns a Iterator object that can be used to list the object names in the current pool
func (ioctx *IOContext) Iter() (*Iter, error) {
	iter := Iter{}
	if cerr := C.rados_nobjects_list_open(ioctx.ioctx, &iter.ctx); cerr < 0 {
		return nil, getError(cerr)
	}
	return &iter, nil
}

// Token returns a token marking the current position of the iterator. To be used in combination with Iter.Seek()
func (iter *Iter) Token() IterToken {
	return IterToken(C.rados_nobjects_list_get_pg_hash_position(iter.ctx))
}

// Seek moves the iterator to the position indicated by the token.
func (iter *Iter) Seek(token IterToken) {
	C.rados_nobjects_list_seek(iter.ctx, C.uint32_t(token))
}

// Next retrieves the next object name in the pool/namespace iterator.
// Upon a successful invocation (return value of true), the Value method should
// be used to obtain the name of the retrieved object name. When the iterator is
// exhausted, Next returns false. The Err method should used to verify whether the
// end of the iterator was reached, or the iterator received an error.
//
// Example:
//
//	iter := pool.Iter()
//	defer iter.Close()
//	for iter.Next() {
//		fmt.Printf("%v\n", iter.Value())
//	}
//	return iter.Err()
func (iter *Iter) Next() bool {
	var cEntry *C.char
	var cNamespace *C.char
	if cerr := C.rados_nobjects_list_next(iter.ctx, &cEntry, nil, &cNamespace); cerr < 0 {
		iter.err = getError(cerr)
		return false
	}
	iter.entry = C.GoString(cEntry)
	iter.namespace = C.GoString(cNamespace)
	return true
}

// Value returns the current value of the iterator (object name), after a successful call to Next.
func (iter *Iter) Value() string {
	if iter.err != nil {
		return ""
	}
	return iter.entry
}

// Namespace returns the namespace associated with the current value of the iterator (object name), after a successful call to Next.
func (iter *Iter) Namespace() string {
	if iter.err != nil {
		return ""
	}
	return iter.namespace
}

// Err checks whether the iterator has encountered an error.
func (iter *Iter) Err() error {
	if iter.err == ErrNotFound {
		return nil
	}
	return iter.err
}

// Close the iterator cursor on the server. Be aware that iterators are not closed automatically
// at the end of iteration.
func (iter *Iter) Close() {
	C.rados_nobjects_list_close(iter.ctx)
}
