package cutil

/*
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

// CommandInput can be used to manage the input to ceph's *_command functions.
type CommandInput struct {
	cmd   []*C.char
	inbuf []byte
}

// NewCommandInput creates C-level pointers from go byte buffers suitable
// for passing off to ceph's *_command functions.
func NewCommandInput(cmd [][]byte, inputBuffer []byte) *CommandInput {
	ci := &CommandInput{
		cmd:   make([]*C.char, len(cmd)),
		inbuf: inputBuffer,
	}
	for i := range cmd {
		ci.cmd[i] = C.CString(string(cmd[i]))
	}
	return ci
}

// Free any C memory managed by this CommandInput.
func (ci *CommandInput) Free() {
	for i := range ci.cmd {
		C.free(unsafe.Pointer(ci.cmd[i]))
	}
	ci.cmd = nil
}

// Cmd returns an unsafe wrapper around an array of C-strings.
func (ci *CommandInput) Cmd() CharPtrPtr {
	ptr := &ci.cmd[0]
	return CharPtrPtr(ptr)
}

// CmdLen returns the length of the array of C-strings returned by
// Cmd.
func (ci *CommandInput) CmdLen() SizeT {
	return SizeT(len(ci.cmd))
}

// InBuf returns an unsafe wrapper to a C char*.
func (ci *CommandInput) InBuf() CharPtr {
	if len(ci.inbuf) == 0 {
		return nil
	}
	return CharPtr(&ci.inbuf[0])
}

// InBufLen returns the length of the buffer returned by InBuf.
func (ci *CommandInput) InBufLen() SizeT {
	return SizeT(len(ci.inbuf))
}
