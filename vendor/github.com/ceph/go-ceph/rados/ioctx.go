package rados

// #cgo LDFLAGS: -lrados
// #include <errno.h>
// #include <stdlib.h>
// #include <rados/librados.h>
//
// char* nextChunk(char **idx) {
// 	char *copy;
// 	copy = strdup(*idx);
// 	*idx += strlen(*idx) + 1;
// 	return copy;
// }
//
// #if __APPLE__
// #define ceph_time_t __darwin_time_t
// #define ceph_suseconds_t __darwin_suseconds_t
// #elif __GLIBC__
// #define ceph_time_t __time_t
// #define ceph_suseconds_t __suseconds_t
// #else
// #define ceph_time_t time_t
// #define ceph_suseconds_t suseconds_t
// #endif
import "C"

import (
	"syscall"
	"time"
	"unsafe"

	"github.com/ceph/go-ceph/internal/retry"
)

// CreateOption is passed to IOContext.Create() and should be one of
// CreateExclusive or CreateIdempotent.
type CreateOption int

const (
	// CreateExclusive if used with IOContext.Create() and the object
	// already exists, the function will return an error.
	CreateExclusive = C.LIBRADOS_CREATE_EXCLUSIVE
	// CreateIdempotent if used with IOContext.Create() and the object
	// already exists, the function will not return an error.
	CreateIdempotent = C.LIBRADOS_CREATE_IDEMPOTENT

	defaultListObjectsResultSize = 1000
	// listEndSentinel is the value returned by rados_list_object_list_is_end
	// when a cursor has reached the end of a pool
	listEndSentinel = 1
)

//revive:disable:var-naming old-yet-exported public api

// PoolStat represents Ceph pool statistics.
type PoolStat struct {
	// space used in bytes
	Num_bytes uint64
	// space used in KB
	Num_kb uint64
	// number of objects in the pool
	Num_objects uint64
	// number of clones of objects
	Num_object_clones uint64
	// num_objects * num_replicas
	Num_object_copies              uint64
	Num_objects_missing_on_primary uint64
	// number of objects found on no OSDs
	Num_objects_unfound uint64
	// number of objects replicated fewer times than they should be
	// (but found on at least one OSD)
	Num_objects_degraded uint64
	Num_rd               uint64
	Num_rd_kb            uint64
	Num_wr               uint64
	Num_wr_kb            uint64
}

//revive:enable:var-naming

// ObjectStat represents an object stat information
type ObjectStat struct {
	// current length in bytes
	Size uint64
	// last modification time
	ModTime time.Time
}

// LockInfo represents information on a current Ceph lock
type LockInfo struct {
	NumLockers int
	Exclusive  bool
	Tag        string
	Clients    []string
	Cookies    []string
	Addrs      []string
}

// IOContext represents a context for performing I/O within a pool.
type IOContext struct {
	ioctx C.rados_ioctx_t

	// Hold a reference back to the connection that the ioctx depends on so
	// that Go's GC doesn't trigger the Conn's finalizer before this
	// IOContext is destroyed.
	conn *Conn
}

// validate returns an error if the ioctx is not ready to be used
// with ceph C calls.
func (ioctx *IOContext) validate() error {
	if ioctx.ioctx == nil {
		return ErrInvalidIOContext
	}
	return nil
}

// Pointer returns a pointer reference to an internal structure.
// This function should NOT be used outside of go-ceph itself.
func (ioctx *IOContext) Pointer() unsafe.Pointer {
	return unsafe.Pointer(ioctx.ioctx)
}

// SetNamespace sets the namespace for objects within this IO context (pool).
// Setting namespace to a empty or zero length string sets the pool to the default namespace.
//
// Implements:
//  void rados_ioctx_set_namespace(rados_ioctx_t io,
//                                 const char *nspace);
func (ioctx *IOContext) SetNamespace(namespace string) {
	var cns *C.char
	if len(namespace) > 0 {
		cns = C.CString(namespace)
		defer C.free(unsafe.Pointer(cns))
	}
	C.rados_ioctx_set_namespace(ioctx.ioctx, cns)
}

// Create a new object with key oid.
//
// Implements:
//  void rados_write_op_create(rados_write_op_t write_op, int exclusive,
//                             const char* category)
func (ioctx *IOContext) Create(oid string, exclusive CreateOption) error {
	op := CreateWriteOp()
	defer op.Release()
	op.Create(exclusive)
	return op.operateCompat(ioctx, oid)
}

