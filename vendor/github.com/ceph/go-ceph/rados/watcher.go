package rados

/*
#cgo LDFLAGS: -lrados
#include <stdlib.h>
#include <rados/librados.h>
extern void watchNotifyCb(void*, uint64_t, uint64_t, uint64_t, void*, size_t);
extern void watchErrorCb(void*, uint64_t, int);
*/
import "C"

import (
	"encoding/binary"
	"fmt"
	"math"
	"sync"
	"time"
	"unsafe"

	"github.com/ceph/go-ceph/internal/log"
)

type (
	// WatcherID is the unique id of a Watcher.
	WatcherID uint64
	// NotifyID is the unique id of a NotifyEvent.
	NotifyID uint64
	// NotifierID is the unique id of a notifying client.
	NotifierID uint64
)

// NotifyEvent is received by a watcher for each notification.
type NotifyEvent struct {
	ID         NotifyID
	WatcherID  WatcherID
	NotifierID NotifierID
	Data       []byte
}

// NotifyAck represents an acknowleged notification.
type NotifyAck struct {
	WatcherID  WatcherID
	NotifierID NotifierID
	Response   []byte
}

// NotifyTimeout represents an unacknowleged notification.
type NotifyTimeout struct {
	WatcherID  WatcherID
	NotifierID NotifierID
}

// Watcher receives all notifications for certain object.
type Watcher struct {
	id     WatcherID
	oid    string
	ioctx  *IOContext
	events chan NotifyEvent
	errors chan error
	done   chan struct{}
}

var (
	watchers    = map[WatcherID]*Watcher{}
	watchersMtx sync.RWMutex
)

// Watch creates a Watcher for the specified object.
//
// A Watcher receives all notifications that are sent to the object on which it
// has been created. It exposes two read-only channels: Events() receives all
// the NotifyEvents and Errors() receives all occuring errors. A typical code
// creating a Watcher could look like this:
//
//  watcher, err := ioctx.Watch(oid)
//  go func() { // event handler
//    for ne := range watcher.Events() {
//      ...
//      ne.Ack([]byte("response data..."))
//      ...
//    }
//  }()
//  go func() { // error handler
//    for err := range watcher.Errors() {
//      ... handle err ...
//    }
//  }()
//
// CAUTION: the Watcher references the IOContext in which it has been created.
// Therefore all watchers must be deleted with the Delete() method before the
// IOContext is being destroyed.
//
// Implements:
//  int rados_watch2(rados_ioctx_t io, const char* o, uint64_t* cookie,
//    rados_watchcb2_t watchcb, rados_watcherrcb_t watcherrcb, void* arg)
func (ioctx *IOContext) Watch(obj string) (*Watcher, error) {
	return ioctx.WatchWithTimeout(obj, 0)
}

// WatchWithTimeout creates a watcher on an object. Same as Watcher(), but
// different timeout than the default can be specified.
//
// Implements:
//  int rados_watch3(rados_ioctx_t io, const char *o, uint64_t *cookie,
// 	  rados_watchcb2_t watchcb, rados_watcherrcb_t watcherrcb, uint32_t timeout,
// 	  void *arg);
func (ioctx *IOContext) WatchWithTimeout(oid string, timeout time.Duration) (*Watcher, error) {
	cObj := C.CString(oid)
	defer C.free(unsafe.Pointer(cObj))
	var id C.uint64_t
	watchersMtx.Lock()
	defer watchersMtx.Unlock()
	ret := C.rados_watch3(
		ioctx.ioctx,
		cObj,
		&id,
		(C.rados_watchcb2_t)(C.watchNotifyCb),
		(C.rados_watcherrcb_t)(C.watchErrorCb),
		C.uint32_t(timeout.Milliseconds()/1000),
		nil,
	)
	if err := getError(ret); err != nil {
		return nil, err
	}
	evCh := make(chan NotifyEvent)
	errCh := make(chan error)
	w := &Watcher{
		id:     WatcherID(id),
		ioctx:  ioctx,
		oid:    oid,
		events: evCh,
		errors: errCh,
		done:   make(chan struct{}),
	}
	watchers[WatcherID(id)] = w
	return w, nil
}

