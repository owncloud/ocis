package command

import (
    "github.com/cernbox/ocis-canary/pkg/command"
    svcconfig "github.com/cernbox/ocis-canary/pkg/config"
    "github.com/cernbox/ocis-canary/pkg/flagset"
    "github.com/micro/cli/v2"
    "github.com/owncloud/ocis/ocis-pkg/config"
    "github.com/owncloud/ocis/ocis/pkg/register"
)

// OcisCanaryCommand is the entry point for the settings command.
func OcisCanaryCommand(cfg *config.Config) *cli.Command {
    return &cli.Command{
        Name:     "ocis-canary",
        Usage:    "Starts CERNBox Canary for OCIS",
        Category: "Canary",
        Flags:    flagset.ServerWithConfig(cfg.CernboxCanary),
        Action: func(c *cli.Context) error {
            origCmd := command.Server(configureCanary(cfg))
            return handleOriginalAction(c, origCmd)
        },
    }
}

func configureCanary(cfg *config.Config) *svcconfig.Config {
    cfg.CernboxCanary.Log.Level = cfg.Log.Level
    cfg.CernboxCanary.Log.Pretty = cfg.Log.Pretty
    cfg.CernboxCanary.Log.Color = cfg.Log.Color

    if cfg.Tracing.Enabled {
        cfg.CernboxCanary.Tracing.Enabled = cfg.Tracing.Enabled
        cfg.CernboxCanary.Tracing.Type = cfg.Tracing.Type
        cfg.CernboxCanary.Tracing.Endpoint = cfg.Tracing.Endpoint
        cfg.CernboxCanary.Tracing.Collector = cfg.Tracing.Collector
    }

    if cfg.TokenManager.JWTSecret != "" {
        cfg.CernboxCanary.TokenManager.JWTSecret = cfg.TokenManager.JWTSecret
    }

    return cfg.CernboxCanary
}

func init() {
    register.AddCommand(OcisCanaryCommand)
}