// Write writes len(data) bytes to the object with key oid starting at byte
// offset offset. It returns an error, if any.
func (ioctx *IOContext) Write(oid string, data []byte, offset uint64) error {
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	dataPointer := unsafe.Pointer(nil)
	if len(data) > 0 {
		dataPointer = unsafe.Pointer(&data[0])
	}

	ret := C.rados_write(ioctx.ioctx, coid,
		(*C.char)(dataPointer),
		(C.size_t)(len(data)),
		(C.uint64_t)(offset))

	return getError(ret)
}

// WriteFull writes len(data) bytes to the object with key oid.
// The object is filled with the provided data. If the object exists,
// it is atomically truncated and then written. It returns an error, if any.
func (ioctx *IOContext) WriteFull(oid string, data []byte) error {
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	ret := C.rados_write_full(ioctx.ioctx, coid,
		(*C.char)(unsafe.Pointer(&data[0])),
		(C.size_t)(len(data)))
	return getError(ret)
}

// Append appends len(data) bytes to the object with key oid.
// The object is appended with the provided data. If the object exists,
// it is atomically appended to. It returns an error, if any.
func (ioctx *IOContext) Append(oid string, data []byte) error {
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	ret := C.rados_append(ioctx.ioctx, coid,
		(*C.char)(unsafe.Pointer(&data[0])),
		(C.size_t)(len(data)))
	return getError(ret)
}

// Read reads up to len(data) bytes from the object with key oid starting at byte
// offset offset. It returns the number of bytes read and an error, if any.
func (ioctx *IOContext) Read(oid string, data []byte, offset uint64) (int, error) {
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	var buf *C.char
	if len(data) > 0 {
		buf = (*C.char)(unsafe.Pointer(&data[0]))
	}

	ret := C.rados_read(
		ioctx.ioctx,
		coid,
		buf,
		(C.size_t)(len(data)),
		(C.uint64_t)(offset))

	if ret >= 0 {
		return int(ret), nil
	}
	return 0, getError(ret)
}

// Delete deletes the object with key oid. It returns an error, if any.
func (ioctx *IOContext) Delete(oid string) error {
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	return getError(C.rados_remove(ioctx.ioctx, coid))
}

// Truncate resizes the object with key oid to size size. If the operation
// enlarges the object, the new area is logically filled with zeroes. If the
// operation shrinks the object, the excess data is removed. It returns an
// error, if any.
func (ioctx *IOContext) Truncate(oid string, size uint64) error {
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	return getError(C.rados_trunc(ioctx.ioctx, coid, (C.uint64_t)(size)))
}

// Destroy informs librados that the I/O context is no longer in use.
// Resources associated with the context may not be freed immediately, and the
// context should not be used again after calling this method.
func (ioctx *IOContext) Destroy() {
	C.rados_ioctx_destroy(ioctx.ioctx)
}

// GetPoolStats returns a set of statistics about the pool associated with this I/O
// context.
//
// Implements:
//  int rados_ioctx_pool_stat(rados_ioctx_t io,
//                            struct rados_pool_stat_t *stats);
func (ioctx *IOContext) GetPoolStats() (stat PoolStat, err error) {
	cStat := C.struct_rados_pool_stat_t{}
	ret := C.rados_ioctx_pool_stat(ioctx.ioctx, &cStat)
	if ret < 0 {
		return PoolStat{}, getError(ret)
	}
	return PoolStat{
		Num_bytes:                      uint64(cStat.num_bytes),
		Num_kb:                         uint64(cStat.num_kb),
		Num_objects:                    uint64(cStat.num_objects),
		Num_object_clones:              uint64(cStat.num_object_clones),
		Num_object_copies:              uint64(cStat.num_object_copies),
		Num_objects_missing_on_primary: uint64(cStat.num_objects_missing_on_primary),
		Num_objects_unfound:            uint64(cStat.num_objects_unfound),
		Num_objects_degraded:           uint64(cStat.num_objects_degraded),
		Num_rd:                         uint64(cStat.num_rd),
		Num_rd_kb:                      uint64(cStat.num_rd_kb),
		Num_wr:                         uint64(cStat.num_wr),
		Num_wr_kb:                      uint64(cStat.num_wr_kb),
	}, nil
}

// GetPoolID returns the pool ID associated with the I/O context.
//
// Implements:
//  int64_t rados_ioctx_get_id(rados_ioctx_t io)
func (ioctx *IOContext) GetPoolID() int64 {
	ret := C.rados_ioctx_get_id(ioctx.ioctx)
	return int64(ret)
}