// ID returns the WatcherId of the Watcher
func (w *Watcher) ID() WatcherID {
	return w.id
}

// Events returns a read-only channel, that receives all notifications that are
// sent to the object of the Watcher.
func (w *Watcher) Events() <-chan NotifyEvent {
	return w.events
}

// Errors returns a read-only channel, that receives all errors for the Watcher.
func (w *Watcher) Errors() <-chan error {
	return w.errors
}

// Check on the status of a Watcher.
//
// Returns the time since it was last confirmed. If there is an error, the
// Watcher is no longer valid, and should be destroyed with the Delete() method.
//
// Implements:
//  int rados_watch_check(rados_ioctx_t io, uint64_t cookie)
func (w *Watcher) Check() (time.Duration, error) {
	ret := C.rados_watch_check(w.ioctx.ioctx, C.uint64_t(w.id))
	if ret < 0 {
		return 0, getError(ret)
	}
	return time.Millisecond * time.Duration(ret), nil
}

// Delete the watcher. This closes both the event and error channel.
//
// Implements:
//  int rados_unwatch2(rados_ioctx_t io, uint64_t cookie)
func (w *Watcher) Delete() error {
	watchersMtx.Lock()
	_, ok := watchers[w.id]
	if ok {
		delete(watchers, w.id)
	}
	watchersMtx.Unlock()
	if !ok {
		return nil
	}
	ret := C.rados_unwatch2(w.ioctx.ioctx, C.uint64_t(w.id))
	if ret != 0 {
		return getError(ret)
	}
	close(w.done) // unblock blocked callbacks
	close(w.events)
	close(w.errors)
	return nil
}

// Notify sends a notification with the provided data to all Watchers of the
// specified object.
//
// CAUTION: even if the error is not nil. the returned slices
// might still contain data.
func (ioctx *IOContext) Notify(obj string, data []byte) ([]NotifyAck, []NotifyTimeout, error) {
	return ioctx.NotifyWithTimeout(obj, data, 0)
}

// NotifyWithTimeout is like Notify() but with a different timeout than the
// default.
//
// Implements:
//  int rados_notify2(rados_ioctx_t io, const char* o, const char* buf, int buf_len,
//    uint64_t timeout_ms, char** reply_buffer, size_t* reply_buffer_len)
func (ioctx *IOContext) NotifyWithTimeout(obj string, data []byte, timeout time.Duration) ([]NotifyAck,
	[]NotifyTimeout, error) {
	cObj := C.CString(obj)
	defer C.free(unsafe.Pointer(cObj))
	var cResponse *C.char
	defer C.rados_buffer_free(cResponse)
	var responseLen C.size_t
	var dataPtr *C.char
	if len(data) > 0 {
		dataPtr = (*C.char)(unsafe.Pointer(&data[0]))
	}
	ret := C.rados_notify2(
		ioctx.ioctx,
		cObj,
		dataPtr,
		C.int(len(data)),
		C.uint64_t(timeout.Milliseconds()),
		&cResponse,
		&responseLen,
	)
	// cResponse has been set even if an error is returned, so we decode it anyway
	acks, timeouts := decodeNotifyResponse(cResponse, responseLen)
	return acks, timeouts, getError(ret)
}

// Ack sends an acknowledgement with the specified response data to the notfier
// of the NotifyEvent. If a notify is not ack'ed, the originating Notify() call
// blocks and eventiually times out.
//
// Implements:
//  int rados_notify_ack(rados_ioctx_t io, const char *o, uint64_t notify_id,
//    uint64_t cookie, const char *buf, int buf_len)
func (ne *NotifyEvent) Ack(response []byte) error {
	watchersMtx.RLock()
	w, ok := watchers[ne.WatcherID]
	watchersMtx.RUnlock()
	if !ok {
		return fmt.Errorf("can't ack on deleted watcher %v", ne.WatcherID)
	}
	cOID := C.CString(w.oid)
	defer C.free(unsafe.Pointer(cOID))
	var respPtr *C.char
	if len(response) > 0 {
		respPtr = (*C.char)(unsafe.Pointer(&response[0]))
	}
	ret := C.rados_notify_ack(
		w.ioctx.ioctx,
		cOID,
		C.uint64_t(ne.ID),
		C.uint64_t(ne.WatcherID),
		respPtr,
		C.int(len(response)),
	)
	return getError(ret)
}

