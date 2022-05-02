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

	"github.com/owncloud/ocis/ocis-pkg/shared"

	mzlog "github.com/go-micro/plugins/v4/logger/zerolog"
	"github.com/mohae/deepcopy"
	"github.com/olekukonko/tablewriter"

	accounts "github.com/owncloud/ocis/extensions/accounts/pkg/command"
	appprovider "github.com/owncloud/ocis/extensions/appprovider/pkg/command"
	authbasic "github.com/owncloud/ocis/extensions/auth-basic/pkg/command"
	authbearer "github.com/owncloud/ocis/extensions/auth-bearer/pkg/command"
	authmachine "github.com/owncloud/ocis/extensions/auth-machine/pkg/command"
	frontend "github.com/owncloud/ocis/extensions/frontend/pkg/command"
	gateway "github.com/owncloud/ocis/extensions/gateway/pkg/command"
	glauth "github.com/owncloud/ocis/extensions/glauth/pkg/command"
	graphExplorer "github.com/owncloud/ocis/extensions/graph-explorer/pkg/command"
	graph "github.com/owncloud/ocis/extensions/graph/pkg/command"
	group "github.com/owncloud/ocis/extensions/group/pkg/command"
	idm "github.com/owncloud/ocis/extensions/idm/pkg/command"
	idp "github.com/owncloud/ocis/extensions/idp/pkg/command"
	nats "github.com/owncloud/ocis/extensions/nats/pkg/command"
	notifications "github.com/owncloud/ocis/extensions/notifications/pkg/command"
	ocdav "github.com/owncloud/ocis/extensions/ocdav/pkg/command"
	ocs "github.com/owncloud/ocis/extensions/ocs/pkg/command"
	proxy "github.com/owncloud/ocis/extensions/proxy/pkg/command"
	settings "github.com/owncloud/ocis/extensions/settings/pkg/command"
	sharing "github.com/owncloud/ocis/extensions/sharing/pkg/command"
	storagemetadata "github.com/owncloud/ocis/extensions/storage-metadata/pkg/command"
	storagepublic "github.com/owncloud/ocis/extensions/storage-publiclink/pkg/command"
	storageshares "github.com/owncloud/ocis/extensions/storage-shares/pkg/command"
	storageusers "github.com/owncloud/ocis/extensions/storage-users/pkg/command"
	store "github.com/owncloud/ocis/extensions/store/pkg/command"
	thumbnails "github.com/owncloud/ocis/extensions/thumbnails/pkg/command"
	user "github.com/owncloud/ocis/extensions/user/pkg/command"
	web "github.com/owncloud/ocis/extensions/web/pkg/command"
	webdav "github.com/owncloud/ocis/extensions/webdav/pkg/command"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/log"
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

	s.ServicesRegistry[opts.Config.Settings.Service.Name] = settings.NewSutureService
	s.ServicesRegistry[opts.Config.Nats.Service.Name] = nats.NewSutureService
	s.ServicesRegistry[opts.Config.StorageMetadata.Service.Name] = storagemetadata.NewStorageMetadata
	s.ServicesRegistry[opts.Config.GLAuth.Service.Name] = glauth.NewSutureService
	s.ServicesRegistry[opts.Config.Graph.Service.Name] = graph.NewSutureService
	s.ServicesRegistry[opts.Config.GraphExplorer.Service.Name] = graphExplorer.NewSutureService
	s.ServicesRegistry[opts.Config.IDM.Service.Name] = idm.NewSutureService
	s.ServicesRegistry[opts.Config.OCS.Service.Name] = ocs.NewSutureService
	s.ServicesRegistry[opts.Config.Store.Service.Name] = store.NewSutureService
	s.ServicesRegistry[opts.Config.Thumbnails.Service.Name] = thumbnails.NewSutureService
	s.ServicesRegistry[opts.Config.Web.Service.Name] = web.NewSutureService
	s.ServicesRegistry[opts.Config.WebDAV.Service.Name] = webdav.NewSutureService
	s.ServicesRegistry[opts.Config.Frontend.Service.Name] = frontend.NewFrontend
	s.ServicesRegistry[opts.Config.OCDav.Service.Name] = ocdav.NewOCDav
	s.ServicesRegistry[opts.Config.Gateway.Service.Name] = gateway.NewGateway
	s.ServicesRegistry[opts.Config.User.Service.Name] = user.NewUserProvider
	s.ServicesRegistry[opts.Config.Group.Service.Name] = group.NewGroupProvider
	s.ServicesRegistry[opts.Config.AuthBasic.Service.Name] = authbasic.NewAuthBasic
	s.ServicesRegistry[opts.Config.AuthBearer.Service.Name] = authbearer.NewAuthBearer
	s.ServicesRegistry[opts.Config.AuthMachine.Service.Name] = authmachine.NewAuthMachine
	s.ServicesRegistry[opts.Config.StorageUsers.Service.Name] = storageusers.NewStorageUsers
	s.ServicesRegistry[opts.Config.StorageShares.Service.Name] = storageshares.NewStorageShares
	s.ServicesRegistry[opts.Config.StoragePublicLink.Service.Name] = storagepublic.NewStoragePublicLink
	s.ServicesRegistry[opts.Config.AppProvider.Service.Name] = appprovider.NewAppProvider
	s.ServicesRegistry[opts.Config.Notifications.Service.Name] = notifications.NewSutureService

	// populate delayed services
	s.Delayed[opts.Config.Sharing.Service.Name] = sharing.NewSharing
	s.Delayed[opts.Config.Accounts.Service.Name] = accounts.NewSutureService
	s.Delayed[opts.Config.Proxy.Service.Name] = proxy.NewSutureService
	s.Delayed[opts.Config.IDP.Service.Name] = idp.NewSutureService

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

	if s.cfg.Commons == nil {
		s.cfg.Commons = &shared.Commons{
			Log: &shared.Log{},
		}
	}

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
			if _, err = net.Dial("tcp", net.JoinHostPort(s.cfg.Runtime.Host, s.cfg.Runtime.Port)); err != nil {
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
func (s *Service) generateRunSet(cfg *ociscfg.Config) {
	if cfg.Runtime.Extensions != "" {
		e := strings.Split(strings.ReplaceAll(cfg.Runtime.Extensions, " ", ""), ",")
		for i := range e {
			runset = append(runset, e[i])
		}
		return
	}

	for name := range s.ServicesRegistry {
		// don't run glauth by default but keep the possiblity to start it via cfg.Runtime.Extensions for now
		if name == "glauth" {
			continue
		}
		runset = append(runset, name)
	}

	for name := range s.Delayed {
		// don't run accounts by default but keep the possiblity to start it via cfg.Runtime.Extensions for now
		if name == "accounts" {
			continue
		}
		runset = append(runset, name)
	}
}

// Start indicates the Service Controller to start a new supervised service as an OS thread.
func (s *Service) Start(name string, reply *int) error {
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
			if err := s.Supervisor.RemoveAndWait(s.serviceToken[name][i], 5000*time.Millisecond); err != nil {
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