// GetPoolName returns the name of the pool associated with the I/O context.
func (ioctx *IOContext) GetPoolName() (name string, err error) {
	var (
		buf []byte
		ret C.int
	)
	retry.WithSizes(128, 8192, func(size int) retry.Hint {
		buf = make([]byte, size)
		ret = C.rados_ioctx_get_pool_name(
			ioctx.ioctx,
			(*C.char)(unsafe.Pointer(&buf[0])),
			C.unsigned(len(buf)))
		err = getErrorIfNegative(ret)
		return retry.DoubleSize.If(err == errRange)
	})
	if err != nil {
		return "", err
	}
	name = C.GoStringN((*C.char)(unsafe.Pointer(&buf[0])), ret)
	return name, nil
}

// ObjectListFunc is the type of the function called for each object visited
// by ListObjects.
type ObjectListFunc func(oid string)

// ListObjects lists all of the objects in the pool associated with the I/O
// context, and called the provided listFn function for each object, passing
// to the function the name of the object. Call SetNamespace with
// RadosAllNamespaces before calling this function to return objects from all
// namespaces
func (ioctx *IOContext) ListObjects(listFn ObjectListFunc) error {
	pageResults := C.size_t(defaultListObjectsResultSize)
	var filterLen C.size_t
	results := make([]C.rados_object_list_item, pageResults)

	next := C.rados_object_list_begin(ioctx.ioctx)
	if next == nil {
		return ErrNotFound
	}
	defer C.rados_object_list_cursor_free(ioctx.ioctx, next)
	finish := C.rados_object_list_end(ioctx.ioctx)
	if finish == nil {
		return ErrNotFound
	}
	defer C.rados_object_list_cursor_free(ioctx.ioctx, finish)

	for {
		ret := C.rados_object_list(ioctx.ioctx, next, finish, pageResults, nil, filterLen, (*C.rados_object_list_item)(unsafe.Pointer(&results[0])), &next)
		if ret < 0 {
			return getError(ret)
		}

		numEntries := int(ret)
		for i := 0; i < numEntries; i++ {
			item := results[i]
			listFn(C.GoStringN(item.oid, (C.int)(item.oid_length)))
		}

		if C.rados_object_list_is_end(ioctx.ioctx, next) == listEndSentinel {
			return nil
		}
	}
}

// Stat returns the size of the object and its last modification time
func (ioctx *IOContext) Stat(object string) (stat ObjectStat, err error) {
	var cPsize C.uint64_t
	var cPmtime C.time_t
	cObject := C.CString(object)
	defer C.free(unsafe.Pointer(cObject))

	ret := C.rados_stat(
		ioctx.ioctx,
		cObject,
		&cPsize,
		&cPmtime)

	if ret < 0 {
		return ObjectStat{}, getError(ret)
	}
	return ObjectStat{
		Size:    uint64(cPsize),
		ModTime: time.Unix(int64(cPmtime), 0),
	}, nil
}

// GetXattr gets an xattr with key `name`, it returns the length of
// the key read or an error if not successful
func (ioctx *IOContext) GetXattr(object string, name string, data []byte) (int, error) {
	cObject := C.CString(object)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cObject))
	defer C.free(unsafe.Pointer(cName))

	ret := C.rados_getxattr(
		ioctx.ioctx,
		cObject,
		cName,
		(*C.char)(unsafe.Pointer(&data[0])),
		(C.size_t)(len(data)))

	if ret >= 0 {
		return int(ret), nil
	}
	return 0, getError(ret)
}

// SetXattr sets an xattr for an object with key `name` with value as `data`
func (ioctx *IOContext) SetXattr(object string, name string, data []byte) error {
	cObject := C.CString(object)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cObject))
	defer C.free(unsafe.Pointer(cName))

	ret := C.rados_setxattr(
		ioctx.ioctx,
		cObject,
		cName,
		(*C.char)(unsafe.Pointer(&data[0])),
		(C.size_t)(len(data)))

	return getError(ret)
}

