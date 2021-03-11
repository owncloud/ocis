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

	"github.com/mohae/deepcopy"

	mzlog "github.com/asim/go-micro/plugins/logger/zerolog/v3"
	accounts "github.com/owncloud/ocis/accounts/pkg/command"
	ocs "github.com/owncloud/ocis/ocs/pkg/command"
	onlyoffice "github.com/owncloud/ocis/onlyoffice/pkg/command"
	proxy "github.com/owncloud/ocis/proxy/pkg/command"
	store "github.com/owncloud/ocis/store/pkg/command"
	thumbnails "github.com/owncloud/ocis/thumbnails/pkg/command"
	web "github.com/owncloud/ocis/web/pkg/command"
	webdav "github.com/owncloud/ocis/webdav/pkg/command"

	"github.com/asim/go-micro/v3/logger"
	"github.com/thejerf/suture"

	"github.com/olekukonko/tablewriter"
	glauth "github.com/owncloud/ocis/glauth/pkg/command"
	idp "github.com/owncloud/ocis/idp/pkg/command"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
	settings "github.com/owncloud/ocis/settings/pkg/command"
	storage "github.com/owncloud/ocis/storage/pkg/command"
	"github.com/rs/zerolog"
)

// Service represents a RPC service.
type Service struct {
	Supervisor       *suture.Supervisor
	ServicesRegistry map[string]func(context.Context, *ociscfg.Config) suture.Service
	Delayed          map[string]func(context.Context, *ociscfg.Config) suture.Service
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
		ServicesRegistry: make(map[string]func(context.Context, *ociscfg.Config) suture.Service),
		Delayed:          make(map[string]func(context.Context, *ociscfg.Config) suture.Service),
		Log:              l,

		serviceToken: make(map[string][]suture.ServiceToken),
		context:      globalCtx,
		cancel:       cancelGlobal,
		cfg:          opts.Config,
	}

	s.ServicesRegistry["settings"] = settings.NewSutureService
	s.ServicesRegistry["storage-metadata"] = storage.NewStorageMetadata
	s.ServicesRegistry["glauth"] = glauth.NewSutureService
	s.ServicesRegistry["idp"] = idp.NewSutureService
	s.ServicesRegistry["ocs"] = ocs.NewSutureService
	s.ServicesRegistry["onlyoffice"] = onlyoffice.NewSutureService
	//s.ServicesRegistry["proxy"] = proxy.NewSutureService
	s.ServicesRegistry["store"] = store.NewSutureService
	s.ServicesRegistry["thumbnails"] = thumbnails.NewSutureService
	s.ServicesRegistry["web"] = web.NewSutureService
	s.ServicesRegistry["webdav"] = webdav.NewSutureService
	s.ServicesRegistry["storage-frontend"] = storage.NewFrontend
	s.ServicesRegistry["storage-gateway"] = storage.NewGateway
	s.ServicesRegistry["storage-users-provider"] = storage.NewUsersProviderService
	s.ServicesRegistry["storage-groupsprovider"] = storage.NewGroupsProvider
	s.ServicesRegistry["storage-authbasic"] = storage.NewAuthBasic
	s.ServicesRegistry["storage-authbearer"] = storage.NewAuthBearer
	s.ServicesRegistry["storage-home"] = storage.NewStorageHome
	s.ServicesRegistry["storage-users"] = storage.NewStorageUsers
	s.ServicesRegistry["storage-public-link"] = storage.NewStoragePublicLink

	// populate delayed services
	s.Delayed["storage-sharing"] = storage.NewSharing
	s.Delayed["accounts"] = accounts.NewSutureService
	s.Delayed["proxy"] = proxy.NewSutureService

	return s, nil
}

// Start an rpc service. By default the package scope Start will run all default extensions to provide with a working
// oCIS instance.
func Start(o ...Option) error {
	// prepare a new rpc Service struct.
	s, err := NewService(o...)
	if err != nil {
		return err
	}

	setMicroLogger()

	// Start creates its own supervisor. Running services under `ocis server` will create its own supervision tree.
	s.Supervisor = suture.New("ocis", suture.Spec{
		Log: func(msg string) {
			s.Log.Info().Str("message", msg).Msg("supervisor: ")
		},
	})

	// reva storages have their own logging. For consistency sake the top level logging will be cascade to reva.
	s.cfg.Storage.Log.Color = s.cfg.Log.Color
	s.cfg.Storage.Log.Level = s.cfg.Log.Level
	s.cfg.Storage.Log.Pretty = s.cfg.Log.Pretty

	// notify goroutines that they are running on supervised mode
	s.cfg.Mode = ociscfg.SUPERVISED

	if err := rpc.Register(s); err != nil {
		if s != nil {
			s.Log.Fatal().Err(err)
		}
	}
	rpc.HandleHTTP()

	// halt listens for interrupt signals and blocks.
	halt := make(chan os.Signal, 1)
	signal.Notify(halt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// TODO(refs) change default port
	l, err := net.Listen("tcp", fmt.Sprintf("%v:%v", "localhost", "6060"))
	if err != nil {
		s.Log.Fatal().Err(err)
	}

	defer func() {
		if r := recover(); r != nil {
			reason := strings.Builder{}
			// TODO(refs) change default port
			if _, err := net.Dial("localhost", "6060"); err != nil {
				reason.WriteString("runtime address already in use")
			}

			fmt.Println(reason.String())
		}
	}()

	for name := range s.ServicesRegistry {
		swap := deepcopy.Copy(s.cfg)
		s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(s.ServicesRegistry[name](s.context, swap.(*ociscfg.Config))))
	}

	go s.Supervisor.ServeBackground()

	// trap will block on halt channel for interruptions.
	go trap(s, halt)

	time.Sleep(2 * time.Second)
	// add delayed
	for name := range s.Delayed {
		swap := deepcopy.Copy(s.cfg)
		s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(s.Delayed[name](s.context, swap.(*ociscfg.Config))))
	}

	return http.Serve(l, nil)
}

// Start indicates the Service Controller to start a new supervised service as an OS thread.
func (s *Service) Start(name string, reply *int) error {
	swap := deepcopy.Copy(s.cfg)
	if _, ok := s.ServicesRegistry[name]; ok {
		*reply = 0
		s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(s.ServicesRegistry[name](s.context, swap.(*ociscfg.Config))))
		return nil
	}

	if _, ok := s.Delayed[name]; ok {
		*reply = 0
		s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(s.Delayed[name](s.context, swap.(*ociscfg.Config))))
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
