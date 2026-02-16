// +build go1.13

package suture

import "errors"

func isErr(err error, target error) bool {
	return errors.Is(err, target)
}

// ErrDoNotRestart can be returned by a service to voluntarily not
// be restarted. Any error that will compare with errors.Is as being this
// error will count as an ErrDoNotRestart.
var ErrDoNotRestart = errors.New("service should not be restarted")

// ErrTerminateSupervisorTree can can be returned by a service to terminate the
// entire supervision tree above it as well. Any error that will compare
// with errors.Is to be ErrTerminateSupervisorTree will count as an
// ErrTerminateSupervisorTree.
var ErrTerminateSupervisorTree = errors.New("tree should be terminated")