// ListXattrs lists all the xattrs for an object. The xattrs are returned as a
// mapping of string keys and byte-slice values.
func (ioctx *IOContext) ListXattrs(oid string) (map[string][]byte, error) {
	coid := C.CString(oid)
	defer C.free(unsafe.Pointer(coid))

	var it C.rados_xattrs_iter_t

	ret := C.rados_getxattrs(ioctx.ioctx, coid, &it)
	if ret < 0 {
		return nil, getError(ret)
	}
	defer func() { C.rados_getxattrs_end(it) }()
	m := make(map[string][]byte)
	for {
		var cName, cVal *C.char
		var cLen C.size_t
		defer C.free(unsafe.Pointer(cName))
		defer C.free(unsafe.Pointer(cVal))

		ret := C.rados_getxattrs_next(it, &cName, &cVal, &cLen)
		if ret < 0 {
			return nil, getError(ret)
		}
		// rados api returns a null name,val & 0-length upon
		// end of iteration
		if cName == nil {
			return m, nil // stop iteration
		}
		m[C.GoString(cName)] = C.GoBytes(unsafe.Pointer(cVal), (C.int)(cLen))
	}
}

// RmXattr removes an xattr with key `name` from object `oid`
func (ioctx *IOContext) RmXattr(oid string, name string) error {
	coid := C.CString(oid)
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(coid))
	defer C.free(unsafe.Pointer(cName))

	ret := C.rados_rmxattr(
		ioctx.ioctx,
		coid,
		cName)

	return getError(ret)
}

// LockExclusive takes an exclusive lock on an object.
func (ioctx *IOContext) LockExclusive(oid, name, cookie, desc string, duration time.Duration, flags *byte) (int, error) {
	coid := C.CString(oid)
	cName := C.CString(name)
	cCookie := C.CString(cookie)
	cDesc := C.CString(desc)

	var cDuration C.struct_timeval
	if duration != 0 {
		tv := syscall.NsecToTimeval(duration.Nanoseconds())
		cDuration = C.struct_timeval{tv_sec: C.ceph_time_t(tv.Sec), tv_usec: C.ceph_suseconds_t(tv.Usec)}
	}

	var cFlags C.uint8_t
	if flags != nil {
		cFlags = C.uint8_t(*flags)
	}

	defer C.free(unsafe.Pointer(coid))
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cCookie))
	defer C.free(unsafe.Pointer(cDesc))

	ret := C.rados_lock_exclusive(
		ioctx.ioctx,
		coid,
		cName,
		cCookie,
		cDesc,
		&cDuration,
		cFlags)

	// 0 on success, negative error code on failure
	// -EBUSY if the lock is already held by another (client, cookie) pair
	// -EEXIST if the lock is already held by the same (client, cookie) pair

	switch ret {
	case 0:
		return int(ret), nil
	case -C.EBUSY:
		return int(ret), nil
	case -C.EEXIST:
		return int(ret), nil
	default:
		return int(ret), getError(ret)
	}
}

// LockShared takes a shared lock on an object.
func (ioctx *IOContext) LockShared(oid, name, cookie, tag, desc string, duration time.Duration, flags *byte) (int, error) {
	coid := C.CString(oid)
	cName := C.CString(name)
	cCookie := C.CString(cookie)
	cTag := C.CString(tag)
	cDesc := C.CString(desc)

	var cDuration C.struct_timeval
	if duration != 0 {
		tv := syscall.NsecToTimeval(duration.Nanoseconds())
		cDuration = C.struct_timeval{tv_sec: C.ceph_time_t(tv.Sec), tv_usec: C.ceph_suseconds_t(tv.Usec)}
	}

	var cFlags C.uint8_t
	if flags != nil {
		cFlags = C.uint8_t(*flags)
	}

	defer C.free(unsafe.Pointer(coid))
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cCookie))
	defer C.free(unsafe.Pointer(cTag))
	defer C.free(unsafe.Pointer(cDesc))

	ret := C.rados_lock_shared(
		ioctx.ioctx,
		coid,
		cName,
		cCookie,
		cTag,
		cDesc,
		&cDuration,
		cFlags)

	// 0 on success, negative error code on failure
	// -EBUSY if the lock is already held by another (client, cookie) pair
	// -EEXIST if the lock is already held by the same (client, cookie) pair

	switch ret {
	case 0:
		return int(ret), nil
	case -C.EBUSY:
		return int(ret), nil
	case -C.EEXIST:
		return int(ret), nil
	default:
		return int(ret), getError(ret)
	}
}

// Unlock releases a shared or exclusive lock on an object.
func (ioctx *IOContext) Unlock(oid, name, cookie string) (int, error) {
	coid := C.CString(oid)
	cName := C.CString(name)
	cCookie := C.CString(cookie)

	defer C.free(unsafe.Pointer(coid))
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cCookie))

	// 0 on success, negative error code on failure
	// -ENOENT if the lock is not held by the specified (client, cookie) pair

	ret := C.rados_unlock(
		ioctx.ioctx,
		coid,
		cName,
		cCookie)

	switch ret {
	case 0:
		return int(ret), nil
	case -C.ENOENT:
		return int(ret), nil
	default:
		return int(ret), getError(ret)
	}
}

