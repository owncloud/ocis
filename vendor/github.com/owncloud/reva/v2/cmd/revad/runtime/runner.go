package runtime

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/owncloud/reva/v2/pkg/registry"
	"github.com/owncloud/reva/v2/pkg/rgrpc"
	"github.com/owncloud/reva/v2/pkg/rhttp"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type servers struct {
	rhttpServer             *rhttp.Server
	rgrpcServer             *rgrpc.Server
	gracefulShutdownTimeout int
	pidFilel                string
	log                     *zerolog.Logger
}

type Runner interface {
	Start() error
	Stop(ctx context.Context) error
}

// RunServerWithOptions runs a revad server w/o wacher with the given config file, pid file and options.
func RunServerWithOptions(mainConf map[string]interface{}, pidFile string, opts ...Option) Runner {
	options := newOptions(opts...)
	parseSharedConfOrDie(mainConf["shared"])
	coreConf := parseCoreConfOrDie(mainConf["core"])

	if err := registry.Init(options.Registry); err != nil {
		panic(err)
	}

	host, _ := os.Hostname()
	log := options.Logger
	log.Info().Msgf("host info: %s", host)

	// Only initialise tracing if we didn't get a tracer provider.
	if options.TraceProvider == nil {
		log.Debug().Msg("no pre-existing tracer given, initializing tracing")
		options.TraceProvider = initTracing(coreConf)
	}
	initCPUCount(coreConf, log)

	server := &servers{
		rhttpServer: initHTTPServer(mainConf, &options),
		rgrpcServer: initGRPCServer(mainConf, &options),
		log:         log,
		pidFilel:    pidFile,
	}
	server.gracefulShutdownTimeout = 30
	if coreConf.GracefulShutdownTimeout > 0 {
		server.gracefulShutdownTimeout = coreConf.GracefulShutdownTimeout
	}

	if server.rhttpServer == nil && server.rgrpcServer == nil {
		log.Fatal().Msg("nothing to do, no grpc/http enabled_services declared in config")
	}
	return server
}

func (s *servers) Start() error {
	eg := new(errgroup.Group)
	if s.rhttpServer != nil {
		eg.Go(s.startHTTPServer)
	}
	if s.rgrpcServer != nil {
		eg.Go(s.startGRPCServer)
	}

	err := eg.Wait()
	if err != nil {
		s.log.Error().Err(err).Msg("error starting servers")
		return err
	}
	s.log.Info().Msg("all servers started")
	return nil

}

func (s *servers) Stop(ctx context.Context) error {
	wg := &sync.WaitGroup{}

	s.gracefulStopHTTPServer(wg)
	s.gracefulStopGRPCServer(wg)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-time.After(time.Duration(s.gracefulShutdownTimeout) * time.Second):
		s.log.Info().Msg("graceful shutdown timeout reached. running hard shutdown")
		s.stopHTTPServer()
		s.stopGRPCServer()
		return nil
	case <-done:
		s.log.Info().Msg("all revad servers gracefully stopped")
		return nil
	}
}

func (s *servers) startHTTPServer() error {
	if s.rhttpServer == nil {
		s.log.Fatal().Msg("http server not initialized")
	}
	ln, err := net.Listen(s.rhttpServer.Network(), s.rhttpServer.Address())
	if err != nil {
		s.log.Fatal().Err(err).Send()
	}
	if err = s.rhttpServer.Start(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Error().Err(err).Msg("http server error")
		return err
	}
	return nil
}

func (s *servers) startGRPCServer() error {
	if s.rgrpcServer == nil {
		s.log.Fatal().Msg("grcp server not initialized")
	}
	ln, err := net.Listen(s.rgrpcServer.Network(), s.rgrpcServer.Address())
	if err != nil {
		s.log.Fatal().Err(err).Send()
	}
	if err = s.rgrpcServer.Start(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.log.Error().Err(err).Msg("grpc server error")
		return err
	}
	return nil
}

func (s *servers) gracefulStopHTTPServer(wg *sync.WaitGroup) {
	if s.rhttpServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.log.Info().Msgf("fd to %s:%s gracefully closing", s.rhttpServer.Network(), s.rhttpServer.Address())
			if err := s.rhttpServer.GracefulStop(); err != nil {
				s.log.Error().Err(err).Msg("error stopping server")
				s.rhttpServer.Stop()
			}
		}()
	}
}

func (s *servers) gracefulStopGRPCServer(wg *sync.WaitGroup) {
	if s.rgrpcServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.log.Info().Msgf("fd to %s:%s gracefully closing", s.rgrpcServer.Network(), s.rgrpcServer.Address())
			if err := s.rgrpcServer.GracefulStop(); err != nil {
				s.log.Error().Err(err).Msg("error gracefully stopping server")
				s.rgrpcServer.Stop()
			}
		}()
	}
}

func (s *servers) stopHTTPServer() {
	if s.rhttpServer == nil {
		return
	}
	s.log.Info().Msgf("fd to %s:%s abruptly closing", s.rhttpServer.Network(), s.rhttpServer.Address())
	err := s.rhttpServer.Stop()
	if err != nil {
		s.log.Error().Err(err).Msg("error stopping server")
	}
}

func (s *servers) stopGRPCServer() {
	if s.rgrpcServer == nil {
		return
	}
	s.log.Info().Msgf("fd to %s:%s abruptly closing", s.rgrpcServer.Network(), s.rgrpcServer.Address())
	err := s.rgrpcServer.Stop()
	if err != nil {
		s.log.Error().Err(err).Msg("error stopping server")
	}
}

func initHTTPServer(mainConf map[string]interface{}, options *Options) *rhttp.Server {
	if isEnabledHTTP(mainConf) {
		s, err := getHTTPServer(mainConf["http"], options.Logger, options.TraceProvider)
		if err != nil {
			options.Logger.Fatal().Err(err).Msg("error creating http server")
		}
		return s
	}
	return nil
}

func initGRPCServer(mainConf map[string]interface{}, options *Options) *rgrpc.Server {
	if isEnabledGRPC(mainConf) {
		s, err := getGRPCServer(mainConf["grpc"], options.Logger, options.TraceProvider)
		if err != nil {
			options.Logger.Fatal().Err(err).Msg("error creating grpc server")
		}
		return s
	}
	return nil
}
