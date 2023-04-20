// Package log is the internal package for go-ceph logging. This package is only
// used from go-ceph code, not from consumers of go-ceph. go-ceph code uses the
// functions in this package to log information that can't be returned as
// errors. The functions default to no-ops and can be set with the external log
// package common/log by the go-ceph consumers.
package log

func noop(string, ...interface{}) {}

// These variables are set by the common log package.
var (
	Warnf  = noop
	Debugf = noop
)
