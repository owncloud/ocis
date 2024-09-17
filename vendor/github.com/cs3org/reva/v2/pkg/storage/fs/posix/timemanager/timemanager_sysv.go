//go:build dragonfly || linux || solaris

package timemanager

import "syscall"

// StatCtime returns the change time
func StatCTime(st *syscall.Stat_t) syscall.Timespec {
	return st.Ctim
}
