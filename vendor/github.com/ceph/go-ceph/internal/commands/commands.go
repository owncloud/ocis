package commands

import (
	"encoding/json"

	ccom "github.com/ceph/go-ceph/common/commands"
	"github.com/ceph/go-ceph/rados"
)

func validate(m interface{}) error {
	if m == nil {
		return rados.ErrNotConnected
	}
	return nil
}

// RawMgrCommand takes a byte buffer and sends it to the MGR as a command.
// The buffer is expected to contain preformatted JSON.
func RawMgrCommand(m ccom.MgrCommander, buf []byte) Response {
	if err := validate(m); err != nil {
		return Response{err: err}
	}
	return NewResponse(m.MgrCommand([][]byte{buf}))
}

// MarshalMgrCommand takes an generic interface{} value, converts it to JSON
// and sends the json to the MGR as a command.
func MarshalMgrCommand(m ccom.MgrCommander, v interface{}) Response {
	b, err := json.Marshal(v)
	if err != nil {
		return Response{err: err}
	}
	return RawMgrCommand(m, b)
}

// RawMonCommand takes a byte buffer and sends it to the MON as a command.
// The buffer is expected to contain preformatted JSON.
func RawMonCommand(m ccom.MonCommander, buf []byte) Response {
	if err := validate(m); err != nil {
		return Response{err: err}
	}
	return NewResponse(m.MonCommand(buf))
}

// MarshalMonCommand takes an generic interface{} value, converts it to JSON
// and sends the json to the MGR as a command.
func MarshalMonCommand(m ccom.MonCommander, v interface{}) Response {
	b, err := json.Marshal(v)
	if err != nil {
		return Response{err: err}
	}
	return RawMonCommand(m, b)
}
