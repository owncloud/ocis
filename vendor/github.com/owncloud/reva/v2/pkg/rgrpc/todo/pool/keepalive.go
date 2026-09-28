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

	// Used when GRPC_CLIENT_KEEPALIVE_TIME / GRPC_CLIENT_KEEPALIVE_TIMEOUT are
	// not present in the environment at all - i.e. the operator hasn't opted
	// into this package's keepalive behavior. These equal the zero-value
	// defaults grpc itself falls back to (google.golang.org/grpc's
	// internal/transport.defaultClientKeepaliveTime/-Timeout), so an
	// unconfigured deployment behaves exactly as it did before this package
	// started setting keepalive params at all.
	_clientDefaultKeepaliveTime    = time.Duration(math.MaxInt64)
	_clientDefaultKeepaliveTimeout = time.Duration(20 * time.Second)

	// Used when the environment variable is present but its value is empty,
	// negative or unparseable - i.e. a configuration mistake rather than a
	// deliberate opt-out. How long a connection with an active rpc may stay
	// silent before we ping the peer, and how long we then wait for the pong.
	// grpc raises anything below 10s to 10s, so 20s is the smallest value that
	// is not silently rewritten while leaving room for a slow but living peer.
	_defaultKeepaliveTime    = 20 * time.Second
	_defaultKeepaliveTimeout = 10 * time.Second
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
// Both values can be overridden with a duration ("30s"). If the environment
// variable is not set at all, the operator hasn't opted in, so grpc's own
// stock ClientParameters zero value is used - which for Time means no pings
// are sent at all. If the variable is set but the value is empty, negative or
// unparseable, that is treated as a configuration mistake rather than an
// explicit opt-out, so it falls back to this package's own working default
// instead of silently disabling detection.
func GetClientKeepaliveParams() keepalive.ClientParameters {
	return keepalive.ClientParameters{
		Time:    keepaliveDuration(_clientKeepaliveTimeEnv, _defaultKeepaliveTime, _clientDefaultKeepaliveTime),
		Timeout: keepaliveDuration(_clientKeepaliveTimeoutEnv, _defaultKeepaliveTimeout, _clientDefaultKeepaliveTimeout),
		// no pings on connections that have no rpc in flight: an idle
		// connection has nothing to rescue, and staying quiet keeps us clear of
		// any server side ping enforcement
		PermitWithoutStream: false,
	}
}

func keepaliveDuration(env string, fallback, clientDefault time.Duration) time.Duration {
	v, ok := os.LookupEnv(env)
	if !ok {
		return clientDefault
	}

	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return fallback
	}
	return d
}