// WatcherFlush flushes all pending notifications of the cluster.
//
// Implements:
//  int rados_watch_flush(rados_t cluster)
func (c *Conn) WatcherFlush() error {
	if !c.connected {
		return ErrNotConnected
	}
	ret := C.rados_watch_flush(c.cluster)
	return getError(ret)
}

// decoder for this notify response format:
//    le32 num_acks
//    {
//      le64 gid     global id for the client (for client.1234 that's 1234)
//      le64 cookie  cookie for the client
//      le32 buflen  length of reply message buffer
//      u8 buflen  payload
//    } num_acks
//    le32 num_timeouts
//    {
//      le64 gid     global id for the client
//      le64 cookie  cookie for the client
//    } num_timeouts
//
// NOTE: starting with pacific this is implemented as a C function and this can
// be replaced later
func decodeNotifyResponse(response *C.char, len C.size_t) ([]NotifyAck, []NotifyTimeout) {
	if len == 0 || response == nil {
		return nil, nil
	}
	b := (*[math.MaxInt32]byte)(unsafe.Pointer(response))[:len:len]
	pos := 0

	num := binary.LittleEndian.Uint32(b[pos:])
	pos += 4
	acks := make([]NotifyAck, num)
	for i := range acks {
		acks[i].NotifierID = NotifierID(binary.LittleEndian.Uint64(b[pos:]))
		pos += 8
		acks[i].WatcherID = WatcherID(binary.LittleEndian.Uint64(b[pos:]))
		pos += 8
		dataLen := binary.LittleEndian.Uint32(b[pos:])
		pos += 4
		if dataLen > 0 {
			acks[i].Response = C.GoBytes(unsafe.Pointer(&b[pos]), C.int(dataLen))
			pos += int(dataLen)
		}
	}

	num = binary.LittleEndian.Uint32(b[pos:])
	pos += 4
	timeouts := make([]NotifyTimeout, num)
	for i := range timeouts {
		timeouts[i].NotifierID = NotifierID(binary.LittleEndian.Uint64(b[pos:]))
		pos += 8
		timeouts[i].WatcherID = WatcherID(binary.LittleEndian.Uint64(b[pos:]))
		pos += 8
	}
	return acks, timeouts
}

//export watchNotifyCb
func watchNotifyCb(_ unsafe.Pointer, notifyID C.uint64_t, id C.uint64_t,
	notifierID C.uint64_t, cData unsafe.Pointer, dataLen C.size_t) {
	ev := NotifyEvent{
		ID:         NotifyID(notifyID),
		WatcherID:  WatcherID(id),
		NotifierID: NotifierID(notifierID),
	}
	if dataLen > 0 {
		ev.Data = C.GoBytes(cData, C.int(dataLen))
	}
	watchersMtx.RLock()
	w, ok := watchers[WatcherID(id)]
	watchersMtx.RUnlock()
	if !ok {
		// usually this should not happen, but who knows
		log.Warnf("received notification for unknown watcher ID: %#v", ev)
		return
	}
	select {
	case <-w.done: // unblock when deleted
	case w.events <- ev:
	}
}

//export watchErrorCb
func watchErrorCb(_ unsafe.Pointer, id C.uint64_t, err C.int) {
	watchersMtx.RLock()
	w, ok := watchers[WatcherID(id)]
	watchersMtx.RUnlock()
	if !ok {
		// usually this should not happen, but who knows
		log.Warnf("received error for unknown watcher ID: id=%d err=%#v", id, err)
		return
	}
	select {
	case <-w.done: // unblock when deleted
	case w.errors <- getError(err):
	}
}
