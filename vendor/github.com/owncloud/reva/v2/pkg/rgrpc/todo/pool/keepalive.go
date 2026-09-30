package pool

import (
	"os"
	"time"

	"google.golang.org/grpc/keepalive"
)

const (
	_clientKeepaliveTimeEnv    = "GRPC_CLIENT_KEEPALIVE_TIME"
	_clientKeepaliveTimeoutEnv = "GRPC_CLIENT_KEEPALIVE_TIMEOUT"

	// Used when a keepalive time is configured but the timeout is not usable, so
	// that a half configured setup still gives up on a dead peer rather than
	// sitting out grpc's longer stock timeout of 20s.
	_defaultKeepaliveTimeout = 10 * time.Second
)

// GetClientKeepaliveParams returns the keepalive parameters for every grpc
// client connection of the pool, read from GRPC_CLIENT_KEEPALIVE_TIME and
// GRPC_CLIENT_KEEPALIVE_TIMEOUT.
//
// Once they are set, a peer that stops answering on an established connection -
// a black-holed node, a wedged process - can be told apart from one that is
// merely slow: grpc pings it while an rpc is in flight and fails the requests on
// that connection when no answer arrives within the timeout, instead of waiting
// for as long as the caller allows. The pings are only sent while an rpc is in
// flight, so a healthy server never sees them often enough to complain.
//
// Keepalive is strictly opt-in and driven by the keepalive time alone. Without a
// usable time this returns the zero ClientParameters, which NewConn reads as
// "add no keepalive dial option at all", so grpc's own behavior applies and no
// pings are sent - a deployment that configured nothing behaves exactly as it
// did before this package started setting keepalive params. A time with no
// usable timeout gets a working default instead, so detection cannot be left
// half configured.
//
// Nothing is salvaged or second-guessed: a value that is empty, negative or not
// a valid duration - "20" without a unit suffix - reads the same as unset. ocis
// configures its go-micro based clients from the same two env vars and applies
// the same rules, so both sets of clients behave identically for every input.
// grpc also raises any time below 10s to 10s on its own.
func GetClientKeepaliveParams() keepalive.ClientParameters {
	keepaliveTime := duration(_clientKeepaliveTimeEnv)
	if keepaliveTime == 0 {
		// not configured: the zero value tells NewConn to skip the keepalive
		// dial option, so grpc's own behavior applies and no pings are sent
		return keepalive.ClientParameters{}
	}
	keepaliveTimeout := duration(_clientKeepaliveTimeoutEnv)
	if keepaliveTimeout <= 0 {
		// a time without a usable timeout: fall back to a working default
		keepaliveTimeout = _defaultKeepaliveTimeout
	}
	return keepalive.ClientParameters{
		// If set below 10s, the grpc library will set a minimum value of 10s.
		Time:    keepaliveTime,
		Timeout: keepaliveTimeout,
		// no pings on connections that have no rpc in flight: an idle
		// connection has nothing to rescue, and staying quiet keeps us clear of
		// any server side ping enforcement
		PermitWithoutStream: false,
	}
}

// duration reads a duration from the environment. Anything that is not a usable
// positive duration - unset, empty, negative or missing its unit suffix - reads
// as zero, which every caller treats as "not configured".
func duration(env string) time.Duration {
	v, ok := os.LookupEnv(env)
	if !ok {
		return 0
	}

	d, err := time.ParseDuration(v)
	if err != nil || d < 0 {
		return 0
	}
	return d
}
