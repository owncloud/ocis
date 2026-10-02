package grpc

import (
	"testing"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/config/envdecode"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// must match the env tags on shared.GRPCClientOptions. GetClientKeepaliveParams itself reads no
// environment, so these are only needed to exercise the decoding the services do before calling it.
const (
	_clientKeepaliveTimeEnv    = "GRPC_CLIENT_KEEPALIVE_TIME"
	_clientKeepaliveTimeoutEnv = "GRPC_CLIENT_KEEPALIVE_TIMEOUT"
)

func TestGetClientKeepaliveParams(t *testing.T) {
	tests := map[string]struct {
		time, timeout       time.Duration
		wantTime, wantTimeo time.Duration
	}{
		"nothing configured disables keepalive entirely": {
			wantTime:  0,
			wantTimeo: 0,
		},
		"nothing configured disables keepalive even if the timeout is set": {
			timeout:   5 * time.Second,
			wantTime:  0,
			wantTimeo: 0,
		},
		"a negative time disables keepalive as well": {
			time:      -5 * time.Second,
			timeout:   5 * time.Second,
			wantTime:  0,
			wantTimeo: 0,
		},
		"valid values are used as configured": {
			time:      30 * time.Second,
			timeout:   5 * time.Second,
			wantTime:  30 * time.Second,
			wantTimeo: 5 * time.Second,
		},
		"a time without a timeout falls back to the working default": {
			time:      30 * time.Second,
			wantTime:  30 * time.Second,
			wantTimeo: 10 * time.Second,
		},
		"a time with a negative timeout falls back to the working default": {
			time:      30 * time.Second,
			timeout:   -5 * time.Second,
			wantTime:  30 * time.Second,
			wantTimeo: 10 * time.Second,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			params := GetClientKeepaliveParams(tt.time, tt.timeout)

			assert.Equal(t, tt.wantTime, params.Time)
			assert.Equal(t, tt.wantTimeo, params.Timeout)
			assert.False(t, params.PermitWithoutStream)
		})
	}
}

// TestGetClientKeepaliveParamsAfterEnvdecode drives the two env vars through the same decoder the
// services use, to pin down what the supported env var spellings actually produce. Only fields
// tagged 'strict' make envdecode report a value it cannot parse, so a duration without a unit
// suffix is dropped and reaches GetClientKeepaliveParams as a zero, which this package treats as
// "not configured". reva's pool reads the same two env vars for the CS3 clients and arrives at the
// same result for these inputs, so the two stay in step.
func TestGetClientKeepaliveParamsAfterEnvdecode(t *testing.T) {
	tests := map[string]struct {
		env                 map[string]string
		wantTime, wantTimeo time.Duration
	}{
		"parseable values are honored": {
			env: map[string]string{
				_clientKeepaliveTimeEnv:    "45s",
				_clientKeepaliveTimeoutEnv: "7s",
			},
			wantTime:  45 * time.Second,
			wantTimeo: 7 * time.Second,
		},
		"a time without a unit suffix is dropped and reads as not configured": {
			env:       map[string]string{_clientKeepaliveTimeEnv: "20"},
			wantTime:  0,
			wantTimeo: 0,
		},
		"an empty time reads as not configured": {
			env:       map[string]string{_clientKeepaliveTimeEnv: ""},
			wantTime:  0,
			wantTimeo: 0,
		},
		"a timeout without a unit suffix falls back but keeps the configured time": {
			env: map[string]string{
				_clientKeepaliveTimeEnv:    "45s",
				_clientKeepaliveTimeoutEnv: "7",
			},
			wantTime:  45 * time.Second,
			wantTimeo: 10 * time.Second,
		},
		"a timeout alone does not enable keepalive": {
			env:       map[string]string{_clientKeepaliveTimeoutEnv: "7s"},
			wantTime:  0,
			wantTimeo: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			var co shared.GRPCClientOptions
			require.NoError(t, envdecode.Decode(&co))

			params := GetClientKeepaliveParams(co.GRPCClientKeepaliveTime, co.GRPCClientKeepaliveTimeout)
			assert.Equal(t, tt.wantTime, params.Time)
			assert.Equal(t, tt.wantTimeo, params.Timeout)
		})
	}
}
