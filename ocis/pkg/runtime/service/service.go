package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	mzlog "github.com/asim/go-micro/plugins/logger/zerolog/v4"
	"github.com/mohae/deepcopy"
	"github.com/olekukonko/tablewriter"
	accounts "github.com/owncloud/ocis/accounts/pkg/command"
	glauth "github.com/owncloud/ocis/glauth/pkg/command"
	graphExplorer "github.com/owncloud/ocis/graph-explorer/pkg/command"
	graph "github.com/owncloud/ocis/graph/pkg/command"
	idp "github.com/owncloud/ocis/idp/pkg/command"
	"github.com/owncloud/ocis/ocis-pkg/config"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	ocs "github.com/owncloud/ocis/ocs/pkg/command"
	proxy "github.com/owncloud/ocis/proxy/pkg/command"
	settings "github.com/owncloud/ocis/settings/pkg/command"
	storage "github.com/owncloud/ocis/storage/pkg/command"
	store "github.com/owncloud/ocis/store/pkg/command"
	thumbnails "github.com/owncloud/ocis/thumbnails/pkg/command"
	web "github.com/owncloud/ocis/web/pkg/command"
	webdav "github.com/owncloud/ocis/webdav/pkg/command"
	"github.com/rs/zerolog"
	"github.com/thejerf/suture/v4"
	"go-micro.dev/v4/logger"
)

var (
	// runset keeps track of which extensions to start supervised.
	runset []string
)

type serviceFuncMap map[string]func(*ociscfg.Config) suture.Service

// Service represents a RPC service.
type Service struct {
	Supervisor       *suture.Supervisor
	ServicesRegistry serviceFuncMap
	Delayed          serviceFuncMap
	Log              log.Logger

	serviceToken map[string][]suture.ServiceToken
	context      context.Context
	cancel       context.CancelFunc
	cfg          *ociscfg.Config
}

// NewService returns a configured service with a controller and a default logger.
// When used as a library, flags are not parsed, and in order to avoid introducing a global state with init functions
// calls are done explicitly to loadFromEnv().
// Since this is the public constructor, options need to be added, at the moment only logging options
// are supported in order to match the running OwnCloud services structured log.
func NewService(options ...Option) (*Service, error) {
	opts := NewOptions()

	for _, f := range options {
		f(opts)
	}

	l := log.NewLogger(
		log.Color(opts.Config.Log.Color),
		log.Pretty(opts.Config.Log.Pretty),
		log.Level(opts.Config.Log.Level),
	)

	globalCtx, cancelGlobal := context.WithCancel(context.Background())

	s := &Service{
		ServicesRegistry: make(serviceFuncMap),
		Delayed:          make(serviceFuncMap),
		Log:              l,

		serviceToken: make(map[string][]suture.ServiceToken),
		context:      globalCtx,
		cancel:       cancelGlobal,
		cfg:          opts.Config,
	}

	s.ServicesRegistry["settings"] = settings.NewSutureService
	s.ServicesRegistry["storage-metadata"] = storage.NewStorageMetadata
	s.ServicesRegistry["glauth"] = glauth.NewSutureService
	s.ServicesRegistry["graph"] = graph.NewSutureService
	s.ServicesRegistry["graph-explorer"] = graphExplorer.NewSutureService
	s.ServicesRegistry["idp"] = idp.NewSutureService
	s.ServicesRegistry["ocs"] = ocs.NewSutureService
	s.ServicesRegistry["store"] = store.NewSutureService
	s.ServicesRegistry["thumbnails"] = thumbnails.NewSutureService
	s.ServicesRegistry["web"] = web.NewSutureService
	s.ServicesRegistry["webdav"] = webdav.NewSutureService
	s.ServicesRegistry["storage-frontend"] = storage.NewFrontend
	s.ServicesRegistry["storage-gateway"] = storage.NewGateway
	s.ServicesRegistry["storage-userprovider"] = storage.NewUserProvider
	s.ServicesRegistry["storage-groupprovider"] = storage.NewGroupProvider
	s.ServicesRegistry["storage-authbasic"] = storage.NewAuthBasic
	s.ServicesRegistry["storage-authbearer"] = storage.NewAuthBearer
	s.ServicesRegistry["storage-authmachine"] = storage.NewAuthMachine
	s.ServicesRegistry["storage-home"] = storage.NewStorageHome
	s.ServicesRegistry["storage-users"] = storage.NewStorageUsers
	s.ServicesRegistry["storage-public-link"] = storage.NewStoragePublicLink
	s.ServicesRegistry["storage-appprovider"] = storage.NewAppProvider

	// populate delayed services
	s.Delayed["storage-sharing"] = storage.NewSharing
	s.Delayed["accounts"] = accounts.NewSutureService
	s.Delayed["proxy"] = proxy.NewSutureService

	return s, nil
}

