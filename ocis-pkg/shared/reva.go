package shared

import "github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"

var defaultRevaConfig = Reva{
	Address: "127.0.0.1:9142",
}

// DefaultRevaConfig returns revas default config.
func DefaultRevaConfig() *Reva {
	ret := defaultRevaConfig
	return &ret
}

// GetRevaOptions
// FIXME: nolint
// nolint: revive
func (r *Reva) GetRevaOptions() []pool.Option {
	tm, _ := pool.StringToTLSMode(r.TLS.Mode)
	opts := []pool.Option{
		pool.WithTLSMode(tm),
	}
	return opts
}

// GetGRPCClientConfig
// FIXME: nolint
// nolint: revive
func (r *Reva) GetGRPCClientConfig() map[string]interface{} {
	return map[string]interface{}{
		"tls_mode":   r.TLS.Mode,
		"tls_cacert": r.TLS.CACert,
	}
}