// ListLockers lists clients that have locked the named object lock and
// information about the lock.
// The number of bytes required in each buffer is put in the corresponding size
// out parameter.  If any of the provided buffers are too short, -ERANGE is
// returned after these sizes are filled in.
func (ioctx *IOContext) ListLockers(oid, name string) (*LockInfo, error) {
	coid := C.CString(oid)
	cName := C.CString(name)

	cTag := (*C.char)(C.malloc(C.size_t(1024)))
	cClients := (*C.char)(C.malloc(C.size_t(1024)))
	cCookies := (*C.char)(C.malloc(C.size_t(1024)))
	cAddrs := (*C.char)(C.malloc(C.size_t(1024)))

	var cExclusive C.int
	cTagLen := C.size_t(1024)
	cClientsLen := C.size_t(1024)
	cCookiesLen := C.size_t(1024)
	cAddrsLen := C.size_t(1024)

	defer C.free(unsafe.Pointer(coid))
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cTag))
	defer C.free(unsafe.Pointer(cClients))
	defer C.free(unsafe.Pointer(cCookies))
	defer C.free(unsafe.Pointer(cAddrs))

	ret := C.rados_list_lockers(
		ioctx.ioctx,
		coid,
		cName,
		&cExclusive,
		cTag,
		&cTagLen,
		cClients,
		&cClientsLen,
		cCookies,
		&cCookiesLen,
		cAddrs,
		&cAddrsLen)

	splitCString := func(items *C.char, itemsLen C.size_t) []string {
		currLen := 0
		clients := []string{}
		for currLen < int(itemsLen) {
			client := C.GoString(C.nextChunk(&items))
			clients = append(clients, client)
			currLen += len(client) + 1
		}
		return clients
	}

	if ret < 0 {
		return nil, radosError(ret)
	}
	return &LockInfo{int(ret), cExclusive == 1, C.GoString(cTag), splitCString(cClients, cClientsLen), splitCString(cCookies, cCookiesLen), splitCString(cAddrs, cAddrsLen)}, nil
}

// BreakLock releases a shared or exclusive lock on an object, which was taken by the specified client.
func (ioctx *IOContext) BreakLock(oid, name, client, cookie string) (int, error) {
	coid := C.CString(oid)
	cName := C.CString(name)
	cClient := C.CString(client)
	cCookie := C.CString(cookie)

	defer C.free(unsafe.Pointer(coid))
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cClient))
	defer C.free(unsafe.Pointer(cCookie))

	// 0 on success, negative error code on failure
	// -ENOENT if the lock is not held by the specified (client, cookie) pair
	// -EINVAL if the client cannot be parsed

	ret := C.rados_break_lock(
		ioctx.ioctx,
		coid,
		cName,
		cClient,
		cCookie)

	switch ret {
	case 0:
		return int(ret), nil
	case -C.ENOENT:
		return int(ret), nil
	case -C.EINVAL: // -EINVAL
		return int(ret), nil
	default:
		return int(ret), getError(ret)
	}
}

// GetLastVersion will return the version number of the last object read or
// written to.
//
// Implements:
//  uint64_t rados_get_last_version(rados_ioctx_t io);
func (ioctx *IOContext) GetLastVersion() (uint64, error) {
	if err := ioctx.validate(); err != nil {
		return 0, err
	}
	v := C.rados_get_last_version(ioctx.ioctx)
	return uint64(v), nil
}

// GetNamespace gets the namespace used for objects within this IO context.
//
// Implements:
//  int rados_ioctx_get_namespace(rados_ioctx_t io, char *buf,
//                                unsigned maxlen);
func (ioctx *IOContext) GetNamespace() (string, error) {
	if err := ioctx.validate(); err != nil {
		return "", err
	}
	var (
		err error
		buf []byte
		ret C.int
	)
	retry.WithSizes(128, 8192, func(size int) retry.Hint {
		buf = make([]byte, size)
		ret = C.rados_ioctx_get_namespace(
			ioctx.ioctx,
			(*C.char)(unsafe.Pointer(&buf[0])),
			C.unsigned(len(buf)))
		err = getErrorIfNegative(ret)
		return retry.DoubleSize.If(err == errRange)
	})
	if err != nil {
		return "", err
	}
	return string(buf[:ret]), nil
}