// Start an rpc service. By default the package scope Start will run all default extensions to provide with a working
// oCIS instance.
func Start(o ...Option) error {
	// Start the runtime. Most likely this was called ONLY by the `ocis server` subcommand, but since we cannot protect
	// from the caller, the previous statement holds truth.

	// prepare a new rpc Service struct.
	s, err := NewService(o...)
	if err != nil {
		return err
	}

	// halt listens for interrupt signals and blocks.
	halt := make(chan os.Signal, 1)
	signal.Notify(halt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// notify goroutines that they are running on supervised mode
	s.cfg.Mode = ociscfg.SUPERVISED

	setMicroLogger()

	// tolerance controls backoff cycles from the supervisor.
	tolerance := 5
	totalBackoff := 0

	// Start creates its own supervisor. Running services under `ocis server` will create its own supervision tree.
	s.Supervisor = suture.New("ocis", suture.Spec{
		EventHook: func(e suture.Event) {
			if e.Type() == suture.EventTypeBackoff {
				totalBackoff++
				if totalBackoff == tolerance {
					halt <- os.Interrupt
				}
			}
			s.Log.Info().Str("event", e.String()).Msg(fmt.Sprintf("supervisor: %v", e.Map()["supervisor_name"]))
		},
		FailureThreshold: 5,
		FailureBackoff:   3 * time.Second,
	})

	if err = rpc.Register(s); err != nil {
		if s != nil {
			s.Log.Fatal().Err(err)
		}
	}
	rpc.HandleHTTP()

	l, err := net.Listen("tcp", net.JoinHostPort(s.cfg.Runtime.Host, s.cfg.Runtime.Port))
	if err != nil {
		s.Log.Fatal().Err(err)
	}

	defer func() {
		if r := recover(); r != nil {
			reason := strings.Builder{}
			if _, err := net.Dial("tcp", net.JoinHostPort(s.cfg.Runtime.Host, s.cfg.Runtime.Port)); err != nil {
				reason.WriteString("runtime address already in use")
			}

			fmt.Println(reason.String())
		}
	}()

	// prepare the set of services to run
	s.generateRunSet(s.cfg)

	// schedule services that we are sure don't have interdependencies.
	scheduleServiceTokens(s, s.ServicesRegistry)

	// there are reasons not to do this, but we have race conditions ourselves. Until we resolve them, mind the following disclaimer:
	// Calling ServeBackground will CORRECTLY start the supervisor running in a new goroutine. It is risky to directly run
	// go supervisor.Serve()
	// because that will briefly create a race condition as it starts up, if you try to .Add() services immediately afterward.
	// https://pkg.go.dev/github.com/thejerf/suture/v4@v4.0.0#Supervisor
	go s.Supervisor.ServeBackground(s.context)

	// trap will block on halt channel for interruptions.
	go trap(s, halt)

	// add services with delayed execution.
	time.Sleep(1 * time.Second)
	scheduleServiceTokens(s, s.Delayed)

	return http.Serve(l, nil)
}

// scheduleServiceTokens adds service tokens to the service supervisor.
func scheduleServiceTokens(s *Service, funcSet serviceFuncMap) {
	for _, name := range runset {
		if _, ok := funcSet[name]; !ok {
			continue
		}

		swap := deepcopy.Copy(s.cfg)
		s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(funcSet[name](swap.(*ociscfg.Config))))
	}
}

