package service

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/owncloud/ocis/ocis/pkg/runtime/config"
	"github.com/owncloud/ocis/ocis/pkg/runtime/controller"
	"github.com/owncloud/ocis/ocis/pkg/runtime/log"
	"github.com/owncloud/ocis/ocis/pkg/runtime/process"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var (
	halt = make(chan os.Signal, 1)
	done = make(chan struct{}, 1)
)

// Service represents a RPC service.
type Service struct {
	Controller controller.Controller
	Log        zerolog.Logger
	wg         *sync.WaitGroup
	done       bool
}

// loadFromEnv would set cmd global variables. This is a workaround spf13/viper since pman used as a library does not
// parse flags.
func loadFromEnv() *config.Config {
	cfg := config.NewConfig()
	viper.AutomaticEnv()

	viper.BindEnv("keep-alive", "RUNTIME_KEEP_ALIVE")
	viper.BindEnv("port", "RUNTIME_PORT")

	cfg.KeepAlive = viper.GetBool("keep-alive")

	if viper.GetString("port") != "" {
		cfg.Port = viper.GetString("port")
	}

	return cfg
}

// NewService returns a configured service with a controller and a default logger.
// When used as a library, flags are not parsed, and in order to avoid introducing a global state with init functions
// calls are done explicitly to loadFromEnv().
// Since this is the public constructor, options need to be added, at the moment only logging options
// are supported in order to match the running OwnCloud services structured log.
func NewService(options ...Option) *Service {
	opts := NewOptions()

	for _, f := range options {
		f(opts)
	}

	cfg := loadFromEnv()
	l := log.NewLogger(
		log.WithPretty(opts.Log.Pretty),
	)

	return &Service{
		wg:  &sync.WaitGroup{},
		Log: l,
		Controller: controller.NewController(
			controller.WithConfig(cfg),
			controller.WithLog(&l),
		),
	}
}

// Start an rpc service.
func Start(o ...Option) error {
	s := NewService(o...)

	if err := rpc.Register(s); err != nil {
		s.Log.Fatal().Err(err)
	}
	rpc.HandleHTTP()

	signal.Notify(halt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", s.Controller.Config.Hostname, s.Controller.Config.Port))
	if err != nil {
		s.Log.Fatal().Err(err)
	}

	// handle panic within the Service scope.
	defer func() {
		if r := recover(); r != nil {
			reason := strings.Builder{}
			if _, err := net.Dial("localhost", s.Controller.Config.Port); err != nil {
				reason.WriteString("runtime address already in use")
			}

			fmt.Println(reason.String())
		}
	}()

	go trap(s)

	return http.Serve(l, nil)
}

// Start indicates the Service Controller to start a new supervised service as an OS thread.
func (s *Service) Start(args process.ProcEntry, reply *int) error {
	if !s.done {
		s.wg.Add(1)
		s.Log.Info().Str("service", args.Extension).Msgf("%v", "started")
		if err := s.Controller.Start(args); err != nil {
			*reply = 1
			return err
		}

		*reply = 0
		s.wg.Done()
	}

	return nil
}

// List running processes for the Service Controller.
func (s *Service) List(args struct{}, reply *string) error {
	*reply = s.Controller.List()
	return nil
}

// Kill a supervised process by subcommand name.
func (s *Service) Kill(args *string, reply *int) error {
	pe := process.ProcEntry{
		Extension: *args,
	}
	if err := s.Controller.Kill(pe); err != nil {
		*reply = 1
		return err
	}

	*reply = 0
	return nil
}

// trap blocks on halt channel. When the runtime is interrupted it
// signals the controller to stop any supervised process.
func trap(s *Service) {
	<-halt
	s.done = true
	s.wg.Wait()
	s.Log.Debug().
		Str("service", "runtime service").
		Msgf("terminating with signal: %v", s)
	if err := s.Controller.Shutdown(done); err != nil {
		s.Log.Err(err)
	}
	close(done)
	os.Exit(0)
}
