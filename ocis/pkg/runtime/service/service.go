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
	"github.com/olekukonko/tablewriter"
	ociscfg "github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
	appProvider "github.com/owncloud/ocis/v2/services/app-provider/pkg/command"
	appRegistry "github.com/owncloud/ocis/v2/services/app-registry/pkg/command"
	authbasic "github.com/owncloud/ocis/v2/services/auth-basic/pkg/command"
	authmachine "github.com/owncloud/ocis/v2/services/auth-machine/pkg/command"
	eventhistory "github.com/owncloud/ocis/v2/services/eventhistory/pkg/command"
	frontend "github.com/owncloud/ocis/v2/services/frontend/pkg/command"
	gateway "github.com/owncloud/ocis/v2/services/gateway/pkg/command"
	graph "github.com/owncloud/ocis/v2/services/graph/pkg/command"
	groups "github.com/owncloud/ocis/v2/services/groups/pkg/command"
	idm "github.com/owncloud/ocis/v2/services/idm/pkg/command"
	idp "github.com/owncloud/ocis/v2/services/idp/pkg/command"
	nats "github.com/owncloud/ocis/v2/services/nats/pkg/command"
	notifications "github.com/owncloud/ocis/v2/services/notifications/pkg/command"
	ocdav "github.com/owncloud/ocis/v2/services/ocdav/pkg/command"
	ocs "github.com/owncloud/ocis/v2/services/ocs/pkg/command"
	postprocessing "github.com/owncloud/ocis/v2/services/postprocessing/pkg/command"
	proxy "github.com/owncloud/ocis/v2/services/proxy/pkg/command"
	search "github.com/owncloud/ocis/v2/services/search/pkg/command"
	settings "github.com/owncloud/ocis/v2/services/settings/pkg/command"
	sharing "github.com/owncloud/ocis/v2/services/sharing/pkg/command"
	storagepublic "github.com/owncloud/ocis/v2/services/storage-publiclink/pkg/command"
	storageshares "github.com/owncloud/ocis/v2/services/storage-shares/pkg/command"
	storageSystem "github.com/owncloud/ocis/v2/services/storage-system/pkg/command"
	storageusers "github.com/owncloud/ocis/v2/services/storage-users/pkg/command"
	store "github.com/owncloud/ocis/v2/services/store/pkg/command"
	thumbnails "github.com/owncloud/ocis/v2/services/thumbnails/pkg/command"
	userlog "github.com/owncloud/ocis/v2/services/userlog/pkg/command"
	users "github.com/owncloud/ocis/v2/services/users/pkg/command"
	web "github.com/owncloud/ocis/v2/services/web/pkg/command"
	webdav "github.com/owncloud/ocis/v2/services/webdav/pkg/command"
	webfinger "github.com/owncloud/ocis/v2/services/webfinger/pkg/command"
	"github.com/thejerf/suture/v4"
)

var (
	// runset keeps track of which services to start supervised.
	runset map[string]struct{}
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
	s.ServicesRegistry[opts.Config.StorageSystem.Service.Name] = storageSystem.NewSutureService
	s.ServicesRegistry[opts.Config.Graph.Service.Name] = graph.NewSutureService
	s.ServicesRegistry[opts.Config.IDM.Service.Name] = idm.NewSutureService
	s.ServicesRegistry[opts.Config.OCS.Service.Name] = ocs.NewSutureService
	s.ServicesRegistry[opts.Config.Store.Service.Name] = store.NewSutureService
	s.ServicesRegistry[opts.Config.Thumbnails.Service.Name] = thumbnails.NewSutureService
	s.ServicesRegistry[opts.Config.Web.Service.Name] = web.NewSutureService
	s.ServicesRegistry[opts.Config.WebDAV.Service.Name] = webdav.NewSutureService
	s.ServicesRegistry[opts.Config.Webfinger.Service.Name] = webfinger.NewSutureService
	s.ServicesRegistry[opts.Config.Frontend.Service.Name] = frontend.NewSutureService
	s.ServicesRegistry[opts.Config.OCDav.Service.Name] = ocdav.NewSutureService
	s.ServicesRegistry[opts.Config.Gateway.Service.Name] = gateway.NewSutureService
	s.ServicesRegistry[opts.Config.AppRegistry.Service.Name] = appRegistry.NewSutureService
	s.ServicesRegistry[opts.Config.Users.Service.Name] = users.NewSutureService
	s.ServicesRegistry[opts.Config.Groups.Service.Name] = groups.NewSutureService
	s.ServicesRegistry[opts.Config.AuthBasic.Service.Name] = authbasic.NewSutureService
	s.ServicesRegistry[opts.Config.AuthMachine.Service.Name] = authmachine.NewSutureService
	s.ServicesRegistry[opts.Config.StorageUsers.Service.Name] = storageusers.NewSutureService
	s.ServicesRegistry[opts.Config.StorageShares.Service.Name] = storageshares.NewSutureService
	s.ServicesRegistry[opts.Config.StoragePublicLink.Service.Name] = storagepublic.NewSutureService
	s.ServicesRegistry[opts.Config.AppProvider.Service.Name] = appProvider.NewSutureService
	s.ServicesRegistry[opts.Config.Notifications.Service.Name] = notifications.NewSutureService
	s.ServicesRegistry[opts.Config.Search.Service.Name] = search.NewSutureService
	s.ServicesRegistry[opts.Config.Postprocessing.Service.Name] = postprocessing.NewSutureService
	s.ServicesRegistry[opts.Config.EventHistory.Service.Name] = eventhistory.NewSutureService
	s.ServicesRegistry[opts.Config.Userlog.Service.Name] = userlog.NewSutureService

	// populate delayed services
	s.Delayed[opts.Config.Sharing.Service.Name] = sharing.NewSutureService
	s.Delayed[opts.Config.Proxy.Service.Name] = proxy.NewSutureService
	s.Delayed[opts.Config.IDP.Service.Name] = idp.NewSutureService

	return s, nil
}

// Start an rpc service. By default the package scope Start will run all default services to provide with a working
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
	for name := range runset {
		if _, ok := funcSet[name]; !ok {
			continue
		}

		swap := deepcopy.Copy(s.cfg)
		s.serviceToken[name] = append(s.serviceToken[name], s.Supervisor.Add(funcSet[name](swap.(*ociscfg.Config))))
	}
}

// generateRunSet interprets the cfg.Runtime.Services config option to cherry-pick which services to start using
// the runtime.
func (s *Service) generateRunSet(cfg *ociscfg.Config) {
	runset = make(map[string]struct{})
	if cfg.Runtime.Services != "" {
		e := strings.Split(strings.ReplaceAll(cfg.Runtime.Services, " ", ""), ",")
		for _, name := range e {
			runset[name] = struct{}{}
		}
		return
	}

	for name := range s.ServicesRegistry {
		runset[name] = struct{}{}
	}

	for name := range s.Delayed {
		runset[name] = struct{}{}
	}

	if cfg.Runtime.Disabled != "" {
		e := strings.Split(strings.ReplaceAll(cfg.Runtime.Disabled, " ", ""), ",")
		for _, name := range e {
			delete(runset, name)
		}
	}
}

// List running processes for the Service Controller.
func (s *Service) List(_ struct{}, reply *string) error {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Service"})

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