// generateRunSet interprets the cfg.Runtime.Extensions config option to cherry-pick which services to start using
// the runtime.
func (s *Service) generateRunSet(cfg *config.Config) {
	if cfg.Runtime.Extensions != "" {
		e := strings.Split(strings.ReplaceAll(cfg.Runtime.Extensions, " ", ""), ",")
		for i := range e {
			runset = append(runset, e[i])
		}
		return
	}

	for name := range s.ServicesRegistry {
		runset = append(runset, name)
	}

	for name := range s.Delayed {
		runset = append(runset, name)
	}
}

// Start indicates the Service Controller to start a new supervised service as an OS thread.
func (s *Service) Start(name string, reply *int) error {
	// RPC calls to a Service object will allow for parsing config. Mind that since the runtime is running on a different
	// machine, the configuration needs to be present in the given machine. RPC does not yet allow providing a config
	// during transport.
	s.cfg.Mode = ociscfg.UNSUPERVISED

	swap := deepcopy.Copy(s.cfg)
	if _, ok := s.ServicesRegistry[name]; ok {
		*reply = 0
		s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(s.ServicesRegistry[name](swap.(*ociscfg.Config))))
		return nil
	}

	if _, ok := s.Delayed[name]; ok {
		*reply = 0
		s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(s.Delayed[name](swap.(*ociscfg.Config))))
		return nil
	}

	*reply = 0
	return fmt.Errorf("cannot start service %s: unknown service", name)
}

// List running processes for the Service Controller.
func (s *Service) List(args struct{}, reply *string) error {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Extension"})

	names := []string{}
	for t := range s.serviceToken {
		if len(s.serviceToken[t]) > 0 {
			names = append(names, t)
		}
	}

	sort.Strings(names)

	for n := range names {
		table.Append([]string{names[n]})
	}

	table.Render()
	*reply = tableString.String()
	return nil
}

// Kill a supervised process by subcommand name.
func (s *Service) Kill(name string, reply *int) error {
	if len(s.serviceToken[name]) > 0 {
		for i := range s.serviceToken[name] {
			if err := s.Supervisor.Remove(s.serviceToken[name][i]); err != nil {
				return err
			}
		}
		delete(s.serviceToken, name)
	} else {
		return fmt.Errorf("service %s not found", name)
	}

	return nil
}

// trap blocks on halt channel. When the runtime is interrupted it
// signals the controller to stop any supervised process.
func trap(s *Service, halt chan os.Signal) {
	<-halt
	s.cancel()
	for sName := range s.serviceToken {
		for i := range s.serviceToken[sName] {
			if err := s.Supervisor.Remove(s.serviceToken[sName][i]); err != nil {
				s.Log.Error().Err(err).Str("service", "runtime service").Msgf("terminating with signal: %v", s)
			}
		}
	}
	s.Log.Debug().Str("service", "runtime service").Msgf("terminating with signal: %v", s)
	os.Exit(0)
}

// for logging reasons we don't want the same logging level on both oCIS and micro. As a framework builder we do not
// want to expose to the end user the internal framework logs unless explicitly specified.
func setMicroLogger() {
	if os.Getenv("MICRO_LOG_LEVEL") == "" {
		_ = os.Setenv("MICRO_LOG_LEVEL", "error")
	}

	lev, err := zerolog.ParseLevel(os.Getenv("MICRO_LOG_LEVEL"))
	if err != nil {
		lev = zerolog.ErrorLevel
	}
	logger.DefaultLogger = mzlog.NewLogger(logger.WithLevel(logger.Level(lev)))
}
