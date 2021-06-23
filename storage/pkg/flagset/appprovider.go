package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// AppProviderWithConfig applies cfg to the root flagset
func AppProviderWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{

		// debug ports are the odd ports
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppProvider.DebugAddr, "0.0.0.0:9165"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"APP_PROVIDER_BASIC_DEBUG_ADDR"},
			Destination: &cfg.Reva.AppProvider.DebugAddr,
		},

		// Auth

		// Services

		// AppProvider

		&cli.StringFlag{
			Name:        "network",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppProvider.GRPCNetwork, "tcp"),
			Usage:       "Network to use for the storage auth-basic service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"APP_PROVIDER_BASIC_GRPC_NETWORK"},
			Destination: &cfg.Reva.AppProvider.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppProvider.GRPCAddr, "0.0.0.0:9164"),
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"APP_PROVIDER_BASIC_GRPC_ADDR"},
			Destination: &cfg.Reva.AppProvider.GRPCAddr,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("appprovider"),
			Usage:   "--service appprovider [--service otherservice]",
			EnvVars: []string{"APP_PROVIDER_BASIC_SERVICES"},
		},
		&cli.StringFlag{
			Name:        "driver",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppProvider.Driver, "demo"),
			Usage:       "app provider driver",
			EnvVars:     []string{"APP_PROVIDER_DRIVER"},
			Destination: &cfg.Reva.AppProvider.Driver,
		},
		&cli.StringFlag{
			Name:        "iopsecret",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppProvider.IopSecret, ""),
			Usage:       "IOP Secret (Shared with WOPI server)",
			EnvVars:     []string{"APP_PROVIDER_IOP_SECRET"},
			Destination: &cfg.Reva.AppProvider.IopSecret,
		},
		&cli.BoolFlag{
			Name:        "wopiinsecure",
			Value:       flags.OverrideDefaultBool(cfg.Reva.AppProvider.WopiInsecure, false),
			Usage:       "Disable SSL certificate verification of WOPI server and WOPI bridge",
			EnvVars:     []string{"APP_PROVIDER_WOPI_INSECURE"},
			Destination: &cfg.Reva.AppProvider.WopiInsecure,
		},
		&cli.StringFlag{
			Name:        "wopiurl",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppProvider.WopiUrl, ""),
			Usage:       "WOPI server URL",
			EnvVars:     []string{"APP_PROVIDER_WOPI_URL"},
			Destination: &cfg.Reva.AppProvider.WopiUrl,
		},
		&cli.StringFlag{
			Name:        "wopibridgeurl",
			Value:       flags.OverrideDefaultString(cfg.Reva.AppProvider.WopiBridgeUrl, ""),
			Usage:       "WOPI bridge URL",
			EnvVars:     []string{"APP_PROVIDER_WOPI_BRIDGE_URL"},
			Destination: &cfg.Reva.AppProvider.WopiBridgeUrl,
		},

		// Gateway

		&cli.StringFlag{
			Name:        "gateway-url",
			Value:       flags.OverrideDefaultString(cfg.Reva.Gateway.Endpoint, "localhost:9142"),
			Usage:       "URL to use for the storage gateway service",
			EnvVars:     []string{"STORAGE_GATEWAY_ENDPOINT"},
			Destination: &cfg.Reva.Gateway.Endpoint,
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)

	return flags
}
