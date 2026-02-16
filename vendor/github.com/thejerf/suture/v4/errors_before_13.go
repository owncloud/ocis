// +build !go1.13

package suture

import "errors"

func isErr(err error, target error) bool {
	return err == target
}

// ErrDoNotRestart can be returned by a service to voluntarily not
// be restarted.
var ErrDoNotRestart = errors.New("service should not be restarted")

// ErrTerminateSupervisorTree can can be returned by a service to terminate the
// entire supervision tree above it as well.
var ErrTerminateSupervisorTree = errors.New("tree should be terminated")
