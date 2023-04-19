/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */

/*
 This package implements system calls, not implemented by the syscall package.
 */
package syscall_x

import "syscall"
import "unsafe"

func Fgetxattr(fd int, attr string, dest []byte) (sz int, err error) {
	attr2 , err := syscall.BytePtrFromString(attr)
	destp := uintptr(0)
	destl := uintptr(len(dest))
	if destl>0 { destp = uintptr(unsafe.Pointer(&dest[0])) }
	if err!=nil { return 0,err }
	sz_,_,err := syscall.Syscall6(
			syscall.SYS_FGETXATTR, uintptr(fd),
			uintptr(unsafe.Pointer(attr2)),
			destp,
			destl,
		0, 0)
	if err==syscall.Errno(0) { err = nil }
	return int(sz_),err
}

func Fsetxattr(fd int, attr string, dest []byte,flags int) error {
	attr2 , err := syscall.BytePtrFromString(attr)
	destp := uintptr(0)
	destl := uintptr(len(dest))
	if destl>0 { destp = uintptr(unsafe.Pointer(&dest[0])) }
	if err!=nil { return err }
	_,_,err = syscall.Syscall6(
			syscall.SYS_FSETXATTR, uintptr(fd),
			uintptr(unsafe.Pointer(attr2)),
			destp,
			destl,
		uintptr(flags), 0)
	if err==syscall.Errno(0) { err = nil }
	return err
}

type winsize struct{
	ws_row 		uint16	/* rows, in characters */
	ws_col 		uint16	/* columns, in characters */
	ws_xpixel	uint16	/* horizontal size, pixels */
	ws_ypixel 	uint16	/* vertical size, pixels */
}

func clamp16(i int) uint16 {
	if i<0 { return 0 }
	if i>=(1<<16) { return ^uint16(0) }
	return uint16(i)
}

/*
 Does (C++):

 winsize ws;

 ws.ws_col = w;

 ws.ws_row = h;

 ws.ws_xpixel = w*5;

 ws.ws_ypixel = h*10;

 ioctl(fd,TIOCSWINSZ,&ws);
*/
func Ioctl_resize(fd int,w, h int) {
	var ws winsize
	ws.ws_col = clamp16(w)
	ws.ws_row = clamp16(h)
	ws.ws_xpixel = clamp16(w*5)
	ws.ws_ypixel = clamp16(h*10)
	syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		syscall.TIOCSWINSZ,
		uintptr(unsafe.Pointer(&ws)) )
}



