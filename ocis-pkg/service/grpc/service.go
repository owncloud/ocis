package grpc

import (
	"strings"
	"time"

	"github.com/owncloud/ocis/ocis-pkg/service/http"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	mclient "github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/server"
	sgrpc "github.com/micro/go-micro/v2/server/grpc"

	"github.com/micro/go-plugins/wrapper/trace/opencensus/v2"
	"github.com/owncloud/ocis/ocis-pkg/registry"
	"github.com/owncloud/ocis/ocis-pkg/wrapper/prometheus"
)

// DefaultClient is a custom ocis grpc configured client.
var DefaultClient = newGrpcClient()

func newGrpcClient() mclient.Client {
	r := *registry.GetRegistry()

	c := grpc.NewClient(
		mclient.RequestTimeout(10*time.Second),
		mclient.Registry(r),
	)
	return c
}

// Service simply wraps the go-micro grpc service.
type Service struct {
	micro.Service
}

// NewService initializes a new grpc service.
func NewService(opts ...Option) Service {
	sopts := newOptions(opts...)
	sopts.Logger.Info().
		Str("transport", "grpc").
		Str("addr", sopts.Address).
		Msg("starting server")

	flags := make([]cli.Flag, 0)
	http.M.Lock()
	// There is a global state somewhere in go-micro. In order to circumvent this, flag registration down to the go-micro
	// flags need to be parsed ONLY once, in other words, this roughly translates into: let go-micro know we're using
	// our own set of flags, so it does not panic when it encounters a flag it does not recognise.
	http.Once.Do(func() {
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
	sname := strings.Join(
		[]string{
			sopts.Namespace,
			sopts.Name,
		},
		".",
	)

	mopts := []micro.Option{
		micro.Name(sname),
		micro.Client(DefaultClient),
		micro.Version(sopts.Version),
		micro.Address(sopts.Address),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
		micro.WrapClient(opencensus.NewClientWrapper()),
		micro.WrapHandler(opencensus.NewHandlerWrapper()),
		micro.WrapSubscriber(opencensus.NewSubscriberWrapper()),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
		micro.Context(sopts.Context),
		micro.Flags(flags...),
		micro.Server(sgrpc.NewServer(server.Name(sname), server.Address(sopts.Address))),
	}

	return Service{
		micro.NewService(
			mopts...,
		),
	}
}
