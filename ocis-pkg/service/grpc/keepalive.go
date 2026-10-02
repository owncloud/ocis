package grpc

import (
	"time"

	"google.golang.org/grpc/keepalive"
)

// GetClientKeepaliveParams turns the configured GRPCClientKeepaliveTime/-Timeout into grpc
// keepalive client parameters for the go-micro based clients built in this package. Once they are
// set, a peer that stops answering on an established connection can be told apart from one that is
// merely slow: grpc pings it while an rpc is in flight and fails the requests on that connection
// when no answer arrives within the timeout, instead of waiting for as long as the caller allows.
//
// Keepalive is strictly opt-in and driven by the keepalive time alone. Without a positive time this
// returns the zero ClientParameters, which NewClient reads as "add no keepalive dial option at
// all", so grpc's own behavior applies and no pings are sent - a deployment that configured nothing
// behaves exactly as it did before this package started setting keepalive params. A positive time
// with no usable timeout gets a working default instead, so detection cannot be left half
// configured.
//
// Values are taken as configured; nothing is salvaged or second-guessed. A malformed value - '20'
// without a unit suffix, say - is dropped by ocis-pkg/config/envdecode before it reaches this
// function and therefore reads as "not configured" here. reva's pool.GetClientKeepaliveParams
// (https://github.com/owncloud/reva/blob/master/pkg/rgrpc/todo/pool/keepalive.go) parses the same
// two env vars itself for the CS3 clients and applies the same rules, so both sets of clients
// behave identically for every input. grpc also raises any time below '10s' to '10s' on its own.
func GetClientKeepaliveParams(grpcClientKeepaliveTime, grpcClientKeepaliveTimeout time.Duration) keepalive.ClientParameters {
	// used when a keepalive time is configured but the timeout is not usable, so that a half
	// configured setup still gives up on a dead peer rather than sitting out grpc's longer
	// stock timeout of 20s
	const _defaultKeepaliveTimeout = 10 * time.Second

	if grpcClientKeepaliveTime <= 0 {
		// not configured: the zero value tells NewClient to skip the keepalive dial option, so
		// grpc's own behavior applies and no pings are sent
		return keepalive.ClientParameters{}
	}
	if grpcClientKeepaliveTimeout <= 0 {
		// a time without a usable timeout: fall back to a working default
		grpcClientKeepaliveTimeout = _defaultKeepaliveTimeout
	}

	return keepalive.ClientParameters{
		Time:    grpcClientKeepaliveTime,
		Timeout: grpcClientKeepaliveTimeout,
		// no pings on connections that have no rpc in flight: an idle connection has nothing
		// to rescue, and staying quiet keeps us clear of any server side ping enforcement
		PermitWithoutStream: false,
	}
}
