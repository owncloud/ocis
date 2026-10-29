package debug

import (
	"context"
	"net/http"

	"github.com/owncloud/ocis/v2/ocis-pkg/checks"
	"github.com/owncloud/ocis/v2/ocis-pkg/handlers"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/reva/v2/pkg/events/stream"
)

// Server initializes the debug service and server.
func Server(opts ...Option) (*http.Server, error) {
	options := newOptions(opts...)

	readyHandlerConfiguration := handlers.NewCheckHandlerConfiguration().
		WithLogger(options.Logger).
		WithCheck("nats reachability", natsCheck(options))

	return debug.NewService(
		debug.Logger(options.Logger),
		debug.Name(options.Config.Service.Name),
		debug.Version(version.GetString()),
		debug.Address(options.Config.Debug.Addr),
		debug.Token(options.Config.Debug.Token),
		debug.Pprof(options.Config.Debug.Pprof),
		debug.Zpages(options.Config.Debug.Zpages),
		debug.Ready(handlers.NewCheckHandler(readyHandlerConfiguration)),
		//debug.CorsAllowedOrigins(options.Config.HTTP.CORS.AllowedOrigins),
		//debug.CorsAllowedMethods(options.Config.HTTP.CORS.AllowedMethods),
		//debug.CorsAllowedHeaders(options.Config.HTTP.CORS.AllowedHeaders),
		//debug.CorsAllowCredentials(options.Config.HTTP.CORS.AllowCredentials),
	), nil
}

// natsCheck returns the check of the provider cache. When the provider cache does not
// use NATS, the check always succeeds.
func natsCheck(options Options) func(context.Context) error {
	c := options.Config.Cache
	if c.ProviderCacheStore != "nats-js-kv" || len(c.ProviderCacheNodes) == 0 {
		return func(context.Context) error { return nil }
	}
	return checks.NewNatsCheckWithConfig(c.ProviderCacheNodes[0], stream.NatsConfig{
		TLSInsecure:          c.ProviderCacheTLSInsecure,
		TLSRootCACertificate: c.ProviderCacheTLSRootCACertificate,
		EnableTLS:            c.ProviderCacheEnableTLS,
		AuthUsername:         c.ProviderCacheAuthUsername,
		AuthPassword:         c.ProviderCacheAuthPassword,
		CredsFile:            c.ProviderCacheCredsFile,
		Creds:                c.ProviderCacheCreds,
	})
}
