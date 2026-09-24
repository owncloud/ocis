package pool

import (
	"math"
	"os"
	"time"

	"google.golang.org/grpc/keepalive"
)

const (
	_clientKeepaliveTimeEnv    = "GRPC_CLIENT_KEEPALIVE_TIME"
	_clientKeepaliveTimeoutEnv = "GRPC_CLIENT_KEEPALIVE_TIMEOUT"

	// How long a connection with an active rpc may stay silent before we ping
	// the peer, and how long we then wait for the pong. grpc raises anything
	// below 10s to 10s, so 20s is the smallest value that is not silently
	// rewritten while leaving room for a slow but living peer.
	_defaultKeepaliveTime    = 20 * time.Second
	_defaultKeepaliveTimeout = 10 * time.Second

	// same value grpc uses to mean "never"
	_keepaliveDisabled = time.Duration(math.MaxInt64)
)

// GetClientKeepaliveParams returns the keepalive parameters for every grpc
// client connection of the pool.
//
// Without them a peer that stops answering on an established connection - a
// black-holed node, a wedged process - is indistinguishable from a peer that
// is merely slow, and every rpc on that connection waits forever. The pings
// are only sent while an rpc is in flight, so a healthy server never sees them
// often enough to complain.
//
// Both values can be overridden with a duration ("30s"); an empty, negative or
// unparseable value yields the default, never "never". Setting the time to "0"
// disables the pings altogether - it cannot mean a zero interval, because grpc
// would raise that to its 10s minimum.
func GetClientKeepaliveParams() keepalive.ClientParameters {
	return keepalive.ClientParameters{
		Time:    keepaliveDuration(_clientKeepaliveTimeEnv, _defaultKeepaliveTime),
		Timeout: keepaliveDuration(_clientKeepaliveTimeoutEnv, _defaultKeepaliveTimeout),
		// no pings on connections that have no rpc in flight: an idle
		// connection has nothing to rescue, and staying quiet keeps us clear of
		// any server side ping enforcement
		PermitWithoutStream: false,
	}
}

func keepaliveDuration(env string, fallback time.Duration) time.Duration {
	v := os.Getenv(env)
	if v == "0" {
		return _keepaliveDisabled
	}

	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return fallback
	}
	return d
}
