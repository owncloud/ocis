package http

import (
	"crypto/tls"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/web"
)

// Service simply wraps the go-micro web service.
type Service struct {
	web.Service
}

var (
	Once sync.Once
	M    sync.Mutex // protects flag setting from competing goroutines
)

// NewService initializes a new http service.
func NewService(opts ...Option) Service {
	sopts := newOptions(opts...)
	sopts.Logger.Info().
		Str("transport", transport(sopts.TLSConfig)).
		Str("addr", sopts.Address).
		Msg("starting server")

	M.Lock()

	flags := make([]cli.Flag, 0)
	// There is a global state somewhere in go-micro. In order to circumvent this, flag registration down to the go-micro
	// flags need to be parsed ONLY once, in other words, this roughly translates into: let go-micro know we're using
	// our own set of flags, so it does not panic when it encounters a flag it does not recognise.
	Once.Do(func() {
		flags = []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Value:   "info",
				Usage:   "Set logging level",
				EnvVars: []string{"OCIS_LOG_LEVEL"},
			},
			&cli.BoolFlag{
				Value:   false,
				Name:    "log-pretty",
				Usage:   "Enable pretty logging",
				EnvVars: []string{"OCIS_LOG_PRETTY"},
			},
			&cli.BoolFlag{
				Value:   false,
				Name:    "log-color",
				Usage:   "Enable colored logging",
				EnvVars: []string{"OCIS_LOG_COLOR"},
			},
		}
	})

	wopts := []web.Option{
		web.Name(
			strings.Join(
				[]string{
					sopts.Namespace,
					sopts.Name,
				},
				".",
			),
		),
		web.Version(sopts.Version),
		web.Address(sopts.Address),
		web.RegisterTTL(time.Second * 30),
		web.RegisterInterval(time.Second * 10),
		web.Context(sopts.Context),
		web.TLSConfig(sopts.TLSConfig),
		web.Handler(sopts.Handler),
		web.Flags(flags...),
	}

	return Service{
		web.NewService(
			wopts...,
		),
	}
}

func transport(secure *tls.Config) string {
	if secure != nil {
		return "https"
	}

	return "http"
}
