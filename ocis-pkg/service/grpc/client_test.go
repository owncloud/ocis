package grpc

import (
	"testing"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go-micro.dev/v4/client"
)

func TestWithDefaultCallOptionsAppliesGivenCallOptionsToDefaults(t *testing.T) {
	var applied int
	fake := client.CallOption(func(_ *client.CallOptions) {
		applied++
	})

	opt := withDefaultCallOptions(fake, fake)
	var opts client.Options
	opt(&opts)

	assert.Equal(t, 2, applied)
}

// TestNewClientKeepaliveDialOptions ensures that keepalive params configured via
// WithKeepaliveParams end up in the client's *default* CallOptions.Context, and only then. That
// context is where the go-micro grpc plugin looks up the dial options on every dial (see
// grpcClient.getGrpcDialOptions in github.com/go-micro/plugins/v4/client/grpc, which reads
// g.opts.CallOptions.Context and ignores per-call CallOptions). Without them a black-holed peer
// would hang the caller for as long as it allows instead of being detected by the keepalive ping.
//
// The dial options themselves cannot be asserted on: the plugin stores them under an unexported
// context key, so a non-nil context is as close as this package can get from the outside.
func TestNewClientKeepaliveDialOptions(t *testing.T) {
	// the plugin needs a registry to construct a client, and an in-memory one keeps the test
	// from reaching out to a nats-js-kv that isn't there
	t.Setenv("MICRO_REGISTRY", "memory")
	// an empty address is treated like an unset one by ocis-pkg/registry
	t.Setenv("MICRO_REGISTRY_ADDRESS", "")

	tests := map[string]struct {
		opts           []ClientOption
		wantContextSet bool
	}{
		"keepalive configured": {
			opts: []ClientOption{
				WithKeepaliveParams(&shared.GRPCClientOptions{GRPCClientKeepaliveTime: 30 * time.Second}),
			},
			wantContextSet: true,
		},
		"no keepalive option leaves grpc's own dial defaults untouched": {
			wantContextSet: false,
		},
		"keepalive option without a configured time leaves grpc's own dial defaults untouched": {
			opts:           []ClientOption{WithKeepaliveParams(&shared.GRPCClientOptions{})},
			wantContextSet: false,
		},
		"keepalive option without any options leaves grpc's own dial defaults untouched": {
			opts:           []ClientOption{WithKeepaliveParams(nil)},
			wantContextSet: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c, err := NewClient(append(GetClientOptions(&shared.GRPCClientTLS{Mode: "off"}), tt.opts...)...)
			require.NoError(t, err)

			if tt.wantContextSet {
				assert.NotNil(t, c.Options().CallOptions.Context)
				return
			}
			assert.Nil(t, c.Options().CallOptions.Context)
		})
	}
}
